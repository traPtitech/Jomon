package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func PutStatus(c echo.Context) error {
	// some program
	return c.String(http.StatusOK, "PutStatus")
}

func PutRepaidStatus(c echo.Context) error {
	// some program
	return c.String(http.StatusOK, "PutRepaidStatus")
}
