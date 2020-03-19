package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func PutStates(c echo.Context) error {
	// some program
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, "PutStates")
}

func PutRepaidStates(c echo.Context) error {
	// some program
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, "PutRepaidStates")
}
