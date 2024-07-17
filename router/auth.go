package router

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand/v2"
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
	"go.uber.org/zap"
)

const (
	sessionDuration        = 24 * 60 * 60 * 7
	sessionCodeVerifierKey = "code_verifier"
)

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (h Handlers) AuthCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if len(code) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "code is required")
	}

	sess, err := session.Get(h.SessionName, c)
	if err != nil {
		h.Logger.Error("failed to get session", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionDuration,
		HttpOnly: true,
	}

	codeVerifier, ok := sess.Values[sessionCodeVerifierKey].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("code_verifier is not found in session"))
	}

	res, err := service.RequestAccessToken(code, codeVerifier)
	if err != nil {
		h.Logger.Error("failed to get access token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	u, err := service.FetchTraQUserInfo(res.AccessToken)
	if err != nil {
		h.Logger.Error("failed to fetch traQ user info", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var modelUser *model.User
	modelUser, err = h.Repository.GetUserByName(c.Request().Context(), u.Name)
	if err != nil {
		if ent.IsNotFound(err) {
			modelUser, err = h.Repository.CreateUser(c.Request().Context(), u.Name, u.DisplayName, false)
			if err != nil {
				h.Logger.Error("failed to create user", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		} else {
			h.Logger.Error("failed to get user by name", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	user := &User{
		ID:          modelUser.ID,
		Name:        modelUser.Name,
		DisplayName: modelUser.DisplayName,
		Admin:       modelUser.Admin,
	}

	sess.Values[sessionUserKey] = user

	if err = sess.Save(c.Request(), c.Response()); err != nil {
		h.Logger.Error("failed to save session", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h Handlers) GeneratePKCE(c echo.Context) error {
	sess, err := session.Get(h.SessionName, c)
	if err != nil {
		h.Logger.Error("failed to get session", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionDuration,
		HttpOnly: true,
	}

	codeVerifier := randAlphabetAndNumberString(43)
	sess.Values[sessionCodeVerifierKey] = codeVerifier

	codeVerifierHash := sha256.Sum256([]byte(codeVerifier))
	encoder := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_").WithPadding(base64.NoPadding)

	codeChallengeMethod := "S256"

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		h.Logger.Error("failed to save session", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("%s/oauth2/authorize?response_type=code&client_id=%s&code_challenge=%s&code_challenge_method=%s", service.TraQBaseURL, service.JomonClientID, encoder.EncodeToString(codeVerifierHash[:]), codeChallengeMethod))
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
	randSrc := randSrcPool.Get().(*rand.Rand)
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
