package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetUsers(c echo.Context) error {
	// some program
	return c.String(http.StatusOK, "GetUsers")
}

func GetMyUser(c echo.Context) error {
	// some program
	return c.String(http.StatusOK, "GetMyUser")
}

func GetAdminUsers(c echo.Context) error {
	// some program
	return c.String(http.StatusOK, "GetAdminUsers")
}

func PutAdminUsers(c echo.Context) error {
	// some program
	return c.String(http.StatusOK, "PutAdminUsers")
}
