package router

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Admin struct {
	ID uuid.UUID `json:"id"`
}

func (h *Handlers) GetAdmins(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) PostAdmin(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) DeleteAdmin(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}
