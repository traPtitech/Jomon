package router

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/internal/logging"
	"github.com/traPtitech/Jomon/internal/model"
	"github.com/traPtitech/Jomon/internal/router/wrapsession"
	"github.com/traPtitech/Jomon/internal/service"
	"github.com/traPtitech/Jomon/internal/traq"
	"go.uber.org/zap"
)

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (h Handlers) AuthCallback(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	code := c.QueryParam("code")
	if len(code) == 0 {
		return service.NewBadInputError("code is required")
	}

	codeVerifier, err := wrapsession.WithSession(
		c, h.SessionName, func(w *wrapsession.W) (string, error) {
			v, ok := w.GetCodeVerifier()
			if !ok {
				err := fmt.Errorf("code_verifier is not found in session")
				return "", service.NewUnexpectedError(err)
			}
			return v, nil
		})
	if err != nil {
		return err
	}

	res, err := traq.RequestAccessToken(code, codeVerifier)
	if err != nil {
		logger.Error("failed to get access token", zap.Error(err))
		return service.NewUnexpectedError(err)
	}

	u, err := traq.FetchTraQUserInfo(res.AccessToken)
	if err != nil {
		logger.Error("failed to fetch traQ user info", zap.Error(err))
		return service.NewUnexpectedError(err)
	}

	var modelUser *model.User
	modelUser, err = h.Repository.GetUserByName(ctx, u.Name)
	if err == nil {
		// User found, do nothing
	} else if nfErr := new(service.NotFoundError); errors.As(err, &nfErr) {
		// User not found, create new user
		modelUser, err = h.Repository.CreateUser(ctx, u.Name, u.DisplayName, false)
		if err != nil {
			logger.Error("failed to create user", zap.Error(err))
			return err
		}
	} else {
		// Some other error occurred
		logger.Error("failed to get user by name", zap.Error(err))
		return err
	}

	_, err = wrapsession.WithSession(c, h.SessionName, func(w *wrapsession.W) (struct{}, error) {
		w.SetUserID(modelUser.ID)
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

func (h Handlers) GeneratePKCE(c echo.Context) error {
	codeVerifier := randAlphabetAndNumberString(43)

	_, err := wrapsession.WithSession(c, h.SessionName, func(w *wrapsession.W) (struct{}, error) {
		w.SetCodeVerifier(codeVerifier)
		return struct{}{}, nil
	})
	if err != nil {
		return err
	}

	codeVerifierHash := sha256.Sum256([]byte(codeVerifier))
	encoder := base64.
		NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_").
		WithPadding(base64.NoPadding)

	codeChallengeMethod := "S256"

	_, err = wrapsession.WithSession(c, h.SessionName, func(w *wrapsession.W) (struct{}, error) {
		w.SetReferer(c.Request().Referer())
		return struct{}{}, nil
	})
	if err != nil {
		return err
	}

	to := traq.AuthorizeURL(
		encoder.EncodeToString(codeVerifierHash[:]),
		codeChallengeMethod,
	)

	return c.Redirect(http.StatusFound, to)
}

var randSrcPool = sync.Pool{
	New: func() interface{} {
		var b [32]byte
		if _, err := crand.Read(b[:]); err != nil {
			panic(err)
		}
		return rand.New(rand.NewChaCha8(b))
	},
}

const (
	rs6Letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	rs6LetterIdxBits = 6
	rs6LetterIdxMask = 1<<rs6LetterIdxBits - 1
	rs6LetterIdxMax  = 63 / rs6LetterIdxBits
)

func randAlphabetAndNumberString(n int) string {
	b := make([]byte, n)
	randSrc, ok := randSrcPool.Get().(*rand.Rand)
	if !ok {
		panic("failed to get rand source")
	}
	cache, remain := randSrc.Int64(), rs6LetterIdxMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int64(), rs6LetterIdxMax
		}
		idx := int(cache & rs6LetterIdxMask)
		if idx < len(rs6Letters) {
			b[i] = rs6Letters[idx]
			i--
		}
		cache >>= rs6LetterIdxBits
		remain--
	}
	randSrcPool.Put(randSrc)
	return string(b)
}
