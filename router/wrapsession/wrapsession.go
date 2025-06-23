// FIXME: Markdownのつもりでドキュメントを書いているが、違う

// wrapsession は `*session.Session` のラッパーを提供します.
//
//	var c echo.Context, sessionName string
//	v, err := wrapsession.WithSession(c, sessionName, func (w *wrapsession.W) (T, error) {
//		return doSomething(w)
//	})
//
// 詳細に関しては `WithSession` を参照してください.
package wrapsession

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// `WithSession` はセッションを使用するスコープを提供します.
//
// 例えば, ログインしているユーザーのIDをセッションから取り出す操作はこのようになります.
//
//	userID, err := wrapsession.WithSession(
//		c, sessionName, func(w *wrapsession.W) (uuid.UUID, error) {
//			id, ok := w.GetUserID()
//			if !ok || id == uuid.Nil {
//				return uuid.Nil, errors.New("unauthenticated")
//			}
//			return id, nil
//		})
//
// 操作の中で使用する `*W` 型の値は, その `WithSession` 呼び出しで提供されるスコープ内でのみ有効です.
// スコープの外側で `*W` の値を使用すると, 思わぬバグに繋がる可能性があります.
//
// 操作中に `(*W).SetUserID` などでセッションへの書き込みが発生していた場合,
// そのセッションはHTTPレスポンスへ反映されます.
// 具体的には, `(*sessions.Session).Save` が呼び出されます.
//
// 以下の場合にエラーが発生します.
//
//   - `session.Get` がエラーを返した場合. エラーの型は *GetSessionError となります.
//   - 与えられた操作がエラーを返した場合. エラーの型は特に変わりません.
//   - `(*sessions.Session).Save` がエラーを返した場合. エラーの型は *SaveSessionError となります.
//
// エラーの型が区別されるため, type switchを用いて詳細なハンドリングが可能です.
//
// nolint:ireturn
func WithSession[T any](c echo.Context, sessionName string, op func(w *W) (T, error)) (T, error) {
	var res T
	sess, err := session.Get(sessionName, c)
	if err != nil {
		return res, newGetSessionError(err)
	}
	w := newW(sess)
	res, err = op(w)
	if err != nil {
		return res, err
	}
	if err := w.drop(c); err != nil {
		return res, err
	}
	return res, nil
}
