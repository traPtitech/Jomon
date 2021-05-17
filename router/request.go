package router

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Request struct {
	CreatedBy uuid.UUID `json:"created_by"`
	Amount    int       `json:"amount"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []*Tag    `json:"tags"`
	Group     *Group    `json:"group"`
}

func (h *Handlers) GetRequests(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) PostRequest(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) GetRequest(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) PutRequest(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) PostComment(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) PutComment(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) DeleteComment(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) PutStatus(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}
