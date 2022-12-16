package router

import (
	"errors"
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
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *Handlers) GetTags(c echo.Context) error {
	ctx := c.Request().Context()
	tags, err := h.Repository.GetTags(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := []*TagOverview{}
	for _, tag := range tags {
		res = append(res, &TagOverview{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) PostTag(c echo.Context) error {
	var tag Tag
	if err := c.Bind(&tag); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	created, err := h.Repository.CreateTag(ctx, tag.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := TagOverview{
		ID:        created.ID,
		Name:      created.Name,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) PutTag(c echo.Context) error {
	tagID, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if tagID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid tag ID"))
	}
	var req Tag
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	tag, err := h.Repository.UpdateTag(ctx, tagID, req.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := &TagOverview{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) DeleteTag(c echo.Context) error {
	tagID, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if tagID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid tag ID"))
	}

	ctx := c.Request().Context()
	err = h.Repository.DeleteTag(ctx, tagID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
