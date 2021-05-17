package router

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Transaction struct {
	Amount int          `json:"amount"`
	Target string       `json:"target"`
	Tags   []*uuid.UUID `json:"tags"`
	Group  *uuid.UUID   `json:"group"`
}

func (h *Handlers) GetTransactions(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) PostTransaction(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) GetTransaction(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) PutTransaction(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}
