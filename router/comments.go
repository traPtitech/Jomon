package router

import (
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
)

type PostCommentRequest struct {
	Comment string `json:"comment"`
}

func (s *Service) PostComments(c echo.Context) error {
	applicationId := uuid.FromStringOrNil(c.Param("applicationId"))
	if applicationId == uuid.Nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if _, err := s.Applications.GetApplication(applicationId, false); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.NoContent(http.StatusNotFound)
		} else {
			return c.NoContent(http.StatusBadRequest)
		}
	}

	var req PostCommentRequest
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if req.Comment == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	user, ok := c.Get("user").(model.User)
	if !ok || user.TrapId == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	comment, err := s.Comments.CreateComment(applicationId, req.Comment, user.TrapId)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, comment)
}

func (s *Service) PutComments(c echo.Context) error {
	applicationId, err := uuid.FromString(c.Param("applicationId"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	comment, err := s.Comments.GetComment(applicationId, commentId)
	if err == gorm.ErrRecordNotFound {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	user, ok := c.Get("user").(model.User)
	if !ok || user.TrapId == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	if comment.UserTrapID.TrapId != user.TrapId {
		return c.NoContent(http.StatusForbidden)
	}

	var req PostCommentRequest
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if req.Comment == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	comment, err = s.Comments.PutComment(applicationId, commentId, req.Comment)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, comment)
}

func (s *Service) DeleteComments(c echo.Context) error {
	applicationId, err := uuid.FromString(c.Param("applicationId"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	comment, err := s.Comments.GetComment(applicationId, commentId)
	if err == gorm.ErrRecordNotFound {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	user, ok := c.Get("user").(model.User)
	if !ok || user.TrapId == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	if comment.UserTrapID.TrapId != user.TrapId {
		return c.NoContent(http.StatusForbidden)
	}

	if err = s.Comments.DeleteComment(applicationId, commentId); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
