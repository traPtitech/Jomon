package traq

import (
	"context"
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/traPtitech/Jomon/internal/service"
	"golang.org/x/oauth2"
)

// https://github.com/coreos/go-oidc/blob/da6b3bf/example/idtoken/app.go#L26-L32
func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(crand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func generateAuthProofs() (*service.AuthProofs, error) {
	state, err := randString(16)
	if err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}
	nonce, err := randString(16)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	// https://github.com/golang/oauth2/blob/4d954e6/pkce.go#L20-L38
	codeVerifier, err := randString(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate code verifier: %w", err)
	}
	return &service.AuthProofs{
		State:        state,
		CodeVerifier: codeVerifier,
		Nonce:        nonce,
	}, nil
}

func (c *Client) NewAuthCodeURL(ctx context.Context) (*url.URL, *service.AuthProofs, error) {
	authProofs, err := generateAuthProofs()
	if err != nil {
		return nil, nil, service.NewUnexpectedError(err)
	}
	rawURL := c.oauth2Config.AuthCodeURL(
		authProofs.State,
		oauth2.S256ChallengeOption(authProofs.CodeVerifier),
		oidc.Nonce(authProofs.Nonce))
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		err = fmt.Errorf("failed to parse auth code URL: %w", err)
		return nil, nil, service.NewUnexpectedError(err)
	}
	return parsedURL, authProofs, nil
}

func (c *Client) ExchangeCodeToToken(
	ctx context.Context, code string, authProofs *service.AuthProofs,
) (service.Token, error) {
	oauth2Token, err := c.oauth2Config.Exchange(ctx, code,
		oauth2.S256ChallengeOption(authProofs.CodeVerifier))
	if err != nil {
		return "", service.NewUnauthenticatedError("failed to exchange code to token").
			WithInternal(err)
	}
	idToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		err = fmt.Errorf("id_token not found in token response")
		return "", service.NewUnexpectedError(err)
	}
	return service.Token(idToken), nil
}

func (c *Client) GetUserInfo(
	ctx context.Context, token service.Token,
) (*service.OIDCUserInfo, error) {
	idToken, err := c.idTokenVerifier.Verify(ctx, string(token))
	if err != nil {
		err = service.NewUnauthenticatedError("failed to verify id token").
			WithInternal(err)
		return nil, err
	}
	var userInfo service.OIDCUserInfo
	if err := idToken.Claims(&userInfo); err != nil {
		err = service.NewUnauthenticatedError("failed to parse id token claims").
			WithInternal(err)
		return nil, err
	}
	return &userInfo, nil
}

var _ service.OIDCClient = (*Client)(nil)
