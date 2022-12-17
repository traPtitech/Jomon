package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
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
	Comments  []*CommentDetail  `json:"comments"`
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

func (h *Handlers) GetRequests(c echo.Context) error {
	ctx := c.Request().Context()
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
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		target = &t
	}
	var since *time.Time
	if c.QueryParam("since") != "" {
		s, err := service.StrToDate(c.QueryParam("since"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		since = &s
	}
	var until *time.Time
	if c.QueryParam("until") != "" {
		u, err := service.StrToDate(c.QueryParam("until"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		until = &u
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
	query := model.RequestQuery{
		Sort:   sort,
		Target: target,
		Status: ss,
		Since:  since,
		Until:  until,
		Tag:    tag,
		Group:  group,
	}

	modelrequests, err := h.Repository.GetRequests(ctx, query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tags := []*TagOverview{}
	requests := []*RequestResponse{}
	for _, request := range modelrequests {
		for _, tag := range request.Tags {
			tags = append(tags, &TagOverview{
				ID:        tag.ID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt,
				UpdatedAt: tag.UpdatedAt,
			})
		}

		restargets := []*TargetOverview{}
		for _, target := range request.Targets {
			restargets = append(restargets, &TargetOverview{
				ID:        target.ID,
				Target:    target.Target,
				Amount:    target.Amount,
				PaidAt:    target.PaidAt,
				CreatedAt: target.CreatedAt,
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
			Title:     request.Title,
			Content:   request.Content,
			Targets:   restargets,
			Tags:      tags,
			Group:     resgroup,
			Comments:  []*CommentDetail{},
		}
		requests = append(requests, res)
	}

	return c.JSON(http.StatusOK, requests)
}

func (h *Handlers) PostRequest(c echo.Context) error {
	var req Request
	var err error
	if err = c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	ctx := c.Request().Context()
	var tags []*model.Tag
	for _, tagID := range req.Tags {
		tag, err := h.Repository.GetTag(ctx, *tagID)
		if err != nil {
			if ent.IsNotFound(err) {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		tags = append(tags, tag)
	}
	var targets []*model.RequestTarget
	for _, target := range req.Targets {
		targets = append(targets, &model.RequestTarget{
			Target: target.Target,
			Amount: target.Amount,
		})
	}
	var group *model.Group
	if req.Group != nil {
		group, err = h.Repository.GetGroup(ctx, *req.Group)
		if err != nil {
			if ent.IsNotFound(err) {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	request, err := h.Repository.CreateRequest(ctx, req.Title, req.Content, tags, targets, group, req.CreatedBy)
	if err != nil {
		if ent.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
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
	var reqtargets []*TargetOverview
	for _, target := range request.Targets {
		reqtargets = append(reqtargets, &TargetOverview{
			ID:        target.ID,
			Target:    target.Target,
			Amount:    target.Amount,
			PaidAt:    target.PaidAt,
			CreatedAt: target.CreatedAt,
		})
	}
	var restags []*TagOverview
	for _, tag := range request.Tags {
		restags = append(restags, &TagOverview{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}
	res := &RequestResponse{
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
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) GetRequest(c echo.Context) error {
	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if requestID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}

	ctx := c.Request().Context()
	request, err := h.Repository.GetRequest(ctx, requestID)
	if err != nil {
		if ent.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	modelcomments, err := h.Repository.GetComments(ctx, requestID)
	if err != nil {
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
	var reqtargets []*TargetOverview
	for _, target := range request.Targets {
		reqtargets = append(reqtargets, &TargetOverview{
			ID:        target.ID,
			Target:    target.Target,
			Amount:    target.Amount,
			PaidAt:    target.PaidAt,
			CreatedAt: target.CreatedAt,
		})
	}
	var restags []*TagOverview
	for _, tag := range request.Tags {
		restags = append(restags, &TagOverview{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}
	res := &RequestResponse{
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
		Comments:  comments,
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) PutRequest(c echo.Context) error {
	var req PutRequest
	var err error
	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if requestID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}

	if err = c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	var tags []*model.Tag
	for _, tagID := range req.Tags {
		ctx := c.Request().Context()
		tag, err := h.Repository.GetTag(ctx, *tagID)
		if err != nil {
			if ent.IsNotFound(err) {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		tags = append(tags, tag)
	}
	var targets []*model.RequestTarget
	for _, target := range req.Targets {
		targets = append(targets, &model.RequestTarget{
			Target: target.Target,
			Amount: target.Amount,
		})
	}
	var group *model.Group
	if req.Group != nil {
		ctx := c.Request().Context()
		group, err = h.Repository.GetGroup(ctx, *req.Group)
		if err != nil {
			if ent.IsNotFound(err) {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	ctx := context.Background()
	request, err := h.Repository.UpdateRequest(ctx, requestID, req.Title, req.Content, tags, targets, group)
	if err != nil {
		if ent.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	modelcomments, err := h.Repository.GetComments(ctx, requestID)
	if err != nil {
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
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}
	var restargets []*TargetOverview
	for _, target := range request.Targets {
		restargets = append(restargets, &TargetOverview{
			ID:        target.ID,
			Target:    target.Target,
			Amount:    target.Amount,
			PaidAt:    target.PaidAt,
			CreatedAt: target.CreatedAt,
		})
	}
	res := &RequestResponse{
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
		Comments:  comments,
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) PostComment(c echo.Context) error {
	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if requestID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var req Comment
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	sess, err := session.Get(h.SessionName, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	user, ok := sess.Values[sessionUserKey].(*User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("sessionUser not found"))
	}

	ctx := c.Request().Context()
	comment, err := h.Repository.CreateComment(ctx, req.Comment, requestID, user.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
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

func (h *Handlers) PutStatus(c echo.Context) error {
	var req PutStatus
	var err error
	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if requestID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	sess, err := session.Get(h.SessionName, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	user, ok := sess.Values[sessionUserKey].(*User)
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, errors.New("sessionUser not found"))
	}

	if err = c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	request, err := h.Repository.GetRequest(ctx, requestID)
	if err != nil {
		if ent.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// judging privilege
	if req.Status == request.Status {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid request: same status"))
	}
	if req.Comment == "" {
		if !IsAbleNoCommentChangeStatus(req.Status, request.Status) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("unable to change %v to %v without comment", request.Status.String(), req.Status.String()))
		}
	}

	u, err := h.Repository.GetUserByID(ctx, user.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if u.Admin {
		if !IsAbleAdminChangeState(req.Status, request.Status) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("admin unable to change %v to %v", request.Status.String(), req.Status.String()))
		}
		if req.Status == model.Submitted && request.Status == model.Accepted {
			targets, err := h.Repository.GetRequestTargets(ctx, requestID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
			var paid bool
			for _, target := range targets {
				if target.PaidAt != nil {
					paid = true
					break
				}
			}
			if paid {
				return echo.NewHTTPError(http.StatusBadRequest, errors.New("someone already paid"))
			}
		}
	}

	if !u.Admin && user.ID == request.CreatedBy {
		if !IsAbleCreatorChangeStatus(req.Status, request.Status) {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("creator unable to change %v to %v", request.Status.String(), req.Status.String()))
		}
	}

	if user.ID != request.CreatedBy && !u.Admin {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	// create status and comment: keep the two in this order
	created, err := h.Repository.CreateStatus(ctx, requestID, user.ID, req.Status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	var resComment CommentDetail
	if req.Comment != "" {
		comment, err := h.Repository.CreateComment(ctx, req.Comment, request.ID, user.ID)
		if err != nil {
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
		CreatedBy: user.ID,
		Status:    created.Status,
		Comment:   resComment,
		CreatedAt: created.CreatedAt,
	}

	return c.JSON(http.StatusOK, res)
}

func IsAbleNoCommentChangeStatus(status, latestStatus model.Status) bool {
	if status == model.FixRequired && latestStatus == model.Submitted ||
		status == model.Rejected && latestStatus == model.Submitted ||
		status == model.Submitted && latestStatus == model.Accepted {
		return false
	}
	return true
}

func IsAbleCreatorChangeStatus(status, latestStatus model.Status) bool {
	if status == model.Submitted && latestStatus == model.FixRequired {
		return true
	}
	return false
}

func IsAbleAdminChangeState(status, latestStatus model.Status) bool {
	if status == model.Rejected && latestStatus == model.Submitted ||
		status == model.Submitted && latestStatus == model.FixRequired ||
		status == model.Accepted && latestStatus == model.Submitted ||
		status == model.Submitted && latestStatus == model.Accepted ||
		status == model.FixRequired && latestStatus == model.Submitted {
		return true
	}
	return false
}
