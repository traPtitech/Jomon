package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
)

func GetApplications(c echo.Context) error {
	allapplications, err := model.GetApplications(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, allapplications)
	}
	return c.JSON(http.StatusOK, allapplications)
}

func PostApplications(c echo.Context) error {
	application, err := model.PostApplications(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, application)
	}
	return c.JSON(http.StatusOK, application)
}

func PatchApplications(c echo.Context) error {
	application, err := model.PatchApplications(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, application)
	}
	return c.JSON(http.StatusOK, application)
}
