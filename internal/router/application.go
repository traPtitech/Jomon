package router

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/internal/logging"
	"github.com/traPtitech/Jomon/internal/nulltime"
	"github.com/traPtitech/Jomon/internal/service"
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

type Application struct {
	CreatedBy uuid.UUID   `json:"created_by"`
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	Tags      []uuid.UUID `json:"tags"`
	Targets   []*Target   `json:"targets"`
}

type PutApplication struct {
	Title   string      `json:"title"`
	Content string      `json:"content"`
	Tags    []uuid.UUID `json:"tags"`
	Targets []*Target   `json:"targets"`
}

type ApplicationResponse struct {
	ID        uuid.UUID         `json:"id"`
	Status    service.Status    `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	CreatedBy uuid.UUID         `json:"created_by"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Tags      []*TagResponse    `json:"tags"`
	Targets   []*TargetOverview `json:"targets"`
}

type ApplicationDetailResponse struct {
	ApplicationResponse
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
	Status  service.Status `json:"status"`
	Comment string         `json:"comment"`
}

type StatusResponseOverview struct {
	CreatedBy uuid.UUID      `json:"created_by"`
	Status    service.Status `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
}

type StatusResponse struct {
	CreatedBy uuid.UUID      `json:"created_by"`
	Status    service.Status `json:"status"`
	Comment   CommentDetail  `json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
}

type Target struct {
	Target uuid.UUID `json:"target"`
	Amount int       `json:"amount"`
}

type TargetOverview struct {
	ID        uuid.UUID         `json:"id"`
	Target    uuid.UUID         `json:"target"`
	Amount    int               `json:"amount"`
	PaidAt    nulltime.NullTime `json:"paid_at"`
	CreatedAt time.Time         `json:"created_at"`
}

func (h Handlers) GetApplications(c *echo.Context) error {
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
			return service.NewBadInputError("invalid status")
		}
	}
	if s := status.String(); s != "" {
		ss = &s
	}
	var target uuid.UUID
	if c.QueryParam("target") != "" {
		t, err := uuid.Parse(c.QueryParam("target"))
		if err != nil {
			logger.Info("could not parse query parameter `target` as UUID", zap.Error(err))
			return service.NewBadInputError("received parameter `target` is not a valid UUID").
				WithInternal(err)
		}
		target = t
	}
	var since nulltime.NullTime
	if c.QueryParam("since") != "" {
		s, err := nulltime.ParseDate(c.QueryParam("since"))
		if err != nil {
			logger.Info("could not parse query parameter `since` as time.Time", zap.Error(err))
			return service.NewBadInputError("received parameter `since` is not a valid date").
				WithInternal(err)
		}
		since = s
	}
	var until nulltime.NullTime
	if c.QueryParam("until") != "" {
		u, err := nulltime.ParseDate(c.QueryParam("until"))
		if err != nil {
			logger.Info("could not parse query parameter `until` as time.Time", zap.Error(err))
			return service.NewBadInputError("received parameter `until` is not a valid date").
				WithInternal(err)
		}
		until = u
	}
	limit := 100
	if limitQuery := c.QueryParam("limit"); limitQuery != "" {
		limitI, err := strconv.Atoi(limitQuery)
		if err != nil {
			logger.Info("could not parse limit as integer", zap.Error(err))
			return service.NewBadInputError("received parameter `limit` is not a valid integer").
				WithInternal(err)
		}
		if limitI < 0 {
			logger.Info("received negative limit", zap.Int("limit", limitI))
			return service.NewBadInputError(fmt.Sprintf("negative limit(=%d) is invalid", limitI))
		}
		limit = limitI
	}
	offset := 0
	if offsetQuery := c.QueryParam("offset"); offsetQuery != "" {
		offsetI, err := strconv.Atoi(offsetQuery)
		if err != nil {
			logger.Info("could not parse offset as integer", zap.Error(err))
			return service.NewBadInputError("received parameter `offset` is not a valid integer").
				WithInternal(err)
		}
		if offsetI < 0 {
			logger.Info("received negative offset", zap.Int("offset", offsetI))
			return service.NewBadInputError(fmt.Sprintf("negative offset(=%d) is invalid", offsetI))
		}
		offset = offsetI
	}
	var tag *string
	if c.QueryParam("tag") != "" {
		t := c.QueryParam("tag")
		tag = &t
	}
	var createdBy uuid.UUID
	if c.QueryParam("created_by") != "" {
		u, err := uuid.Parse(c.QueryParam("created_by"))
		if err != nil {
			logger.Info("could not parse query parameter `created_by` as UUID", zap.Error(err))
			return service.NewBadInputError("received parameter `created_by` is not a valid UUID").
				WithInternal(err)
		}
		createdBy = u
	}
	query := service.ApplicationQuery{
		Sort:      sort,
		Target:    target,
		Status:    ss,
		Since:     since,
		Until:     until,
		Limit:     limit,
		Offset:    offset,
		Tag:       tag,
		CreatedBy: createdBy,
	}

	modelapplications, err := h.Service.GetApplications(ctx, query)
	if err != nil {
		logger.Error("failed to get applications from service", zap.Error(err))
		return err
	}

	applications := lo.Map(
		modelapplications,
		func(application *service.ApplicationResponse, _ int) *ApplicationResponse {
			restags := lo.Map(application.Tags, func(tag *service.Tag, _ int) *TagResponse {
				return &TagResponse{
					ID:        tag.ID,
					Name:      tag.Name,
					CreatedAt: tag.CreatedAt,
					UpdatedAt: tag.UpdatedAt,
				}
			})

			restargets := lo.Map(
				application.Targets,
				func(target *service.ApplicationTargetDetail, _ int) *TargetOverview {
					return &TargetOverview{
						ID:        target.ID,
						Target:    target.Target,
						Amount:    target.Amount,
						PaidAt:    target.PaidAt,
						CreatedAt: target.CreatedAt,
					}
				},
			)

			return &ApplicationResponse{
				ID:        application.ID,
				Status:    application.Status,
				CreatedAt: application.CreatedAt,
				UpdatedAt: application.UpdatedAt,
				CreatedBy: application.CreatedBy,
				Title:     application.Title,
				Content:   application.Content,
				Targets:   restargets,
				Tags:      restags,
			}
		},
	)

	return c.JSON(http.StatusOK, applications)
}

