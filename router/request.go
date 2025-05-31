package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/logging"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
	"go.uber.org/zap"
)

type Status string

const (
	Submitted   Status = "submitted"
	FixRequired Status = "fix_required"
	Accepted    Status = "accepted"
	Completed   Status = "completed"
	Rejected    Status = "rejected"
)

func (s Status) Valid() bool {
	switch s {
	case Submitted, FixRequired, Accepted, Completed, Rejected:
		return true
	default:
		return false
	}
}

func (s Status) String() string {
	return string(s)
}

type Request struct {
	CreatedBy uuid.UUID    `json:"created_by"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Tags      []*uuid.UUID `json:"tags"`
	Targets   []*Target    `json:"targets"`
	Group     *uuid.UUID   `json:"group"`
}

type PutRequest struct {
	Title   string       `json:"title"`
	Content string       `json:"content"`
	Tags    []*uuid.UUID `json:"tags"`
	Targets []*Target    `json:"targets"`
	Group   *uuid.UUID   `json:"group"`
}

type RequestResponse struct {
	ID        uuid.UUID         `json:"id"`
	Status    model.Status      `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	CreatedBy uuid.UUID         `json:"created_by"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Tags      []*TagOverview    `json:"tags"`
	Targets   []*TargetOverview `json:"targets"`
	Group     *GroupOverview    `json:"group"`
}

type RequestDetailResponse struct {
	RequestResponse
	Comments []*CommentDetail          `json:"comments"`
	Statuses []*StatusResponseOverview `json:"statuses"`
	Files    []uuid.UUID               `json:"files"`
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

type PutStatus struct {
	Status  model.Status `json:"status"`
	Comment string       `json:"comment"`
}

type StatusResponseOverview struct {
	CreatedBy uuid.UUID    `json:"created_by"`
	Status    model.Status `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
}

type StatusResponse struct {
	CreatedBy uuid.UUID     `json:"created_by"`
	Status    model.Status  `json:"status"`
	Comment   CommentDetail `json:"comment"`
	CreatedAt time.Time     `json:"created_at"`
}

type Target struct {
	Target uuid.UUID `json:"target"`
	Amount int       `json:"amount"`
}

