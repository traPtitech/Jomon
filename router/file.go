package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) PostFile(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) GetFile(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) DeleteFile(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}
