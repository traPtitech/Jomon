package router

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/logging"
	"github.com/traPtitech/Jomon/model"
	"go.uber.org/zap"
)

type PostTagRequest struct {
	Name string `json:"name"`
}

type PutTagRequest = PostTagRequest

type TagResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h Handlers) GetTags(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	tags, err := h.Repository.GetTags(ctx)
	if err != nil {
		logger.Error("failed to get tags from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := lo.Map(tags, func(tag *model.Tag, _ int) *TagResponse {
		return &TagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}
	})

	return c.JSON(http.StatusOK, res)
}

func (h Handlers) PostTag(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	var tag PostTagRequest
	if err := c.Bind(&tag); err != nil {
		logger.Info("could not get tag from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	created, err := h.Repository.CreateTag(ctx, tag.Name)
	if err != nil {
		logger.Error("failed to create tag in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := TagResponse{
		ID:        created.ID,
		Name:      created.Name,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}

	return c.JSON(http.StatusOK, res)
}

func (h Handlers) PutTag(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	tagID, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		logger.Info("could not parse query parameter `tagID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if tagID == uuid.Nil {
		logger.Info("invalid tag ID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid tag ID"))
	}
	var req PutTagRequest
	if err := c.Bind(&req); err != nil {
		logger.Info("could not get tag from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tag, err := h.Repository.UpdateTag(ctx, tagID, req.Name)
	if err != nil {
		logger.Error("failed to update tag in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := &TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}

	return c.JSON(http.StatusOK, res)
}

func (h Handlers) DeleteTag(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	tagID, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		logger.Info("could not parse query parameter `tagID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if tagID == uuid.Nil {
		logger.Info("invalid tag ID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid tag ID"))
	}

	err = h.Repository.DeleteTag(ctx, tagID)
	if err != nil {
		logger.Error("failed to delete tag in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
