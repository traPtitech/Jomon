package router

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TagOverview struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (h *Handlers) GetTags(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) PostTag(c echo.Context) error {
	var tag Tag
	if err := c.Bind(&tag); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	ctx := context.Background()
	created, err := h.EntCli.Tag.
		Create().
		SetName(tag.Name).
		SetDescription(tag.Description).
		Save(ctx)

	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	res := TagOverview{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) GetTag(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) PutTag(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) DeleteTag(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}
