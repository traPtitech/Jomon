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

type TagResponse struct {
	Tags []*TagOverview `json:"tags"`
}

func (h *Handlers) GetTags(c echo.Context) error {
	ctx := context.Background()
	tags, err := h.Repository.GetTags(ctx)
	if err != nil {
		return internalServerError(err)
	}

	res := []*TagOverview{}
	for _, tag := range tags {
		res = append(res, &TagOverview{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
			CreatedAt:   tag.CreatedAt,
			UpdatedAt:   tag.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, &TagResponse{res})
}

func (h *Handlers) PostTag(c echo.Context) error {
	var tag Tag
	if err := c.Bind(&tag); err != nil {
		return badRequest(err)
	}

	ctx := context.Background()
	created, err := h.Repository.CreateTag(ctx, tag.Name, tag.Description)
	if err != nil {
		return internalServerError(err)
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

func (h *Handlers) PutTag(c echo.Context) error {
	tagID, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		return badRequest(err)
	}
	if tagID == uuid.Nil {
		return badRequest(err)
	}
	var req Tag
	if err := c.Bind(&req); err != nil {
		return badRequest(err)
	}

	ctx := context.Background()
	tag, err := h.Repository.UpdateTag(ctx, tagID, req.Name, req.Description)
	if err != nil {
		return internalServerError(err)
	}

	res := &TagOverview{
		ID:          tag.ID,
		Name:        tag.Name,
		Description: tag.Description,
		CreatedAt:   tag.CreatedAt,
		UpdatedAt:   tag.UpdatedAt,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) DeleteTag(c echo.Context) error {
	tagID, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		return badRequest(err)
	}
	if tagID == uuid.Nil {
		return badRequest(err)
	}

	ctx := context.Background()
	err = h.Repository.DeleteTag(ctx, tagID)
	if err != nil {
		return internalServerError(err)
	}

	return c.NoContent(http.StatusOK)
}
