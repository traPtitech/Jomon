package router

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/traPtitech/Jomon/internal/router/wrapsession"
	"github.com/traPtitech/Jomon/internal/service"
)

func HTTPErrorHandlerInner(err error) *echo.HTTPError {
	if err == nil {
		return nil
	}
	// FIXME(error-handling): type switchで最上位のエラーに対してhandlingする
	if he := new(echo.HTTPError); errors.As(err, &he) {
		return he
	}
	var herr error
	if e := new(service.BadInputError); errors.As(err, &e) {
		herr = echo.NewHTTPError(http.StatusBadRequest, e.Message).Wrap(e)
	} else if e := new(service.NotFoundError); errors.As(err, &e) {
		herr = echo.NewHTTPError(http.StatusNotFound, e.Message).Wrap(e)
	} else if e := new(service.ForbiddenError); errors.As(err, &e) {
		herr = echo.NewHTTPError(http.StatusForbidden, e.Message).Wrap(e)
	} else if e := new(service.UnauthenticatedError); errors.As(err, &e) {
		herr = echo.NewHTTPError(http.StatusUnauthorized, e.Message).Wrap(e)
	} else if e := new(service.UnexpectedError); errors.As(err, &e) {
		herr = echo.ErrInternalServerError.Wrap(e)
	} else if e := new(wrapsession.GetSessionError); errors.As(err, &e) {
		herr = echo.ErrInternalServerError.Wrap(e)
	} else if e := new(wrapsession.SaveSessionError); errors.As(err, &e) {
		herr = echo.ErrInternalServerError.Wrap(e)
	} else {
		herr = echo.ErrInternalServerError.Wrap(err)
	}
	he := new(echo.HTTPError)
	if !errors.As(herr, &he) {
		// ここには来ないはず
		code := http.StatusInternalServerError
		return echo.NewHTTPError(code, http.StatusText(code))
	}
	return he
}
