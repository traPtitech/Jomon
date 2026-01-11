package router

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

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
		return echo.NewHTTPError(http.StatusBadRequest, e.Message).SetInternal(e)
	}
	if e := new(service.NotFoundError); errors.As(err, &e) {
		return echo.NewHTTPError(http.StatusNotFound, e.Message).SetInternal(e)
	}
	if e := new(service.ForbiddenError); errors.As(err, &e) {
		return echo.NewHTTPError(http.StatusForbidden, e.Message).SetInternal(e)
	}
	if e := new(service.UnauthenticatedError); errors.As(err, &e) {
		return echo.NewHTTPError(http.StatusUnauthorized, e.Message).SetInternal(e)
	}
	if e := new(service.UnexpectedError); errors.As(err, &e) {
		return echo.ErrInternalServerError.WithInternal(e)
	}
	return echo.ErrInternalServerError.WithInternal(err)
}
