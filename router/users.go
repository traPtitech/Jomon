package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetUsers(c echo.Context) error {
	// some program
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, "GetUsers")
}

func GetMyUser(c echo.Context) error {
	// some program
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, "GetMyUser")
}

func GetAdminUsers(c echo.Context) error {
	// some program
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, "GetAdminUsers")
}

func PutAdminUsers(c echo.Context) error {
	// some program
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, "PutAdminUsers")
}