func (h Handlers) PostApplication(c *echo.Context) error {
	var req Application
	var err error
	if err = c.Bind(&req); err != nil {
		return service.NewBadInputError("failed to get application from request").
			WithInternal(err)
	}
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	targets := lo.Map(req.Targets, func(target *Target, _ int) *service.ApplicationTarget {
		return &service.ApplicationTarget{
			Target: target.Target,
			Amount: target.Amount,
		}
	})
	tags := make([]*service.Tag, 0, len(req.Tags))
	for _, tagID := range req.Tags {
		tag, err := h.Service.GetTag(ctx, tagID)
		if err != nil {
			return err
		}
		tags = append(tags, tag)
	}
	inputs := service.CreateApplicationInputs{
		Title:   req.Title,
		Content: req.Content,
		Tags:    tags,
		Targets: targets,
	}
	application, err := h.Service.CreateApplication(ctx, req.CreatedBy, inputs)
	if err != nil {
		logger.Error("failed to create application in service", zap.Error(err))
		return err
	}
	restargets := lo.Map(
		application.Targets,
		func(target *service.ApplicationTargetDetail, _ int) *TargetOverview {
			return &TargetOverview{
				ID:        target.ID,
				Target:    target.Target,
				Amount:    target.Amount,
				PaidAt:    target.PaidAt,
				CreatedAt: target.CreatedAt,
			}
		},
	)
	restags := lo.Map(application.Tags, func(tag *service.Tag, _ int) *TagResponse {
		return &TagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}
	})
	comments := lo.Map(
		application.Comments,
		func(comment *service.Comment, _ int) *CommentDetail {
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
		application.Statuses,
		func(status *service.ApplicationStatus, _ int) *StatusResponseOverview {
			return &StatusResponseOverview{
				Status:    status.Status,
				CreatedAt: status.CreatedAt,
				CreatedBy: status.CreatedBy,
			}
		},
	)
	files := application.Files

	res := &ApplicationDetailResponse{
		ApplicationResponse: ApplicationResponse{
			ID:        application.ID,
			Status:    application.Status,
			CreatedAt: application.CreatedAt,
			UpdatedAt: application.UpdatedAt,
			CreatedBy: application.CreatedBy,
			Title:     application.Title,
			Content:   application.Content,
			Tags:      restags,
			Targets:   restargets,
		},
		Comments: comments,
		Statuses: statuses,
		Files:    files,
	}
	return c.JSON(http.StatusOK, res)
}

