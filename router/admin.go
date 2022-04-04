package router

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/ent"
)

type Admin struct {
	ID uuid.UUID `json:"id"`
}

func (h *Handlers) GetAdmins(c echo.Context) error {
	ctx := context.Background()
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

func (h *Handlers) PostAdmin(c echo.Context) error {
	var admin Admin
	if err := c.Bind(&admin); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := context.Background()
	created, err := h.Repository.CreateAdmin(ctx, admin.ID)
	if err != nil {
		if ent.IsConstraintError(err) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := Admin{
		ID: created.ID,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) DeleteAdmin(c echo.Context) error {
	id, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := context.Background()
	err = h.Repository.DeleteAdmin(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
