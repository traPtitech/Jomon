package router

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
)

type Request struct {
	CreatedBy uuid.UUID    `json:"created_by"`
	Amount    int          `json:"amount"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Tags      []*uuid.UUID `json:"tags"`
	Group     *uuid.UUID   `json:"group"`
}

type RequestResponse struct {
	ID        uuid.UUID      `json:"id"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedBy uuid.UUID      `json:"created_by"`
	Amount    int            `json:"amount"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Tags      []*TagOverview `json:"tags"`
	Group     *GroupOverview `json:"group"`
}

type Comment struct {
	Comment string `json:"comment"`
}
type CommentDetail struct {
	ID        uuid.UUID `json:"id"`
	User      uuid.UUID `json:"user"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *Handlers) GetRequests(c echo.Context) error {
	ctx := context.Background()
	modelrequests, err := h.Repository.GetRequests(ctx, model.RequestQuery{})
	if err != nil {
		return internalServerError(err)
	}

	var tags []*TagOverview
	var requests []*RequestResponse
	for _, request := range modelrequests {
		for _, tag := range request.Tags {
			tags = append(tags, &TagOverview{
				ID:          tag.ID,
				Name:        tag.Name,
				Description: tag.Description,
				CreatedAt:   tag.CreatedAt,
				UpdatedAt:   tag.UpdatedAt,
			})
		}

		var resgroup *GroupOverview
		if request.Group != nil {
			resgroup = &GroupOverview{
				ID:          request.Group.ID,
				Name:        request.Group.Name,
				Description: request.Group.Description,
				Budget:      request.Group.Budget,
				CreatedAt:   request.Group.CreatedAt,
				UpdatedAt:   request.Group.UpdatedAt,
			}
		}

		res := &RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Amount:    request.Amount,
			Title:     request.Title,
			Content:   request.Content,
			Tags:      tags,
			Group:     resgroup,
		}
		requests = append(requests, res)
	}

	return c.JSON(http.StatusOK, requests)
}

func (h *Handlers) PostRequest(c echo.Context) error {
	var req Request
	var err error
	if err := c.Bind(&req); err != nil {
		return badRequest(err)
	}
	var tags []*model.Tag
	for _, tagID := range req.Tags {
		ctx := context.Background()
		tag, err := h.Repository.GetTag(ctx, *tagID)
		if err != nil {
			return internalServerError(err)
		}
		tags = append(tags, tag)
	}
	var group *model.Group
	if req.Group != nil {
		ctx := context.Background()
		group, err = h.Repository.GetGroup(ctx, *req.Group)
		if err != nil {
			return internalServerError(err)
		}
	}
	ctx := context.Background()
	request, err := h.Repository.CreateRequest(ctx, req.Amount, req.Title, req.Content, tags, group, req.CreatedBy)
	if err != nil {
		return internalServerError(err)
	}
	var resgroup *GroupOverview
	if group != nil {
		resgroup = &GroupOverview{
			ID:          request.Group.ID,
			Name:        request.Group.Name,
			Description: request.Group.Description,
			Budget:      request.Group.Budget,
			CreatedAt:   request.Group.CreatedAt,
			UpdatedAt:   request.Group.UpdatedAt,
		}
	}
	var restags []*TagOverview
	for _, tag := range request.Tags {
		restags = append(restags, &TagOverview{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
			CreatedAt:   tag.CreatedAt,
			UpdatedAt:   tag.UpdatedAt,
		})
	}
	res := &RequestResponse{
		ID:        request.ID,
		Status:    request.Status,
		CreatedAt: request.CreatedAt,
		UpdatedAt: request.UpdatedAt,
		CreatedBy: request.CreatedBy,
		Amount:    request.Amount,
		Title:     request.Title,
		Content:   request.Content,
		Tags:      restags,
		Group:     resgroup,
	}
	return c.JSON(http.StatusOK, res)
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
	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		return badRequest(err)
	}
	if requestID == uuid.Nil {
		return badRequest(err)
	}

	var req Comment
	if err := c.Bind(&req); err != nil {
		return badRequest(err)
	}

	user, ok := c.Get(contextUserKey).(model.User)
	if !ok || user.ID == uuid.Nil {
		return unauthorized(err)
	}
	ctx := context.Background()
	comment, err := h.Repository.CreateComment(ctx, req.Comment, requestID, user.ID)
	if err != nil {
		return internalServerError(err)
	}
	res := &CommentDetail{
		ID:        comment.ID,
		User:      comment.User,
		Comment:   comment.Comment,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) PutComment(c echo.Context) error {
	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		return badRequest(err)
	}
	if requestID == uuid.Nil {
		return badRequest(err)
	}
	commentID, err := uuid.Parse(c.Param("commentID"))
	if err != nil {
		return badRequest(err)
	}
	if commentID == uuid.Nil {
		return badRequest(err)
	}

	var req Comment
	if err := c.Bind(&req); err != nil {
		return badRequest(err)
	}

	ctx := context.Background()
	comment, err := h.Repository.UpdateComment(ctx, req.Comment, requestID, commentID)
	if err != nil {
		return internalServerError(err)
	}
	res := &CommentDetail{
		ID:        comment.ID,
		User:      comment.User,
		Comment:   comment.Comment,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) DeleteComment(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) PutStatus(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}
