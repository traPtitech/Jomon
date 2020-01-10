package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetImages(c echo.Context) error {
	// some program
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, "GetImages")
}

func DeleteImages(c echo.Context) error {
	// some program
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, "DeleteImages")
}
