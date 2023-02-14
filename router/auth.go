package router

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
)

const (
	sessionDuration        = 24 * 60 * 60 * 7
	sessionKey             = "sessions"
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
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	u, err := service.FetchTraQUserInfo(res.AccessToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var modelUser *model.User
	modelUser, err = h.Repository.GetUserByName(c.Request().Context(), u.Name)
	if err != nil {
		if ent.IsNotFound(err) {
			modelUser, err = h.Repository.CreateUser(c.Request().Context(), u.Name, u.DisplayName, false)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	user := &User{
		ID:          modelUser.ID,
		Name:        modelUser.Name,
		DisplayName: modelUser.DisplayName,
		Admin:       modelUser.Admin,
	}
	gob.Register(&User{})
	sess.Values[sessionUserKey] = user

	if err = sess.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h Handlers) GeneratePKCE(c echo.Context) error {
	sess, err := session.Get(h.SessionName, c)
	if err != nil {
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
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("%s/oauth2/authorize?response_type=code&client_id=%s&code_challenge=%s&code_challenge_method=%s", service.TraQBaseURL, service.JomonClientID, encoder.EncodeToString(codeVerifierHash[:]), codeChallengeMethod))
}

var randSrcPool = sync.Pool{
	New: func() interface{} {
		return rand.NewSource(time.Now().UnixNano())
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
	randSrc := randSrcPool.Get().(rand.Source)
	cache, remain := randSrc.Int63(), rs6LetterIdxMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), rs6LetterIdxMax
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
