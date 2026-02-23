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
	if he := new(echo.HTTPError); errors.As(err, &he) {
		return he
	}
	if e := new(service.BadInputError); errors.As(err, &e) {
		return echo.NewHTTPError(http.StatusBadRequest, e.Message).Wrap(e).(*echo.HTTPError)
	}
	if e := new(service.NotFoundError); errors.As(err, &e) {
		return echo.NewHTTPError(http.StatusNotFound, e.Message).Wrap(e).(*echo.HTTPError)
	}
	if e := new(service.ForbiddenError); errors.As(err, &e) {
		return echo.NewHTTPError(http.StatusForbidden, e.Message).Wrap(e).(*echo.HTTPError)
	}
	if e := new(service.UnauthenticatedError); errors.As(err, &e) {
		return echo.NewHTTPError(http.StatusUnauthorized, e.Message).Wrap(e).(*echo.HTTPError)
	}
	if e := new(service.UnexpectedError); errors.As(err, &e) {
		return echo.ErrInternalServerError.Wrap(e).(*echo.HTTPError)
	}
	if e := new(wrapsession.GetSessionError); errors.As(err, &e) {
		return echo.ErrInternalServerError.Wrap(e).(*echo.HTTPError)
	}
	if e := new(wrapsession.SaveSessionError); errors.As(err, &e) {
		return echo.ErrInternalServerError.Wrap(e).(*echo.HTTPError)
	}
	return echo.ErrInternalServerError.Wrap(err).(*echo.HTTPError)
}
