package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
)

type Request struct {
	CreatedBy uuid.UUID    `json:"created_by"`
	Amount    int          `json:"amount"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Tags      []*uuid.UUID `json:"tags"`
	Group     *uuid.UUID   `json:"group"`
}

type PutRequest struct {
	Amount  int          `json:"amount"`
	Title   string       `json:"title"`
	Content string       `json:"content"`
	Tags    []*uuid.UUID `json:"tags"`
	Group   *uuid.UUID   `json:"group"`
}

type RequestResponse struct {
	ID        uuid.UUID        `json:"id"`
	Status    string           `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	CreatedBy uuid.UUID        `json:"created_by"`
	Amount    int              `json:"amount"`
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	Tags      []*TagOverview   `json:"tags"`
	Group     *GroupOverview   `json:"group"`
	Comments  []*CommentDetail `json:"comments"`
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
	sort := c.QueryParam("sort")
	target := c.QueryParam("target")
	var year int
	var err error
	if c.QueryParam("year") != "" {
		year, err = strconv.Atoi(c.QueryParam("year"))
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}
	var since time.Time
	if c.QueryParam("since") != "" {
		since, err = service.StrToDate(c.QueryParam("since"))
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}
	var until time.Time
	if c.QueryParam("until") != "" {
		until, err = service.StrToDate(c.QueryParam("until"))
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}
	tag := c.QueryParam("tag")
	group := c.QueryParam("group")
	query := model.RequestQuery{
		Sort:   &sort,
		Target: &target,
		Year:   &year,
		Since:  &since,
		Until:  &until,
		Tag:    &tag,
		Group:  &group,
	}

	modelrequests, err := h.Repository.GetRequests(ctx, query)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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
	if err = c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	var tags []*model.Tag
	for _, tagID := range req.Tags {
		ctx := context.Background()
		tag, err := h.Repository.GetTag(ctx, *tagID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		tags = append(tags, tag)
	}
	var group *model.Group
	if req.Group != nil {
		ctx := context.Background()
		group, err = h.Repository.GetGroup(ctx, *req.Group)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	ctx := context.Background()
	request, err := h.Repository.CreateRequest(ctx, req.Amount, req.Title, req.Content, tags, group, req.CreatedBy)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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
	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	ctx := context.Background()
	request, err := h.Repository.GetRequest(ctx, requestID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	modelcomments, err := h.Repository.GetComments(ctx, requestID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	var comments []*CommentDetail
	for _, modelcomment := range modelcomments {
		comment := &CommentDetail{
			ID:        modelcomment.ID,
			User:      modelcomment.User,
			Comment:   modelcomment.Comment,
			CreatedAt: modelcomment.CreatedAt,
			UpdatedAt: modelcomment.UpdatedAt,
		}
		comments = append(comments, comment)
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
		Comments:  comments,
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) PutRequest(c echo.Context) error {
	var req PutRequest
	var err error
	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if requestID == uuid.Nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	var tags []*model.Tag
	for _, tagID := range req.Tags {
		ctx := context.Background()
		tag, err := h.Repository.GetTag(ctx, *tagID)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		tags = append(tags, tag)
	}
	var group *model.Group
	if req.Group != nil {
		ctx := context.Background()
		group, err = h.Repository.GetGroup(ctx, *req.Group)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	ctx := context.Background()
	request, err := h.Repository.UpdateRequest(ctx, requestID, req.Amount, req.Title, req.Content, tags, group)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	modelcomments, err := h.Repository.GetComments(ctx, requestID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	var comments []*CommentDetail
	for _, modelcomment := range modelcomments {
		comment := &CommentDetail{
			ID:        modelcomment.ID,
			User:      modelcomment.User,
			Comment:   modelcomment.Comment,
			CreatedAt: modelcomment.CreatedAt,
			UpdatedAt: modelcomment.UpdatedAt,
		}
		comments = append(comments, comment)
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
		Comments:  comments,
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) PostComment(c echo.Context) error {
	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if requestID == uuid.Nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var req Comment
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	sess, err := h.SessionStore.Get(c.Request(), h.SessionName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	bodyUser, ok := sess.Values[sessionUserKey].([]byte)
	if !ok {
		c.Logger().Error(errors.New("sessionUser not found"))
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("sessionUser not found"))
	}
	user := new(User)
	err = json.Unmarshal(bodyUser, user)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	ctx := context.Background()
	comment, err := h.Repository.CreateComment(ctx, req.Comment, requestID, user.ID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if requestID == uuid.Nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	commentID, err := uuid.Parse(c.Param("commentID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if commentID == uuid.Nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var req Comment
	if err := c.Bind(&req); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := context.Background()
	comment, err := h.Repository.UpdateComment(ctx, req.Comment, requestID, commentID)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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
