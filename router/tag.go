package router

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/ent/tag"
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

type TagDetail struct {
	ID           uuid.UUID   `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Transactions []uuid.UUID `json:"transactions"`
	Requests     []uuid.UUID `json:"requests"`
}

func (h *Handlers) GetTags(c echo.Context) error {
	ctx := context.Background()
	tags, err := h.EntCli.Tag.
		Query().
		All(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
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

	return c.JSON(http.StatusOK, res)
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
	tagID, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if tagID == uuid.Nil {
		return c.NoContent(http.StatusBadRequest)
	}

	ctx := context.Background()
	tag, err := h.EntCli.Tag.
		Query().
		Where(tag.IDEQ(tagID)).
		Only(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
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

func (h *Handlers) PutTag(c echo.Context) error {
	tagID, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if tagID == uuid.Nil {
		return c.NoContent(http.StatusBadRequest)
	}
	var req Tag
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	ctx := context.Background()
	tag, err := h.EntCli.Tag.
		UpdateOneID(tagID).
		SetName(req.Name).
		SetDescription(req.Description).
		Save(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	transactions, err := tag.
		QueryTransaction().
		All(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	transactionIDs := []uuid.UUID{}
	for _, transaction := range transactions {
		transactionIDs = append(transactionIDs, transaction.ID)
	}

	requests, err := tag.
		QueryRequest().
		All(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	requestIDs := []uuid.UUID{}
	for _, request := range requests {
		requestIDs = append(requestIDs, request.ID)
	}

	res := &TagDetail{
		ID:           tag.ID,
		Name:         tag.Name,
		Description:  tag.Description,
		CreatedAt:    tag.CreatedAt,
		UpdatedAt:    tag.UpdatedAt,
		Transactions: transactionIDs,
		Requests:     requestIDs,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) DeleteTag(c echo.Context) error {
	tagID, err := uuid.Parse(c.Param("tagID"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if tagID == uuid.Nil {
		return c.NoContent(http.StatusBadRequest)
	}

	ctx := context.Background()
	tag, err := h.EntCli.Tag.
		Query().
		Where(tag.IDEQ(tagID)).
		Only(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	err = h.EntCli.Tag.
		DeleteOne(tag).
		Exec(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
	// TODO: Implement
}
