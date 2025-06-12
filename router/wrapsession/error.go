package wrapsession

import "fmt"

type GetSessionError struct {
	inner error
}

func newGetSessionError(err error) *GetSessionError {
	return &GetSessionError{
		inner: err,
	}
}

func (e *GetSessionError) Error() string {
	return fmt.Sprintf("failed to get session from cookie: %s", e.inner.Error())
}

func (e *GetSessionError) Unwrap() error {
	return e.inner
}

type SaveSessionError struct {
	inner error
}

func newSaveSessionError(err error) *SaveSessionError {
	return &SaveSessionError{
		inner: err,
	}
}

func (e *SaveSessionError) Error() string {
	return fmt.Sprintf("failed to save session to cookie: %s", e.inner.Error())
}

func (e *SaveSessionError) Unwrap() error {
	return e.inner
}

var (
	_ error = (*GetSessionError)(nil)
	_ error = (*SaveSessionError)(nil)
)
