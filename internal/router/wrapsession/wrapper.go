package wrapsession

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"
)

// NOTE: ここの値は外から変更する必要がないので非公開にしている

const (
	stateKey        = "state"
	codeVerifierKey = "code_verifier"
	nonceKey        = "nonce"
	refererKey      = "referer"
	idTokenKey      = "id_token"
)

// TODO: この設定は外から与えられるようにしたい
var defaultSessionOptions = &sessions.Options{
	Path:     "/",
	MaxAge:   24 * 60 * 60 * 7,
	HttpOnly: true,
}

// *W は `*session.Session` のラッパーです.
//
// セッション操作として, 2つの値の操作が提供されています.
//
//   - ログインしているユーザーのID: `(*W).GetUserID`, `(*W).SetUserID`
//   - OAuth2.0 Authorization Code Flowで使用される `code_verifier`:
//     `(*W).GetCodeVerifier`, `(*W).SetCodeVerifier`
//   - OAuth2.0 Authorization Code Flowで使用される `state`:
//     `(*W).GetState`, `(*W).SetState`
type W struct {
	inner   *sessions.Session
	changed bool
}

func newW(sess *sessions.Session) *W {
	return &W{
		inner:   sess,
		changed: false,
	}
}

func (w *W) getValue(key string) interface{} {
	return w.inner.Values[key]
}

func (w *W) setValue(key string, value interface{}) {
	w.inner.Values[key] = value
	w.changed = true
}

func (w *W) drop(c *echo.Context) error {
	if !w.changed {
		return nil
	}
	w.inner.Options = defaultSessionOptions
	if err := w.inner.Save(c.Request(), c.Response()); err != nil {
		return newSaveSessionError(err)
	}
	return nil
}

func (w *W) GetState() (string, bool) {
	v, ok := w.getValue(stateKey).(string)
	return v, ok
}

func (w *W) SetState(value string) {
	w.setValue(stateKey, value)
}

func (w *W) GetCodeVerifier() (string, bool) {
	v, ok := w.getValue(codeVerifierKey).(string)
	return v, ok
}

func (w *W) SetCodeVerifier(value string) {
	w.setValue(codeVerifierKey, value)
}

func (w *W) GetNonce() (string, bool) {
	v, ok := w.getValue(nonceKey).(string)
	return v, ok
}

func (w *W) SetNonce(value string) {
	w.setValue(nonceKey, value)
}

func (w *W) GetReferer() (string, bool) {
	v, ok := w.getValue(refererKey).(string)
	return v, ok
}

func (w *W) SetReferer(value string) {
	w.setValue(refererKey, value)
}

func (w *W) GetIDToken() (string, bool) {
	v, ok := w.getValue(idTokenKey).(string)
	return v, ok
}

func (w *W) SetIDToken(value string) {
	w.setValue(idTokenKey, value)
}
