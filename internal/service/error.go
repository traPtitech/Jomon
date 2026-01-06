package service

// BadInputError は不正な入力があった場合に返されるエラーです.
// ユーザーに表示しても問題ないメッセージを含みます.
// 内部エラーを含むことができ, Unwrap メソッドで取得できます. この内部エラーはユーザーに表示されるべきではありません.
//
// Example:
//
//	return service.NewBadInputError("invalid input").WithInternal(err)
//
// 対応するHTTPステータスコード: [400 Bad Request]
//
// [400 Bad Request]: https://developer.mozilla.org/ja/docs/Web/HTTP/Reference/Status/400
type BadInputError struct {
	Message  string
	internal error
}

func NewBadInputError(message string) *BadInputError {
	return &BadInputError{
		Message:  message,
		internal: nil,
	}
}

func (e *BadInputError) Error() string {
	return e.Message
}

func (e *BadInputError) Unwrap() error {
	return e.internal
}

func (e *BadInputError) WithInternal(err error) *BadInputError {
	return &BadInputError{
		Message:  e.Message,
		internal: err,
	}
}

// NotFoundError はリソースが見つからなかった場合に返されるエラーです.
// ユーザーに表示しても問題ないメッセージを含みます.
// 内部エラーを含むことができ, Unwrap メソッドで取得できます. この内部エラーはユーザーに表示されるべきではありません.
//
// Example:
//
//	return service.NewNotFoundError("resource not found").WithInternal(err)
//
// 対応するHTTPステータスコード: [404 Not Found]
//
// [404 Not Found]: https://developer.mozilla.org/ja/docs/Web/HTTP/Reference/Status/404
type NotFoundError struct {
	Message  string
	internal error
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Message:  message,
		internal: nil,
	}
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *NotFoundError) Unwrap() error {
	return e.internal
}

func (e *NotFoundError) WithInternal(err error) *NotFoundError {
	return &NotFoundError{
		Message:  e.Message,
		internal: err,
	}
}

// ForbiddenError はリソースへの操作が禁止されている場合に返されるエラーです.
// ユーザーに表示しても問題ないメッセージを含みます.
// 内部エラーを含むことができ, Unwrap メソッドで取得できます. この内部エラーはユーザーに表示されるべきではありません.
//
// Example:
//
//	return service.NewForbiddenError("access forbidden").WithInternal(err)
//
// **note**: セキュリティ上の理由から, リソースを取得できなかった場合には NotFoundError を返すことを推奨します.
//
// 対応するHTTPステータスコード: [403 Forbidden]
//
// [403 Forbidden]: https://developer.mozilla.org/ja/docs/Web/HTTP/Reference/Status/403
type ForbiddenError struct {
	Message  string
	internal error
}

func NewForbiddenError(message string) *ForbiddenError {
	return &ForbiddenError{
		Message:  message,
		internal: nil,
	}
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func (e *ForbiddenError) Unwrap() error {
	return e.internal
}

func (e *ForbiddenError) WithInternal(err error) *ForbiddenError {
	return &ForbiddenError{
		Message:  e.Message,
		internal: err,
	}
}

// UnauthenticatedError は認証が必要な操作で認証されていない場合に返されるエラーです.
// ユーザーに表示しても問題ないメッセージを含みます.
// 内部エラーを含むことができ, Unwrap メソッドで取得できます. この内部エラーはユーザーに表示されるべきではありません.
//
// Example:
//
//	return service.NewUnauthenticatedError("authentication required").WithInternal(err)
//
// 対応するHTTPステータスコード: [401 Unauthorized]
//
// [401 Unauthorized]: https://developer.mozilla.org/ja/docs/Web/HTTP/Reference/Status/401
type UnauthenticatedError struct {
	Message  string
	internal error
}

func NewUnauthenticatedError(message string) *UnauthenticatedError {
	return &UnauthenticatedError{
		Message:  message,
		internal: nil,
	}
}

func (e *UnauthenticatedError) Error() string {
	return e.Message
}

func (e *UnauthenticatedError) Unwrap() error {
	return e.internal
}

func (e *UnauthenticatedError) WithInternal(err error) *UnauthenticatedError {
	return &UnauthenticatedError{
		Message:  e.Message,
		internal: err,
	}
}

// UnexpectedError は予期しないエラーが発生した場合に返されるエラーです.
// ユーザーに表示しても問題ない一般的なメッセージを含みます.
// 内部エラーを含むことができ, Unwrap メソッドで取得できます. この内部エラーはユーザーに表示されるべきではありません.
//
// Example:
//
//	return service.NewUnexpectedError(err)
//
// 対応するHTTPステータスコード: [500 Internal Server Error]
//
// [500 Internal Server Error]: https://developer.mozilla.org/ja/docs/Web/HTTP/Reference/Status/500
type UnexpectedError struct {
	internal error
}

func NewUnexpectedError(err error) *UnexpectedError {
	return &UnexpectedError{
		internal: err,
	}
}

func (e *UnexpectedError) Error() string {
	return "an unexpected error occurred"
}

func (e *UnexpectedError) Unwrap() error {
	return e.internal
}

// interface guards
var (
	_ error = (*BadInputError)(nil)
	_ error = (*NotFoundError)(nil)
	_ error = (*ForbiddenError)(nil)
	_ error = (*UnauthenticatedError)(nil)
	_ error = (*UnexpectedError)(nil)
)
