package router

import (
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

func (h *Handlers) PostTag(c echo.Context) error {
	// TODO: Implement
	/***** TEMPORARY IMPLEMENTED! *****/
	var tag Tag
	if err := c.Bind(&tag); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	res := TagOverview{
		ID:          uuid.New(),
		Name:        tag.Name,
		Description: tag.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return c.JSON(http.StatusOK, res)
}