func (h Handlers) GetApplication(c *echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	applicationID, err := uuid.Parse(c.Param("applicationID"))
	if err != nil {
		logger.Info("could not parse query parameter `applicationID` as UUID", zap.Error(err))
		return service.NewBadInputError("received parameter `applicationID` is not a valid UUID").
			WithInternal(err)
	}
	if applicationID == uuid.Nil {
		logger.Info("invalid UUID")
		return service.NewBadInputError("received parameter `applicationID` is not a valid UUID")
	}

	application, err := h.Service.GetApplication(ctx, applicationID)
	if err != nil {
		logger.Error("failed to get application from service", zap.Error(err))
		return err
	}
	restargets := lo.Map(
		application.Targets,
		func(target *service.ApplicationTargetDetail, _ int) *TargetOverview {
			return &TargetOverview{
				ID:        target.ID,
				Target:    target.Target,
				Amount:    target.Amount,
				PaidAt:    target.PaidAt,
				CreatedAt: target.CreatedAt,
			}
		},
	)
	restags := lo.Map(application.Tags, func(tag *service.Tag, _ int) *TagResponse {
		return &TagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}
	})

	comments := lo.Map(application.Comments, func(comment *service.Comment, _ int) *CommentDetail {
		return &CommentDetail{
			ID:        comment.ID,
			User:      comment.User,
			Comment:   comment.Comment,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}
	})
	statuses := lo.Map(
		application.Statuses,
		func(status *service.ApplicationStatus, _ int) *StatusResponseOverview {
			return &StatusResponseOverview{
				CreatedBy: status.CreatedBy,
				Status:    status.Status,
				CreatedAt: status.CreatedAt,
			}
		},
	)
	files := application.Files

	res := &ApplicationDetailResponse{
		ApplicationResponse: ApplicationResponse{
			ID:        application.ID,
			Status:    application.Status,
			CreatedAt: application.CreatedAt,
			UpdatedAt: application.UpdatedAt,
			CreatedBy: application.CreatedBy,
			Title:     application.Title,
			Content:   application.Content,
			Tags:      restags,
			Targets:   restargets,
		},
		Statuses: statuses,
		Comments: comments,
		Files:    files,
	}
	return c.JSON(http.StatusOK, res)
}

