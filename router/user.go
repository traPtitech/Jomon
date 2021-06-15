package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)


func (h *Handlers) GetUsers(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) PutUsers(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) GetMe(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}
