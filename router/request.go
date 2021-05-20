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
	CreatedBy uuid.UUID `json:"created_by"`
	Amount    int       `json:"amount"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []*Tag    `json:"tags"`
	Group     *Group    `json:"group"`
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
	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if requestID == uuid.Nil {
		return c.NoContent(http.StatusBadRequest)
	}

	var req Comment
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	user, ok := c.Get(contextUserKey).(model.User)
	if !ok || user.ID == uuid.Nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	ctx := context.Background()
	comment, err := h.Repository.CreateComment(ctx, req.Comment, requestID, user.ID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
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
