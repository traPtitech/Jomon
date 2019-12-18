package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func PutStatus(c echo.Context) error {
	// some program
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, "PutStatus")
}

func PutRepaidStatus(c echo.Context) error {
	// some program
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, "PutRepaidStatus")
}
