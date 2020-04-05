package router

import (
	"crypto/sha256"
	"github.com/dvsekhvalnov/jose2go/base64url"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
	"time"
)

const sessionDuration = 24 * 60 * 60 * 1000

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type PKCEParams struct {
	CodeChallenge       string
	CodeChallengeMethod string
	ClientID            string
	ResponseType        string
}

func (s Service) AuthCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if len(code) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	codeVerifier := sess.Values["codeVerifier"].(string)
	res, err := s.TraQAuth.GetAccessToken(code, codeVerifier)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionDuration,
		HttpOnly: true,
	}
	sess.Values["accessToken"] = res.AccessToken
	sess.Values["refreshToken"] = res.RefreshToken

	user, err := s.Users.GetMyUser(res.AccessToken)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	sess.Values["user"] = user

	if err = sess.Save(c.Request(), c.Response()); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (s Service) GeneratePKCE(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	bytesCodeVerifier := generateCodeVerifier()
	codeVerifier := string(bytesCodeVerifier)
	sess.Values["code_verifier"] = codeVerifier
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	params := PKCEParams{
		CodeChallenge:       getCodeChallenge(bytesCodeVerifier),
		CodeChallengeMethod: "S256",
		ClientID:            s.TraQAuth.GetClientId(),
		ResponseType:        "code",
	}

	return c.JSON(http.StatusOK, params)
}

var randSrc = rand.NewSource(time.Now().UnixNano())

const (
	// omit `.` and `~` to improve performance
	letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func generateCodeVerifier() []byte {
	bytesCodeVerifier := make([]byte, 128)
	cache, remain := randSrc.Int63(), letterIdxMax
	for i := 0; i < 128; i++ {
		if remain == 0 {
			cache, remain = randSrc.Int63(), letterIdxMax
		}
		idx := int(cache & letterIdxMask)
		bytesCodeVerifier[i] = letters[idx]
		cache >>= letterIdxBits
		remain--
	}

	return bytesCodeVerifier
}

func getCodeChallenge(cv []byte) string {
	bytesCodeChallenge := sha256.Sum256(cv)
	return base64url.Encode(bytesCodeChallenge[:])
}
