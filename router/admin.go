package router

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/ent"
)

func (h *Handlers) GetAdmins(c echo.Context) error {
	ctx := c.Request().Context()
	admins, err := h.Repository.GetAdmins(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := []*uuid.UUID{}
	for _, admin := range admins {
		res = append(res, &admin.ID)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) PostAdmins(c echo.Context) error {
	var admin []uuid.UUID
	if err := c.Bind(&admin); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	err := h.Repository.AddAdmins(ctx, admin)
	if err != nil {
		if ent.IsConstraintError(err) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handlers) DeleteAdmins(c echo.Context) error {
	var admin []uuid.UUID
	if err := c.Bind(&admin); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	err := h.Repository.DeleteAdmins(ctx, admin)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
