package traq

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/traPtitech/Jomon/internal/service"
	"golang.org/x/oauth2"
)

type Client struct {
	clientID        string
	oauth2Config    *oauth2.Config
	idTokenVerifier *oidc.IDTokenVerifier
}

type ClientConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	ProviderURL  string
}

func LoadClient(ctx context.Context, config ClientConfig) (*Client, error) {
	provider, err := oidc.NewProvider(ctx, config.ProviderURL)
	if err != nil {
		err = fmt.Errorf("failed to create OIDC provider: %w", err)
		return nil, service.NewUnexpectedError(err)
	}
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
	}
	idTokenVerifier := provider.Verifier(&oidc.Config{
		ClientID: config.ClientID,
		// 他のアプリが発行したトークンも受け入れるため、クライアントIDのチェックはスキップする
		// 必要に応じてaudのチェックを入れる; checkIDTokenAudienceを参照
		SkipClientIDCheck:          true,
		SkipExpiryCheck:            false,
		SkipIssuerCheck:            false,
		InsecureSkipSignatureCheck: false,
	})
	return &Client{
		clientID:        config.ClientID,
		oauth2Config:    oauth2Config,
		idTokenVerifier: idTokenVerifier,
	}, nil
}
