package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Service) GetImages(c echo.Context) error {
	// some program
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, "GetImages")
}

func (s *Service) DeleteImages(c echo.Context) error {
	// some program
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, "DeleteImages")
}
