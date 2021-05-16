package router

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math/rand"
	"net/http"

	"github.com/dvsekhvalnov/jose2go/base64url"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/ent"
)

const (
	sessionDuration        = 24 * 60 * 60 * 7
	sessionKey             = "sessions"
	sessionAccessTokenKey  = "access_token"
	sessionCodeVerifierKey = "code_verifier"
	sessionRefreshTokenKey = "refresh_token"
	sessionUserKey         = "user"
)

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type PKCEParams struct {
	CodeChallenge       string `json:"code_challenge"`
	CodeChallengeMethod string `json:"code_challenge_method"`
	ClientID            string `json:"client_id"`
	ResponseType        string `json:"response_type"`
}

func (h Handlers) AuthUser(c echo.Context) (echo.Context, error) {
	sess, err := session.Get(sessionKey, c)
	if err != nil {
		return nil, c.NoContent(http.StatusInternalServerError)
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionDuration,
		HttpOnly: true,
	}

	accTok, ok := sess.Values[sessionAccessTokenKey].(string)
	if !ok || accTok == "" {
		return nil, c.NoContent(http.StatusUnauthorized)
	}
	c.Set(contextAccessTokenKey, accTok)

	user, ok := sess.Values[sessionUserKey].(ent.User)
	if !ok {
		user, err = h.Service.Users.GetMyUser(accTok)
		sess.Values[sessionUserKey] = user
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return nil, c.NoContent(http.StatusInternalServerError)
		}

		if err != nil {
			return nil, c.NoContent(http.StatusInternalServerError)
		}
	}

	admins, err := h.Service.Administrators.GetAdministratorList()
	if err != nil {
		return nil, c.NoContent(http.StatusInternalServerError)
	}
	user.GiveIsUserAdmin(admins)

	c.Set(contextUserKey, user)

	return c, nil
}

func (h Handlers) AuthCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if len(code) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	sess, err := session.Get(sessionKey, c)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionDuration,
		HttpOnly: true,
	}

	codeVerifier, ok := sess.Values[sessionCodeVerifierKey].(string)
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	res, err := h.Service.TraQAuth.GetAccessToken(code, codeVerifier)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	sess.Values[sessionAccessTokenKey] = res.AccessToken
	sess.Values[sessionRefreshTokenKey] = res.RefreshToken

	user, err := h.Service.Users.GetMyUser(res.AccessToken)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	sess.Values[sessionUserKey] = user

	if err = sess.Save(c.Request(), c.Response()); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h Handlers) GeneratePKCE(c echo.Context) error {
	sess, err := session.Get(sessionKey, c)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionDuration,
		HttpOnly: true,
	}

	bytesCodeVerifier := generateCodeVerifier()
	codeVerifier := string(bytesCodeVerifier)
	sess.Values[sessionCodeVerifierKey] = codeVerifier
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	params := PKCEParams{
		CodeChallenge:       getCodeChallenge(bytesCodeVerifier),
		CodeChallengeMethod: "S256",
		ClientID:            h.Service.GetClientId(),
		ResponseType:        "code",
	}

	return c.JSON(http.StatusOK, params)
}

var src cryptoSource
var randSrc = rand.New(src)

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

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}
