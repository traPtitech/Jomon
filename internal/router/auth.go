package router

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/traPtitech/Jomon/internal/logging"
	"github.com/traPtitech/Jomon/internal/router/wrapsession"
	"github.com/traPtitech/Jomon/internal/service"
	"go.uber.org/zap"
)

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type authCallbackQueryParams struct {
	Code  string `query:"code"`
	State string `query:"state"`
}

func (h Handlers) bindAuthCodeWithProofs(c *echo.Context) (string, *service.AuthProofs, error) {
	var params authCallbackQueryParams
	if err := c.Bind(&params); err != nil {
		return "", nil, service.NewUnauthenticatedError("unexpected query params").WithInternal(err)
	}
	codeVerifier, err := wrapsession.WithSession(
		c, h.SessionName, func(w *wrapsession.W) (string, error) {
			v, ok := w.GetCodeVerifier()
			if !ok {
				return "", service.NewUnauthenticatedError("code_verifier not found in session")
			}
			return v, nil
		})
	if err != nil {
		return "", nil, err
	}
	nonce, err := wrapsession.WithSession(
		c, h.SessionName, func(w *wrapsession.W) (string, error) {
			v, ok := w.GetNonce()
			if !ok {
				return "", service.NewUnauthenticatedError("nonce not found in session")
			}
			return v, nil
		})
	if err != nil {
		return "", nil, err
	}
	authProofs := &service.AuthProofs{
		State:        params.State,
		CodeVerifier: codeVerifier,
		Nonce:        nonce,
	}
	return params.Code, authProofs, nil
}

func (h Handlers) AuthCallback(c *echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	code, authProofs, err := h.bindAuthCodeWithProofs(c)
	if err != nil {
		logger.Error("failed to bind auth code with proofs", zap.Error(err))
		return err
	}
	token, err := h.Service.ExchangeCodeToToken(ctx, code, authProofs)
	if err != nil {
		logger.Error("failed to exchange authorization code to token", zap.Error(err))
		return err
	}
	_, err = wrapsession.WithSession(c, h.SessionName, func(w *wrapsession.W) (struct{}, error) {
		w.SetIDToken(string(token))
		return struct{}{}, nil
	})
	if err != nil {
		return err
	}

	location, err := wrapsession.WithSession(
		c, h.SessionName, func(w *wrapsession.W) (string, error) {
			v, ok := w.GetReferer()
			if !ok {
				return "/", nil
			}
			return v, nil
		})
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, location)
}

func (h Handlers) BeginAuth(c *echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	authURL, authProofs, err := h.Service.NewAuthCodeURL(ctx)
	if err != nil {
		logger.Error("failed to create auth code URL", zap.Error(err))
		return err
	}
	_, err = wrapsession.WithSession(c, h.SessionName, func(w *wrapsession.W) (struct{}, error) {
		w.SetCodeVerifier(authProofs.CodeVerifier)
		w.SetNonce(authProofs.Nonce)
		referer := c.Request().Referer()
		if referer != "" {
			w.SetReferer(referer)
		}
		return struct{}{}, nil
	})
	if err != nil {
		logger.Error("failed to save auth proofs in session", zap.Error(err))
		return err
	}
	return c.Redirect(http.StatusSeeOther, authURL.String())
}
