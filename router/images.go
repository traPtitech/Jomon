package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetImages(c echo.Context) error {
	return c.String(http.StatusOK, "GetImages")
}

func DeleteImages(c echo.Context) error {
	return c.String(http.StatusOK, "DeleteImages")
}
