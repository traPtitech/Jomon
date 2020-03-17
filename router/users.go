package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
)

func GetUsers(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	admins := []string{}
	users, err := model.GetUsers(token, admins, false)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, users)
}

func GetMyUser(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	admins := []string{}
	user, err := model.GetMyUser(token, admins)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, user)
}

func GetAdminUsers(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	admins := []string{}
	users, err := model.GetUsers(token, admins, true)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, users)
}

func PutAdminUsers(c echo.Context) error {
	// some program
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	return c.String(http.StatusOK, "PutAdminUsers")
}