type TargetOverview struct {
	ID        uuid.UUID  `json:"id"`
	Target    uuid.UUID  `json:"target"`
	Amount    int        `json:"amount"`
	PaidAt    *time.Time `json:"paid_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (h Handlers) GetRequests(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	var sort *string
	if s := c.QueryParam("sort"); s != "" {
		sort = &s
	}
	var status Status
	var ss *string
	if s := c.QueryParam("status"); s != "" {
		status = Status(s)
		if !status.Valid() {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid status")
		}
	}
	if s := status.String(); s != "" {
		ss = &s
	}
	var target *uuid.UUID
	if c.QueryParam("target") != "" {
		t, err := uuid.Parse(c.QueryParam("target"))
		if err != nil {
			logger.Info("could not parse query parameter `target` as UUID", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		target = &t
	}
	var since *time.Time
	if c.QueryParam("since") != "" {
		s, err := service.StrToDate(c.QueryParam("since"))
		if err != nil {
			logger.Info("could not parse query parameter `since` as time.Time", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		since = &s
	}
	var until *time.Time
	if c.QueryParam("until") != "" {
		u, err := service.StrToDate(c.QueryParam("until"))
		if err != nil {
			logger.Info("could not parse query parameter `until` as time.Time", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		until = &u
	}
	limit := 100
	if limitQuery := c.QueryParam("limit"); limitQuery != "" {
		limitI, err := strconv.Atoi(limitQuery)
		if err != nil {
			logger.Info("could not parse limit as integer", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		if limitI < 0 {
			logger.Info("received negative limit", zap.Int("limit", limitI))
			return echo.NewHTTPError(
				http.StatusBadRequest,
				fmt.Errorf("negative limit(=%d) is invalid", limitI),
			)
		}
		limit = limitI
	}
	offset := 0
	if offsetQuery := c.QueryParam("offset"); offsetQuery != "" {
		offsetI, err := strconv.Atoi(offsetQuery)
		if err != nil {
			logger.Info("could not parse offset as integer", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		if offsetI < 0 {
			logger.Info("received negative offset", zap.Int("offset", offsetI))
			return echo.NewHTTPError(
				http.StatusBadRequest,
				fmt.Errorf("negative offset(=%d) is invalid", offsetI),
			)
		}
		offset = offsetI
	}
	var tag *string
	if c.QueryParam("tag") != "" {
		t := c.QueryParam("tag")
		tag = &t
	}
	var group *string
	if c.QueryParam("group") != "" {
		g := c.QueryParam("group")
		group = &g
	}
	var cratedBy *uuid.UUID
	if c.QueryParam("created_by") != "" {
		u, err := uuid.Parse(c.QueryParam("created_by"))
		if err != nil {
			logger.Info("could not parse query parameter `created_by` as UUID", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		cratedBy = &u
	}
	query := model.RequestQuery{
		Sort:      sort,
		Target:    target,
		Status:    ss,
		Since:     since,
		Until:     until,
		Limit:     limit,
		Offset:    offset,
		Tag:       tag,
		Group:     group,
		CreatedBy: cratedBy,
	}

	modelrequests, err := h.Repository.GetRequests(ctx, query)
	if err != nil {
		logger.Error("failed to get requests from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	requests := lo.Map(
		modelrequests,
		func(request *model.RequestResponse, _ int) *RequestResponse {
			tags := lo.Map(request.Tags, func(tag *model.Tag, _ int) *TagOverview {
				return &TagOverview{
					ID:        tag.ID,
					Name:      tag.Name,
					CreatedAt: tag.CreatedAt,
					UpdatedAt: tag.UpdatedAt,
				}
			})

			restargets := lo.Map(
				request.Targets,
				func(target *model.RequestTargetDetail, _ int) *TargetOverview {
					return &TargetOverview{
						ID:        target.ID,
						Target:    target.Target,
						Amount:    target.Amount,
						PaidAt:    target.PaidAt,
						CreatedAt: target.CreatedAt,
					}
				},
			)

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

			return &RequestResponse{
				ID:        request.ID,
				Status:    request.Status,
				CreatedAt: request.CreatedAt,
				UpdatedAt: request.UpdatedAt,
				CreatedBy: request.CreatedBy,
				Title:     request.Title,
				Content:   request.Content,
				Targets:   restargets,
				Tags:      tags,
				Group:     resgroup,
			}
		},
	)

	return c.JSON(http.StatusOK, requests)
}

func (h Handlers) PostRequest(c echo.Context) error {
	var req Request
	var err error
	if err = c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	tags := []*model.Tag{}
	for _, tagID := range req.Tags {
		tag, err := h.Repository.GetTag(ctx, *tagID)
		if err != nil {
			if ent.IsNotFound(err) {
				logger.Info("could not find tag in repository", zap.String("ID", tagID.String()))
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			logger.Error("failed to get tag from repository", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		tags = append(tags, tag)
	}
	targets := lo.Map(req.Targets, func(target *Target, _ int) *model.RequestTarget {
		return &model.RequestTarget{
			Target: target.Target,
			Amount: target.Amount,
		}
	})
	var group *model.Group
	if req.Group != nil {
		group, err = h.Repository.GetGroup(ctx, *req.Group)
		if err != nil {
			if ent.IsNotFound(err) {
				logger.Info(
					"could not find group in repository",
					zap.String("ID", req.Group.String()))
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			logger.Error("failed to get group from repository", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	request, err := h.Repository.CreateRequest(
		ctx,
		req.Title, req.Content, tags, targets, group, req.CreatedBy)
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Info(
				"could not find request in repository",
				zap.String("ID", req.CreatedBy.String()))
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		logger.Error("failed to create request in repository", zap.Error(err))
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
	reqtargets := lo.Map(
		request.Targets,
		func(target *model.RequestTargetDetail, _ int) *TargetOverview {
			return &TargetOverview{
				ID:        target.ID,
				Target:    target.Target,
				Amount:    target.Amount,
				PaidAt:    target.PaidAt,
				CreatedAt: target.CreatedAt,
			}
		},
	)
	restags := lo.Map(request.Tags, func(tag *model.Tag, _ int) *TagOverview {
		return &TagOverview{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}
	})
	comments := lo.Map(
		request.Comments,
		func(comment *model.Comment, _ int) *CommentDetail {
			return &CommentDetail{
				ID:        comment.ID,
				User:      comment.User,
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			}
		},
	)
	statuses := lo.Map(
		request.Statuses,
		func(status *model.RequestStatus, _ int) *StatusResponseOverview {
			return &StatusResponseOverview{
				Status:    status.Status,
				CreatedAt: status.CreatedAt,
				CreatedBy: status.CreatedBy,
			}
		},
	)
	files := lo.Map(
		request.Files,
		func(file *uuid.UUID, _ int) uuid.UUID {
			return *file
		},
	)

	res := &RequestDetailResponse{
		RequestResponse: RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
			Tags:      restags,
			Targets:   reqtargets,
			Group:     resgroup,
		},
		Comments: comments,
		Statuses: statuses,
		Files:    files,
	}
	return c.JSON(http.StatusOK, res)
}

func (h Handlers) GetRequest(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		logger.Info("could not parse query parameter `requestID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if requestID == uuid.Nil {
		logger.Info("invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}

	request, err := h.Repository.GetRequest(ctx, requestID)
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Info(
				"could not find request in repository",
				zap.String("ID", requestID.String()),
				zap.Error(err))
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		logger.Error("failed to get request from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	modelcomments, err := h.Repository.GetComments(ctx, requestID)
	if err != nil {
		logger.Error("failed to get comments from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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
	reqtargets := lo.Map(
		request.Targets,
		func(target *model.RequestTargetDetail, _ int) *TargetOverview {
			return &TargetOverview{
				ID:        target.ID,
				Target:    target.Target,
				Amount:    target.Amount,
				PaidAt:    target.PaidAt,
				CreatedAt: target.CreatedAt,
			}
		},
	)
	restags := lo.Map(request.Tags, func(tag *model.Tag, _ int) *TagOverview {
		return &TagOverview{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}
	})

	comments := lo.Map(modelcomments, func(modelcomment *model.Comment, _ int) *CommentDetail {
		return &CommentDetail{
			ID:        modelcomment.ID,
			User:      modelcomment.User,
			Comment:   modelcomment.Comment,
			CreatedAt: modelcomment.CreatedAt,
			UpdatedAt: modelcomment.UpdatedAt,
		}
	})
	statuses := lo.Map(
		request.Statuses,
		func(status *model.RequestStatus, _ int) *StatusResponseOverview {
			return &StatusResponseOverview{
				CreatedBy: status.CreatedBy,
				Status:    status.Status,
				CreatedAt: status.CreatedAt,
			}
		},
	)
	files := lo.Map(request.Files, func(file *uuid.UUID, _ int) uuid.UUID {
		return *file
	})

	res := &RequestDetailResponse{
		RequestResponse: RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
			Tags:      restags,
			Targets:   reqtargets,
			Group:     resgroup,
		},
		Statuses: statuses,
		Comments: comments,
		Files:    files,
	}
	return c.JSON(http.StatusOK, res)
}

func (h Handlers) PutRequest(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	loginUser, _ := c.Get(loginUserKey).(User)
	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		logger.Info("could not parse query parameter `requestID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if requestID == uuid.Nil {
		logger.Info("invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	isRequestCreator, err := h.isRequestCreator(ctx, loginUser.ID, requestID)
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Info(
				"could not find request in repository",
				zap.String("ID", requestID.String()),
				zap.Error(err))
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		logger.Error("failed to check request creator", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if !isRequestCreator {
		logger.Info("user is not request creator", zap.String("ID", loginUser.ID.String()))
		return echo.NewHTTPError(http.StatusForbidden, "you are not request creator")
	}

	var req PutRequest
	if err = c.Bind(&req); err != nil {
		logger.Info("failed to get request from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	tags := []*model.Tag{}
	for _, tagID := range req.Tags {
		tag, err := h.Repository.GetTag(ctx, *tagID)
		if err != nil {
			if ent.IsNotFound(err) {
				logger.Info("could not find tag in repository", zap.String("ID", tagID.String()))
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			logger.Error("failed to get tag from repository", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		tags = append(tags, tag)
	}
	targets := lo.Map(req.Targets, func(target *Target, _ int) *model.RequestTarget {
		return &model.RequestTarget{
			Target: target.Target,
			Amount: target.Amount,
		}
	})
	var group *model.Group
	if req.Group != nil {
		group, err = h.Repository.GetGroup(ctx, *req.Group)
		if err != nil {
			if ent.IsNotFound(err) {
				logger.Info(
					"could not find group in repository",
					zap.String("ID", req.Group.String()))
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			logger.Error("failed to get group from repository", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	request, err := h.Repository.UpdateRequest(
		ctx,
		requestID, req.Title, req.Content, tags, targets, group)
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Info(
				"could not find request in repository",
				zap.String("ID", requestID.String()))
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		logger.Error("failed to update request in repository", zap.Error(err))
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
	restags := lo.Map(request.Tags, func(tag *model.Tag, _ int) *TagOverview {
		return &TagOverview{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}
	})

	restargets := lo.Map(
		request.Targets,
		func(target *model.RequestTargetDetail, _ int) *TargetOverview {
			return &TargetOverview{
				ID:        target.ID,
				Target:    target.Target,
				Amount:    target.Amount,
				PaidAt:    target.PaidAt,
				CreatedAt: target.CreatedAt,
			}
		},
	)

	comments := lo.Map(request.Comments, func(c *model.Comment, _ int) *CommentDetail {
		return &CommentDetail{
			ID:        c.ID,
			User:      c.User,
			Comment:   c.Comment,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	})
	statuses := lo.Map(
		request.Statuses,
		func(status *model.RequestStatus, _ int) *StatusResponseOverview {
			return &StatusResponseOverview{
				CreatedBy: status.CreatedBy,
				Status:    status.Status,
				CreatedAt: status.CreatedAt,
			}
		},
	)
	files := lo.Map(request.Files, func(file *uuid.UUID, _ int) uuid.UUID {
		return *file
	})

	res := &RequestDetailResponse{
		RequestResponse: RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
			Tags:      restags,
			Targets:   restargets,
			Group:     resgroup,
		},
		Comments: comments,
		Statuses: statuses,
		Files:    files,
	}
	return c.JSON(http.StatusOK, res)
}

func (h Handlers) PostComment(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		logger.Info("could not parse query parameter `requestID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if requestID == uuid.Nil {
		logger.Info("invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var req Comment
	if err := c.Bind(&req); err != nil {
		logger.Info("failed to get comment from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	sess, err := session.Get(h.SessionName, c)
	if err != nil {
		logger.Error("failed to get session", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	user, ok := sess.Values[sessionUserKey].(User)
	if !ok {
		logger.Info("could not find use in session")
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("sessionUser not found"))
	}

	comment, err := h.Repository.CreateComment(ctx, req.Comment, requestID, user.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Info(
				"could not find request in repository",
				zap.String("ID", requestID.String()))
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		logger.Error("failed to create comment in repository", zap.Error(err))
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

func (h Handlers) PutStatus(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	loginUser, _ := c.Get(loginUserKey).(User)
	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		logger.Info("could not parse query parameter `requestID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if requestID == uuid.Nil {
		logger.Info("invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := h.filterAdminOrRequestCreator(ctx, loginUser, requestID); err != nil {
		return err
	}

	var req PutStatus
	if err = c.Bind(&req); err != nil {
		logger.Info("could not get status from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	request, err := h.Repository.GetRequest(ctx, requestID)
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Info(
				"could not find request in repository",
				zap.String("ID", requestID.String()))
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		logger.Error("failed to get request from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// judging privilege
	if req.Status == request.Status {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid request: same status"))
	}
	if req.Comment == "" {
		if !IsAbleNoCommentChangeStatus(req.Status, request.Status) {
			err := fmt.Errorf(
				"unable to change %v to %v without comment",
				request.Status.String(), req.Status.String())
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	u, err := h.Repository.GetUserByID(ctx, loginUser.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Info("could not find user in repository", zap.String("ID", loginUser.ID.String()))
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		logger.Error("failed to get user from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if u.Admin {
		if !IsAbleAdminChangeState(req.Status, request.Status) {
			logger.Info("admin unable to change status")
			err := fmt.Errorf(
				"admin unable to change %v to %v",
				request.Status.String(), req.Status.String())
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		if req.Status == model.Submitted && request.Status == model.Accepted {
			targets, err := h.Repository.GetRequestTargets(ctx, requestID)
			if err != nil {
				logger.Error("failed to get request targets from repository", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			paid := lo.Reduce(targets, func(p bool, target *model.RequestTargetDetail, _ int) bool {
				return p || target.PaidAt != nil
			}, false)
			if paid {
				logger.Info("someone already paid")
				return echo.NewHTTPError(http.StatusBadRequest, errors.New("someone already paid"))
			}
		}
	}

	if !u.Admin && loginUser.ID == request.CreatedBy {
		if !IsAbleCreatorChangeStatus(req.Status, request.Status) {
			logger.Info("creator unable to change status")
			err := fmt.Errorf(
				"creator unable to change %v to %v",
				request.Status.String(), req.Status.String())
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	if loginUser.ID != request.CreatedBy && !u.Admin {
		logger.Info("use is not creator or admin")
		return echo.NewHTTPError(http.StatusForbidden)
	}

	// create status and comment: keep the two in this order
	created, err := h.Repository.CreateStatus(ctx, requestID, loginUser.ID, req.Status)
	if err != nil {
		logger.Error("failed to create status in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	var resComment CommentDetail
	if req.Comment != "" {
		comment, err := h.Repository.CreateComment(ctx, req.Comment, request.ID, loginUser.ID)
		if err != nil {
			logger.Error("failed to create comment in repository", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		resComment = CommentDetail{
			ID:        comment.ID,
			User:      comment.User,
			Comment:   comment.Comment,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}
	}

	res := &StatusResponse{
		CreatedBy: loginUser.ID,
		Status:    created.Status,
		Comment:   resComment,
		CreatedAt: created.CreatedAt,
	}

	return c.JSON(http.StatusOK, res)
}

func IsAbleNoCommentChangeStatus(status, latestStatus model.Status) bool {
	switch latestStatus {
	case model.Submitted:
		return status != model.FixRequired && status != model.Rejected
	case model.Accepted:
		return status != model.Submitted
	case model.FixRequired, model.Completed, model.Rejected:
		return true
	}
	// the switch above performs exhaustive check
	panic("unreachable")
}

func IsAbleCreatorChangeStatus(status, latestStatus model.Status) bool {
	return status == model.Submitted && latestStatus == model.FixRequired
}

func IsAbleAdminChangeState(status, latestStatus model.Status) bool {
	return status == model.Rejected && latestStatus == model.Submitted ||
		status == model.Submitted && latestStatus == model.FixRequired ||
		status == model.Accepted && latestStatus == model.Submitted ||
		status == model.Submitted && latestStatus == model.Accepted ||
		status == model.FixRequired && latestStatus == model.Submitted
}

func (h Handlers) isRequestCreator(
	ctx context.Context, userID, requestID uuid.UUID,
) (bool, error) {
	request, err := h.Repository.GetRequest(ctx, requestID)
	if err != nil {
		return false, err
	}
	return request.CreatedBy == userID, nil
}

func (h Handlers) filterAdminOrRequestCreator(
	ctx context.Context, user User, requestID uuid.UUID,
) *echo.HTTPError {
	logger := logging.GetLogger(ctx)
	if user.Admin {
		return nil
	}
	request, err := h.Repository.GetRequest(ctx, requestID)
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Info(
				"could not find request in repository",
				zap.String("ID", requestID.String()))
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		logger.Error("failed to get request from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if request.CreatedBy == user.ID {
		return nil
	}
	logger.Info("user is not admin or request creator", zap.String("ID", user.ID.String()))
	return echo.NewHTTPError(http.StatusForbidden, "you are not request creator")
}