func (h Handlers) PutApplication(c *echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	loginUser, _ := c.Get(loginUserKey).(User)
	applicationID, err := uuid.Parse(c.Param("applicationID"))
	if err != nil {
		logger.Info("could not parse query parameter `applicationID` as UUID", zap.Error(err))
		return service.NewBadInputError("received parameter `applicationID` is not a valid UUID").
			WithInternal(err)
	}
	if applicationID == uuid.Nil {
		logger.Info("invalid UUID")
		return service.NewBadInputError("received parameter `applicationID` is not a valid UUID")
	}
	var req PutApplication
	if err = c.Bind(&req); err != nil {
		logger.Info("failed to get application from request", zap.Error(err))
		return service.NewBadInputError("failed to get application from request").
			WithInternal(err)
	}
	targets := lo.Map(req.Targets, func(target *Target, _ int) *service.ApplicationTarget {
		return &service.ApplicationTarget{
			Target: target.Target,
			Amount: target.Amount,
		}
	})
	inputs := service.UpdateApplicationInputs{
		Title:   req.Title,
		Content: req.Content,
		TagIDs:  req.Tags,
		Targets: targets,
	}
	application, err := h.Service.UpdateApplication(ctx, applicationID, loginUser.ID, inputs)
	if err != nil {
		logger.Error("failed to update application in service", zap.Error(err))
		return err
	}
	restags := lo.Map(application.Tags, func(tag *service.Tag, _ int) *TagResponse {
		return &TagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}
	})

	restargets := lo.Map(
		application.Targets,
		func(target *service.ApplicationTargetDetail, _ int) *TargetOverview {
			return &TargetOverview{
				ID:        target.ID,
				Target:    target.Target,
				Amount:    target.Amount,
				PaidAt:    target.PaidAt,
				CreatedAt: target.CreatedAt,
			}
		},
	)

	comments := lo.Map(application.Comments, func(c *service.Comment, _ int) *CommentDetail {
		return &CommentDetail{
			ID:        c.ID,
			User:      c.User,
			Comment:   c.Comment,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	})
	statuses := lo.Map(
		application.Statuses,
		func(status *service.ApplicationStatus, _ int) *StatusResponseOverview {
			return &StatusResponseOverview{
				CreatedBy: status.CreatedBy,
				Status:    status.Status,
				CreatedAt: status.CreatedAt,
			}
		},
	)
	files := application.Files

	res := &ApplicationDetailResponse{
		ApplicationResponse: ApplicationResponse{
			ID:        application.ID,
			Status:    application.Status,
			CreatedAt: application.CreatedAt,
			UpdatedAt: application.UpdatedAt,
			CreatedBy: application.CreatedBy,
			Title:     application.Title,
			Content:   application.Content,
			Tags:      restags,
			Targets:   restargets,
		},
		Comments: comments,
		Statuses: statuses,
		Files:    files,
	}
	return c.JSON(http.StatusOK, res)
}

func (h Handlers) PostComment(c *echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	loginUser, _ := c.Get(loginUserKey).(User)
	applicationID, err := uuid.Parse(c.Param("applicationID"))
	if err != nil {
		logger.Info("could not parse query parameter `applicationID` as UUID", zap.Error(err))
		return service.NewBadInputError("received parameter `applicationID` is not a valid UUID").
			WithInternal(err)
	}
	if applicationID == uuid.Nil {
		logger.Info("invalid UUID")
		return service.NewBadInputError("received parameter `applicationID` is not a valid UUID")
	}

	var req Comment
	if err := c.Bind(&req); err != nil {
		logger.Info("failed to get comment from request", zap.Error(err))
		return service.NewBadInputError("failed to get comment from request").
			WithInternal(err)
	}

	comment, err := h.Service.CreateCommentToApplication(
		ctx, applicationID, loginUser.ID, req.Comment)
	if err != nil {
		logger.Error("failed to create comment in service", zap.Error(err))
		return err
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

func (h Handlers) PutStatus(c *echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	loginUser, _ := c.Get(loginUserKey).(User)
	applicationID, err := uuid.Parse(c.Param("applicationID"))
	if err != nil {
		logger.Info("could not parse query parameter `applicationID` as UUID", zap.Error(err))
		return service.NewBadInputError("received parameter `applicationID` is not a valid UUID").
			WithInternal(err)
	}
	if applicationID == uuid.Nil {
		logger.Info("invalid UUID")
		return service.NewBadInputError("received parameter `applicationID` is not a valid UUID")
	}

	var req PutStatus
	if err = c.Bind(&req); err != nil {
		logger.Info("could not get status from request", zap.Error(err))
		return service.NewBadInputError("could not get status from request").
			WithInternal(err)
	}

	inputs := service.UpdateStatusInputs{
		Status:  req.Status,
		Comment: req.Comment,
	}
	created, comment, err := h.Service.UpdateApplicationStatus(
		ctx, applicationID, loginUser.ID, inputs)
	if err != nil {
		logger.Error("failed to update application status in service", zap.Error(err))
		return err
	}
	resComment := CommentDetail{}
	if comment != nil {
		resComment = CommentDetail{
			ID:        comment.ID,
			User:      comment.User,
			Comment:   comment.Comment,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}
	}

	res := &StatusResponse{
		CreatedBy: created.CreatedBy,
		Status:    created.Status,
		Comment:   resComment,
		CreatedAt: created.CreatedAt,
	}

	return c.JSON(http.StatusOK, res)
}
