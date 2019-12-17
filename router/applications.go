package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetApplications(c echo.Context) error {
	return c.String(http.StatusOK, "GetApplications")
}

func PostApplications(c echo.Context) error {
	return c.String(http.StatusOK, "PostApplications")
}

func PatchApplications(c echo.Context) error {
	return c.String(http.StatusOK, "PatchApplications")
}
