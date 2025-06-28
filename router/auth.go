package router

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand/v2"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/logging"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/router/wrapsession"
	"github.com/traPtitech/Jomon/service"
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
		return echo.NewHTTPError(http.StatusBadRequest, "code is required")
	}

	codeVerifier, err := wrapsession.WithSession(
		c, h.SessionName, func(w *wrapsession.W) (string, error) {
			v, ok := w.GetCodeVerifier()
			if !ok {
				err := echo.NewHTTPError(
					http.StatusInternalServerError,
					fmt.Errorf("code_verifier is not found in session"))
				return "", err
			}
			return v, nil
		})
	if err != nil {
		return err
	}

	res, err := service.RequestAccessToken(code, codeVerifier)
	if err != nil {
		logger.Error("failed to get access token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	u, err := service.FetchTraQUserInfo(res.AccessToken)
	if err != nil {
		logger.Error("failed to fetch traQ user info", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var modelUser *model.User
	modelUser, err = h.Repository.GetUserByName(ctx, u.Name)
	if err != nil {
		if ent.IsNotFound(err) {
			modelUser, err = h.Repository.CreateUser(ctx, u.Name, u.DisplayName, false)
			if err != nil {
				logger.Error("failed to create user", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		} else {
			logger.Error("failed to get user by name", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
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
				err := echo.NewHTTPError(
					http.StatusInternalServerError,
					fmt.Errorf("referer is not found in session"))
				return "/", err
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
		w.SetReferer(c)
		return struct{}{}, nil
	})
	if err != nil {
		return err
	}

	// nolint:lll
	to := fmt.Sprintf(
		"%s/oauth2/authorize?response_type=code&client_id=%s&code_challenge=%s&code_challenge_method=%s",
		service.TraQBaseURL, service.JomonClientID,
		encoder.EncodeToString(codeVerifierHash[:]), codeChallengeMethod)

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
