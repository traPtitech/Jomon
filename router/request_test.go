package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
	"github.com/traPtitech/Jomon/testutil"
	"github.com/traPtitech/Jomon/testutil/random"
	"go.uber.org/mock/gomock"
)

func modelTagToTagOverview(t *model.Tag) *TagResponse {
	return &TagResponse{
		ID:        t.ID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func modelRequestTargetDetailToTargetOverview(t *model.RequestTargetDetail) *TargetOverview {
	return &TargetOverview{
		ID:        t.ID,
		Target:    t.Target,
		Amount:    t.Amount,
		CreatedAt: t.CreatedAt,
	}
}

func modelRequestStatusToStatusResponseOverview(s *model.RequestStatus) *StatusResponseOverview {
	return &StatusResponseOverview{
		CreatedBy: s.CreatedBy,
		Status:    s.Status,
		CreatedAt: s.CreatedAt,
	}
}

func modelCommentToCommentDetail(c *model.Comment) *CommentDetail {
	return &CommentDetail{
		ID:        c.ID,
		User:      c.User,
		Comment:   c.Comment,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// FIXME: この処理はrequest.goにも書かれてある
func modelRequestDetailToRequestResponse(r *model.RequestDetail) *RequestDetailResponse {
	var group *GroupOverview
	if r.Group != nil {
		group = &GroupOverview{
			ID:          r.Group.ID,
			Name:        r.Group.Name,
			Description: r.Group.Description,
			Budget:      r.Group.Budget,
			CreatedAt:   r.Group.CreatedAt,
			UpdatedAt:   r.Group.UpdatedAt,
		}
	}
	return &RequestDetailResponse{
		RequestResponse: RequestResponse{
			ID:        r.ID,
			Status:    r.Status,
			CreatedBy: r.CreatedBy,
			Title:     r.Title,
			Content:   r.Content,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			Targets: lo.Map(r.Targets, func(m *model.RequestTargetDetail, _ int) *TargetOverview {
				return modelRequestTargetDetailToTargetOverview(m)
			}),
			Tags: lo.Map(r.Tags, func(m *model.Tag, _ int) *TagResponse {
				return modelTagToTagOverview(m)
			}),
			Group: group,
		},
		Statuses: lo.Map(r.Statuses, func(m *model.RequestStatus, _ int) *StatusResponseOverview {
			return modelRequestStatusToStatusResponseOverview(m)
		}),
		Comments: lo.Map(r.Comments, func(m *model.Comment, _ int) *CommentDetail {
			return modelCommentToCommentDetail(m)
		}),
		Files: lo.Map(r.Files, func(f *uuid.UUID, _ int) uuid.UUID {
			return *f
		}),
	}
}

func TestHandlers_GetRequests(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date1 := time.Now()
		date2 := date1.Add(time.Hour)
		request1 := &model.RequestResponse{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date1,
			UpdatedAt: date1,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses:  []*model.RequestStatus{},
			Group:     nil,
		}
		request2 := &model.RequestResponse{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date2,
			UpdatedAt: date2,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses:  []*model.RequestStatus{},
			Group:     nil,
		}
		requests := []*model.RequestResponse{request2, request1}

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/requests", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), model.RequestQuery{
				Limit:  100,
				Offset: 0,
			}).
			Return(requests, nil)

		require.NoError(t, h.Handlers.GetRequests(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*RequestResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := []*RequestResponse{
			{
				ID:        request2.ID,
				Status:    request2.Status,
				CreatedAt: request2.CreatedAt,
				UpdatedAt: request2.UpdatedAt,
				CreatedBy: request2.CreatedBy,
				Title:     request2.Title,
				Content:   request2.Content,
				Targets:   []*TargetOverview{},
				Tags:      []*TagResponse{},
				Group:     nil,
			},
			{
				ID:        request1.ID,
				Status:    request1.Status,
				CreatedAt: request1.CreatedAt,
				UpdatedAt: request1.UpdatedAt,
				CreatedBy: request1.CreatedBy,
				Title:     request1.Title,
				Content:   request1.Content,
				Targets:   []*TargetOverview{},
				Tags:      []*TagResponse{},
				Group:     nil,
			},
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		requests := []*model.RequestResponse{}

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/requests", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), model.RequestQuery{
				Limit:  100,
				Offset: 0,
			}).
			Return(requests, nil)

		require.NoError(t, h.Handlers.GetRequests(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*RequestResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("Success3", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date1 := time.Now()
		request1 := &model.RequestResponse{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date1,
			UpdatedAt: date1,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses:  []*model.RequestStatus{},
			Group:     nil,
		}
		requests := []*model.RequestResponse{request1}

		e := echo.New()
		status := "submitted"
		path := fmt.Sprintf("/api/requests?status=%s", status)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), model.RequestQuery{
				Status: &status,
				Limit:  100,
				Offset: 0,
			}).
			Return(requests, nil)

		require.NoError(t, h.Handlers.GetRequests(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*RequestResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := []*RequestResponse{
			{
				ID:        request1.ID,
				Status:    request1.Status,
				CreatedAt: request1.CreatedAt,
				UpdatedAt: request1.UpdatedAt,
				CreatedBy: request1.CreatedBy,
				Title:     request1.Title,
				Content:   request1.Content,
				Tags:      []*TagResponse{},
				Targets:   []*TargetOverview{},
				Group:     nil,
			},
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success4", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date1 := time.Now()
		date2str := date1.Add(time.Hour).Format("2006-01-02")
		date2, err := service.StrToDate(date2str)
		require.NoError(t, err)
		request1 := &model.RequestResponse{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date1,
			UpdatedAt: date1,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses:  []*model.RequestStatus{},
			Group:     nil,
		}
		requests := []*model.RequestResponse{request1}

		e := echo.New()
		path := fmt.Sprintf("/api/requests?until=%s", date2str)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), model.RequestQuery{
				Until:  &date2,
				Limit:  100,
				Offset: 0,
			}).
			Return(requests, nil)

		require.NoError(t, h.Handlers.GetRequests(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*RequestResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := []*RequestResponse{
			{
				ID:        request1.ID,
				Status:    request1.Status,
				CreatedAt: request1.CreatedAt,
				UpdatedAt: request1.UpdatedAt,
				CreatedBy: request1.CreatedBy,
				Title:     request1.Title,
				Content:   request1.Content,
				Tags:      []*TagResponse{},
				Targets:   []*TargetOverview{},
				Group:     nil,
			},
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success5", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date1 := time.Now()
		date2str := date1.Add(-time.Hour).Format("2006-01-02")
		date2, err := service.StrToDate(date2str)
		require.NoError(t, err)
		request1 := &model.RequestResponse{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date1,
			UpdatedAt: date1,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses:  []*model.RequestStatus{},
			Group:     nil,
		}
		requests := []*model.RequestResponse{request1}

		e := echo.New()
		path := fmt.Sprintf("/api/requests?since=%s", date2str)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), model.RequestQuery{
				Since:  &date2,
				Limit:  100,
				Offset: 0,
			}).
			Return(requests, nil)
		require.NoError(t, h.Handlers.GetRequests(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*RequestResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := []*RequestResponse{
			{
				ID:        request1.ID,
				Status:    request1.Status,
				CreatedAt: request1.CreatedAt,
				UpdatedAt: request1.UpdatedAt,
				CreatedBy: request1.CreatedBy,
				Title:     request1.Title,
				Content:   request1.Content,
				Targets:   []*TargetOverview{},
				Tags:      []*TagResponse{},
				Group:     nil,
			},
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success6", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date1 := time.Now()
		tag1 := model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 10),
			CreatedAt: date1,
			UpdatedAt: date1,
		}
		tag1ov := TagResponse{
			ID:        tag1.ID,
			Name:      tag1.Name,
			CreatedAt: tag1.CreatedAt,
			UpdatedAt: tag1.UpdatedAt,
		}
		request1 := &model.RequestResponse{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date1,
			UpdatedAt: date1,
			Tags:      []*model.Tag{&tag1},
			Targets:   []*model.RequestTargetDetail{},
			Statuses:  []*model.RequestStatus{},
			Group:     nil,
		}
		requests := []*model.RequestResponse{request1}

		e := echo.New()
		path := fmt.Sprintf("/api/requests?tag=%s", tag1.Name)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), model.RequestQuery{
				Tag:    &tag1.Name,
				Limit:  100,
				Offset: 0,
			}).
			Return(requests, nil)
		require.NoError(t, h.Handlers.GetRequests(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*RequestResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := []*RequestResponse{
			{
				ID:        request1.ID,
				Status:    request1.Status,
				CreatedAt: request1.CreatedAt,
				UpdatedAt: request1.UpdatedAt,
				CreatedBy: request1.CreatedBy,
				Title:     request1.Title,
				Content:   request1.Content,
				Tags:      []*TagResponse{&tag1ov},
				Targets:   []*TargetOverview{},
				Group:     nil,
			},
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success7", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		request := &model.RequestResponse{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses:  []*model.RequestStatus{},
			Group:     nil,
		}
		modelRequests := []*model.RequestResponse{request}

		e := echo.New()
		path := fmt.Sprintf("/api/requests?created_by=%s", request.CreatedBy.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), model.RequestQuery{
				Limit:     100,
				Offset:    0,
				CreatedBy: &request.CreatedBy},
			).
			Return(modelRequests, nil)

		err = h.Handlers.GetRequests(c)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*RequestResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := []*RequestResponse{
			{
				ID:        request.ID,
				Status:    request.Status,
				CreatedAt: date,
				UpdatedAt: date,
				CreatedBy: request.CreatedBy,
				Title:     request.Title,
				Content:   request.Content,
				Tags:      []*TagResponse{},
				Targets:   []*TargetOverview{},
				Group:     nil,
			},
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("InvaildStatus", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		e := echo.New()
		path := "/api/requests?status=invalid-status"
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.GetRequests(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, "invalid status"), err)
	})

	t.Run("FailedToGetRequests", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/requests", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		resErr := errors.New("Failed to get requests.")
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), model.RequestQuery{
				Limit:  100,
				Offset: 0,
			}).
			Return(nil, resErr)

		err = h.Handlers.GetRequests(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})
}

func TestHandlers_PostRequest(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				CreatedBy: uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqRequest := Request{
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		tags := []*model.Tag{}
		targets := []*model.RequestTarget{}
		var group *model.Group

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			CreateRequest(
				c.Request().Context(),
				reqRequest.Title, reqRequest.Content,
				tags, targets,
				group, reqRequest.CreatedBy).
			Return(request, nil)

		require.NoError(t, h.Handlers.PostRequest(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *RequestDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelRequestDetailToRequestResponse(request)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithTags", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		tags := []*model.Tag{tag}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
			Tags:      tags,
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				CreatedBy: uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqRequest := Request{
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
			Tags:      []*uuid.UUID{&tag.ID},
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		targets := []*model.RequestTarget{}
		var group *model.Group

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTag(c.Request().Context(), tag.ID).
			Return(tag, nil)
		h.Repository.MockRequestRepository.
			EXPECT().
			CreateRequest(
				c.Request().Context(),
				reqRequest.Title, reqRequest.Content,
				tags, targets,
				group, reqRequest.CreatedBy).
			Return(request, nil)

		require.NoError(t, h.Handlers.PostRequest(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *RequestDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelRequestDetailToRequestResponse(request)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithGroup", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		budget := random.Numeric(t, 100000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				CreatedBy: uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
			}},
			Group:    group,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqRequest := Request{
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
			Group:     &group.ID,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		tags := []*model.Tag{}
		targets := []*model.RequestTarget{}

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), group.ID).
			Return(group, nil)
		h.Repository.MockRequestRepository.
			EXPECT().
			CreateRequest(
				c.Request().Context(),
				reqRequest.Title, reqRequest.Content,
				tags, targets,
				group, reqRequest.CreatedBy).
			Return(request, nil)

		require.NoError(t, h.Handlers.PostRequest(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *RequestDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelRequestDetailToRequestResponse(request)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithTarget", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		target := &model.RequestTarget{
			Target: uuid.New(),
			Amount: random.Numeric(t, 1000000),
		}
		tgd := &model.RequestTargetDetail{
			ID:        uuid.New(),
			Target:    target.Target,
			Amount:    target.Amount,
			CreatedAt: date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{tgd},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				CreatedBy: uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		tg := &Target{
			Target: target.Target,
			Amount: target.Amount,
		}
		reqRequest := Request{
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
			Targets:   []*Target{tg},
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		tags := []*model.Tag{}
		var group *model.Group

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			CreateRequest(
				c.Request().Context(),
				reqRequest.Title, reqRequest.Content,
				tags, []*model.RequestTarget{target},
				group, reqRequest.CreatedBy).
			Return(request, nil)

		require.NoError(t, h.Handlers.PostRequest(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *RequestDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelRequestDetailToRequestResponse(request)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("UnknownTagID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}
		unknownTagID := uuid.New()
		reqRequest := Request{
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
			Tags:      []*uuid.UUID{&unknownTagID},
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTag(c.Request().Context(), unknownTagID).
			Return(nil, resErr)

		err = h.Handlers.PostRequest(c)
		require.Error(t, err)
		// FIXME: http.StatusNotFoundだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}
		unknownGroupID := uuid.New()
		reqRequest := Request{
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
			Group:     &unknownGroupID,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), unknownGroupID).
			Return(nil, resErr)

		err = h.Handlers.PostRequest(c)
		require.Error(t, err)
		// FIXME: http.StatusNotFoundだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
	})

	t.Run("UnknownUserID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}
		reqRequest := Request{
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		tags := []*model.Tag{}
		targets := []*model.RequestTarget{}
		var group *model.Group

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			CreateRequest(
				c.Request().Context(),
				reqRequest.Title, reqRequest.Content,
				tags, targets,
				group, reqRequest.CreatedBy).
			Return(nil, resErr)

		err = h.Handlers.PostRequest(c)
		require.Error(t, err)
		// FIXME: http.StatusNotFoundだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
	})
}

func TestHandlers_GetRequest(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(c.Request().Context(), request.ID).
			Return(request, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			GetComments(c.Request().Context(), request.ID).
			Return(nil, nil)

		require.NoError(t, h.Handlers.GetRequest(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *RequestDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelRequestDetailToRequestResponse(request)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithComments", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)
		date := time.Now()

		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		comment1 := &model.Comment{
			ID:        uuid.New(),
			User:      request.CreatedBy,
			Comment:   random.AlphaNumeric(t, 100),
			CreatedAt: date,
			UpdatedAt: date,
		}
		comment2 := &model.Comment{
			ID:        uuid.New(),
			User:      request.CreatedBy,
			Comment:   random.AlphaNumeric(t, 100),
			CreatedAt: date,
			UpdatedAt: date,
		}
		comments := []*model.Comment{comment1, comment2}

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(c.Request().Context(), request.ID).
			Return(request, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			GetComments(c.Request().Context(), request.ID).
			Return(comments, nil)

		require.NoError(t, h.Handlers.GetRequest(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *RequestDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelRequestDetailToRequestResponse(request)
		exp.Comments = lo.Map(comments, func(c *model.Comment, _ int) *CommentDetail {
			return modelCommentToCommentDetail(c)
		})
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithTarget", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		target := &model.RequestTargetDetail{
			ID:        uuid.New(),
			Target:    uuid.New(),
			Amount:    random.Numeric(t, 1000000),
			PaidAt:    nil,
			CreatedAt: date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{target},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(c.Request().Context(), request.ID).
			Return(request, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			GetComments(c.Request().Context(), request.ID).
			Return(nil, nil)
		require.NoError(t, h.Handlers.GetRequest(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *RequestDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelRequestDetailToRequestResponse(request)
		exp.Targets = []*TargetOverview{
			modelRequestTargetDetailToTargetOverview(target),
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", invalidUUID)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(invalidUUID)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.GetRequest(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", uuid.Nil)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(uuid.Nil.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		resErr := errors.New("invalid UUID")

		err = h.Handlers.GetRequest(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("UnknownID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		unknownID := uuid.New()

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", unknownID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(unknownID.String())

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(c.Request().Context(), unknownID).
			Return(nil, resErr)

		err = h.Handlers.GetRequest(c)
		require.Error(t, err)
		// FIXME: http.StatusNotFoundだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
	})
}

func TestHandlers_PutRequest(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqRequest := PutRequest{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
			Tags:    []*uuid.UUID{},
			Targets: []*Target{},
			Group:   nil,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		tags := []*model.Tag{}
		var group *model.Group
		updateRequest := &model.RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			CreatedBy: request.CreatedBy,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			Title:     reqRequest.Title,
			Content:   reqRequest.Content,
			Tags:      tags,
			Targets:   []*model.RequestTargetDetail{},
			Statuses:  request.Statuses,
			Group:     group,
			Comments:  request.Comments,
			Files:     request.Files,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", request.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		targets := lo.Map(
			updateRequest.Targets,
			func(t *model.RequestTargetDetail, _ int) *model.RequestTarget {
				return &model.RequestTarget{
					Target: t.Target,
					Amount: t.Amount,
				}
			},
		)
		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(c.Request().Context(), request.ID).
			Return(request, nil)
		h.Repository.MockRequestRepository.
			EXPECT().
			UpdateRequest(
				c.Request().Context(),
				request.ID,
				reqRequest.Title, reqRequest.Content,
				tags, targets,
				group).
			Return(updateRequest, nil)

		require.NoError(t, h.Handlers.PutRequest(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *RequestDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelRequestDetailToRequestResponse(updateRequest)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithTag", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		tag1 := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		tag2 := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		tags := []*model.Tag{tag1, tag2}
		reqRequest := PutRequest{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
			Tags:    []*uuid.UUID{&tag1.ID, &tag2.ID},
			Targets: []*Target{},
			Group:   nil,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		var group *model.Group
		updateRequest := &model.RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   reqRequest.Content,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			Tags:      tags,
			Targets:   []*model.RequestTargetDetail{},
			Statuses:  request.Statuses,
			Group:     group,
			Comments:  request.Comments,
			Files:     request.Files,
		}

		targets := lo.Map(
			updateRequest.Targets,
			func(t *model.RequestTargetDetail, _ int) *model.RequestTarget {
				return &model.RequestTarget{
					Target: t.Target,
					Amount: t.Amount,
				}
			},
		)
		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(c.Request().Context(), request.ID).
			Return(request, nil)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTag(c.Request().Context(), tag1.ID).
			Return(tag1, nil)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTag(c.Request().Context(), tag2.ID).
			Return(tag2, nil)
		h.Repository.MockRequestRepository.
			EXPECT().
			UpdateRequest(
				c.Request().Context(),
				request.ID,
				reqRequest.Title, reqRequest.Content,
				tags, targets,
				group).
			Return(updateRequest, nil)

		require.NoError(t, h.Handlers.PutRequest(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *RequestDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelRequestDetailToRequestResponse(updateRequest)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithTarget", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		target1 := &model.RequestTargetDetail{
			ID:        uuid.New(),
			Target:    uuid.New(),
			Amount:    random.Numeric(t, 100000),
			PaidAt:    nil,
			CreatedAt: date,
		}
		target2 := &model.RequestTargetDetail{
			ID:        uuid.New(),
			Target:    uuid.New(),
			Amount:    random.Numeric(t, 100000),
			PaidAt:    nil,
			CreatedAt: date,
		}
		targetDetails := []*model.RequestTargetDetail{target1, target2}
		reqRequest := PutRequest{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
			Tags:    []*uuid.UUID{},
			Targets: []*Target{
				{Target: target1.Target, Amount: target1.Amount},
				{Target: target2.Target, Amount: target2.Amount},
			},
			Group: nil,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		tags := []*model.Tag{}
		var group *model.Group
		updateRequest := &model.RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   reqRequest.Content,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			Tags:      tags,
			Targets:   targetDetails,
			Statuses:  request.Statuses,
			Group:     group,
			Comments:  request.Comments,
			Files:     request.Files,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		targets := lo.Map(
			updateRequest.Targets,
			func(t *model.RequestTargetDetail, _ int) *model.RequestTarget {
				return &model.RequestTarget{
					Target: t.Target,
					Amount: t.Amount,
				}
			},
		)
		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(c.Request().Context(), request.ID).
			Return(request, nil)
		h.Repository.MockRequestRepository.
			EXPECT().
			UpdateRequest(
				c.Request().Context(),
				request.ID,
				reqRequest.Title, reqRequest.Content,
				tags, targets,
				group).
			Return(updateRequest, nil)

		require.NoError(t, h.Handlers.PutRequest(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *RequestDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelRequestDetailToRequestResponse(updateRequest)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithGroup", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		budget := random.Numeric(t, 100000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		reqRequest := PutRequest{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
			Tags:    []*uuid.UUID{},
			Targets: []*Target{},
			Group:   &group.ID,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		tags := []*model.Tag{}
		updateRequest := &model.RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   reqRequest.Content,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			Tags:      tags,
			Targets:   []*model.RequestTargetDetail{},
			Statuses:  request.Statuses,
			Group:     group,
			Comments:  request.Comments,
			Files:     request.Files,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		targets := lo.Map(
			updateRequest.Targets,
			func(t *model.RequestTargetDetail, _ int) *model.RequestTarget {
				return &model.RequestTarget{
					Target: t.Target,
					Amount: t.Amount,
				}
			},
		)
		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(c.Request().Context(), request.ID).
			Return(request, nil)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), group.ID).
			Return(group, nil)
		h.Repository.MockRequestRepository.
			EXPECT().
			UpdateRequest(
				c.Request().Context(),
				request.ID,
				reqRequest.Title, reqRequest.Content,
				tags, targets,
				group).
			Return(updateRequest, nil)

		require.NoError(t, h.Handlers.PutRequest(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *RequestDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelRequestDetailToRequestResponse(updateRequest)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithComment", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqRequest := PutRequest{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
			Tags:    []*uuid.UUID{},
			Targets: []*Target{},
			Group:   nil,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		tags := []*model.Tag{}
		var group *model.Group
		comment1 := &model.Comment{
			ID:        uuid.New(),
			User:      request.CreatedBy,
			Comment:   random.AlphaNumeric(t, 100),
			CreatedAt: date,
			UpdatedAt: date,
		}
		comment2 := &model.Comment{
			ID:        uuid.New(),
			User:      request.CreatedBy,
			Comment:   random.AlphaNumeric(t, 100),
			CreatedAt: date,
			UpdatedAt: date,
		}
		comments := []*model.Comment{comment1, comment2}
		updateRequest := &model.RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			CreatedBy: request.CreatedBy,
			Title:     reqRequest.Title,
			Content:   reqRequest.Content,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			Tags:      tags,
			Targets:   []*model.RequestTargetDetail{},
			Statuses:  request.Statuses,
			Group:     nil,
			Comments:  comments,
			Files:     request.Files,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", request.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		targets := lo.Map(
			updateRequest.Targets,
			func(t *model.RequestTargetDetail, _ int) *model.RequestTarget {
				return &model.RequestTarget{
					Target: t.Target,
					Amount: t.Amount,
				}
			},
		)
		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(c.Request().Context(), request.ID).
			Return(request, nil)
		h.Repository.MockRequestRepository.
			EXPECT().
			UpdateRequest(
				c.Request().Context(),
				request.ID,
				reqRequest.Title, reqRequest.Content,
				tags, targets,
				group).
			Return(updateRequest, nil)

		require.NoError(t, h.Handlers.PutRequest(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *RequestDetailResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelRequestDetailToRequestResponse(updateRequest)
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", invalidUUID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(invalidUUID)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.PutRequest(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", uuid.Nil)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(uuid.Nil.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		resErr := errors.New("invalid UUID")

		err = h.Handlers.PutRequest(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("UnknownID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		unknownID := uuid.New()
		reqRequest := PutRequest{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
			Tags:    []*uuid.UUID{},
			Targets: []*Target{},
			Group:   nil,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", unknownID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(unknownID.String())
		c.Set(loginUserKey, user)

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(c.Request().Context(), unknownID).
			Return(nil, resErr)

		err = h.Handlers.PutRequest(c)
		require.Error(t, err)
		// FIXME: http.StatusNotFoundだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
	})

	t.Run("UnknownTagID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		reqRequest := PutRequest{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
			Tags:    []*uuid.UUID{&tag.ID},
			Targets: []*Target{},
			Group:   nil,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(c.Request().Context(), request.ID).
			Return(request, nil)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTag(c.Request().Context(), tag.ID).
			Return(nil, resErr)

		err = h.Handlers.PutRequest(c)
		require.Error(t, err)
		// FIXME: http.StatusNotFoundだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		budget := random.Numeric(t, 100000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		reqRequest := PutRequest{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
			Tags:    []*uuid.UUID{},
			Targets: []*Target{},
			Group:   &group.ID,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(c.Request().Context(), request.ID).
			Return(request, nil)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), group.ID).
			Return(nil, resErr)

		err = h.Handlers.PutRequest(c)
		require.Error(t, err)
		// FIXME: http.StatusNotFoundだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
	})
}

func TestHandlers_PutStatus(t *testing.T) {
	t.Parallel()

	t.Run("SuccessByCreator", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.FixRequired,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      user.ID,
			Comment:   reqStatus.Comment,
			CreatedAt: date,
			UpdatedAt: date,
		}
		status := &model.RequestStatus{
			ID:        uuid.New(),
			CreatedBy: user.ID,
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockRequestStatusRepository.
			EXPECT().
			CreateStatus(ctx, request.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, request.ID, user.ID).
			Return(comment, nil)

		require.NoError(t, h.Handlers.PutStatus(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *StatusResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &StatusResponse{
			CreatedBy: user.ID,
			Status:    status.Status,
			Comment: CommentDetail{
				ID:        comment.ID,
				User:      comment.User,
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			},
			CreatedAt: status.CreatedAt,
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessByAdminFromSubmittedToFixRequired", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.FixRequired,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      user.ID,
			Comment:   reqStatus.Comment,
			CreatedAt: date,
			UpdatedAt: date,
		}
		status := &model.RequestStatus{
			ID:        uuid.New(),
			CreatedBy: user.ID,
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockRequestStatusRepository.
			EXPECT().
			CreateStatus(ctx, request.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, request.ID, user.ID).
			Return(comment, nil)

		require.NoError(t, h.Handlers.PutStatus(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *StatusResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &StatusResponse{
			CreatedBy: user.ID,
			Status:    status.Status,
			Comment: CommentDetail{
				ID:        comment.ID,
				User:      comment.User,
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			},
			CreatedAt: status.CreatedAt,
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessByAdminFromSubmittedToAccepted", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.Accepted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      user.ID,
			Comment:   reqStatus.Comment,
			CreatedAt: date,
			UpdatedAt: date,
		}
		status := &model.RequestStatus{
			ID:        uuid.New(),
			CreatedBy: user.ID,
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockRequestStatusRepository.
			EXPECT().
			CreateStatus(ctx, request.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, request.ID, user.ID).
			Return(comment, nil)

		require.NoError(t, h.Handlers.PutStatus(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *StatusResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &StatusResponse{
			CreatedBy: user.ID,
			Status:    status.Status,
			Comment: CommentDetail{
				ID:        comment.ID,
				User:      comment.User,
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			},
			CreatedAt: status.CreatedAt,
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessByAdminFromSubmittedToFixRequired", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.FixRequired,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		comment := &model.Comment{
			ID:        uuid.New(),
			User:      user.ID,
			Comment:   reqStatus.Comment,
			CreatedAt: date,
			UpdatedAt: date,
		}
		status := &model.RequestStatus{
			ID:        uuid.New(),
			CreatedBy: user.ID,
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockRequestStatusRepository.
			EXPECT().
			CreateStatus(ctx, request.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, request.ID, user.ID).
			Return(comment, nil)

		require.NoError(t, h.Handlers.PutStatus(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *StatusResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &StatusResponse{
			CreatedBy: user.ID,
			Status:    status.Status,
			Comment: CommentDetail{
				ID:        comment.ID,
				User:      comment.User,
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			},
			CreatedAt: status.CreatedAt,
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessByAdminFromFixRequiredToSubmitted", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.FixRequired,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      user.ID,
			Comment:   reqStatus.Comment,
			CreatedAt: date,
			UpdatedAt: date,
		}
		status := &model.RequestStatus{
			ID:        uuid.New(),
			CreatedBy: user.ID,
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockRequestStatusRepository.
			EXPECT().
			CreateStatus(ctx, request.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, request.ID, user.ID).
			Return(comment, nil)

		require.NoError(t, h.Handlers.PutStatus(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *StatusResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &StatusResponse{
			CreatedBy: user.ID,
			Status:    status.Status,
			Comment: CommentDetail{
				ID:        comment.ID,
				User:      comment.User,
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			},
			CreatedAt: status.CreatedAt,
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessByAdminFromAcceptedToSubmitted", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Accepted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Accepted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		target := &model.RequestTargetDetail{
			ID:        uuid.New(),
			Target:    uuid.New(),
			PaidAt:    nil,
			CreatedAt: date,
		}
		targets := []*model.RequestTargetDetail{target}
		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      user.ID,
			Comment:   reqStatus.Comment,
			CreatedAt: date,
			UpdatedAt: date,
		}
		status := &model.RequestStatus{
			ID:        uuid.New(),
			CreatedBy: user.ID,
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockRequestTargetRepository.
			EXPECT().
			GetRequestTargets(ctx, request.ID).
			Return(targets, nil)
		h.Repository.MockRequestStatusRepository.
			EXPECT().
			CreateStatus(ctx, request.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, request.ID, user.ID).
			Return(comment, nil)

		require.NoError(t, h.Handlers.PutStatus(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *StatusResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &StatusResponse{
			CreatedBy: user.ID,
			Status:    status.Status,
			Comment: CommentDetail{
				ID:        comment.ID,
				User:      comment.User,
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt,
				UpdatedAt: comment.UpdatedAt,
			},
			CreatedAt: status.CreatedAt,
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("InvalidStatus", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.FixRequired,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		invalidStatus := random.AlphaNumeric(t, 20)
		reqBody, err := json.Marshal(&struct {
			Status  string `json:"status"`
			Comment string `json:"comment"`
		}{
			Status:  invalidStatus,
			Comment: random.AlphaNumeric(t, 20),
		})
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		resErr := echo.NewHTTPError(http.StatusBadRequest)
		resErrMessage := echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("invalid Status %s", invalidStatus))
		resErrMessage.Internal = fmt.Errorf("invalid Status %s", invalidStatus)
		resErr.Message = resErrMessage

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい
		require.Equal(t, resErr, err)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		invalidUUID := "invalid-uuid"
		_, resErr := uuid.Parse(invalidUUID)
		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", invalidUUID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(invalidUUID)
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("NillUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", uuid.Nil)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(uuid.Nil.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		_, resErr := uuid.Parse(c.Param("requestID"))

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("SameStatusError", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Status(random.Numeric(t, 5) + 1),
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Status(random.Numeric(t, 5) + 1),
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  request.Status,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)

		resErr := errors.New("invalid request: same status")

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("CommentRequiredErrorFromSubmittedToFixRequired", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)
		date := time.Now()

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}

		reqStatus := PutStatus{
			Status: model.FixRequired,
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)

		resErr := fmt.Errorf(
			"unable to change %v to %v without comment",
			request.Status.String(),
			reqStatus.Status.String())

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("CommentRequiredErrorFromSubmittedToRejected", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status: model.Rejected,
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)

		resErr := fmt.Errorf(
			"unable to change %v to %v without comment",
			request.Status.String(),
			reqStatus.Status.String())

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("CommentRequiredErrorFromAcceptedToSubmitted", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Accepted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Accepted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status: model.Submitted,
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)

		resErr := fmt.Errorf(
			"unable to change %v to %v without comment",
			request.Status.String(),
			reqStatus.Status.String())

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("AdminNoPrivilege", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.FixRequired,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.Accepted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)

		resErr := fmt.Errorf(
			"admin unable to change %v to %v",
			request.Status.String(),
			reqStatus.Status.String())

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusForbiddenだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("AlreadyPaid", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, true)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Accepted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Accepted,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		target := &model.RequestTargetDetail{
			ID:        uuid.New(),
			Target:    uuid.New(),
			PaidAt:    &date,
			CreatedAt: date,
		}
		targets := []*model.RequestTargetDetail{target}
		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		ctx = c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockRequestTargetRepository.
			EXPECT().
			GetRequestTargets(ctx, request.ID).
			Return(targets, nil)

		resErr := errors.New("someone already paid")

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("CreatorNoPrivilege", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: user.ID,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.Submitted,
				CreatedAt: date,
				CreatedBy: user.ID,
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.Accepted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)

		resErr := fmt.Errorf(
			"creator unable to change %v to %v",
			request.Status.String(), reqStatus.Status.String())

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusForbiddenだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})

	t.Run("NoPrivilege", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, false)
		user := userFromModelUser(*accessUser)
		date := time.Now()
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			Tags:      []*model.Tag{},
			Targets:   []*model.RequestTargetDetail{},
			Statuses: []*model.RequestStatus{{
				ID:        uuid.New(),
				Status:    model.FixRequired,
				CreatedAt: date,
				CreatedBy: uuid.New(),
			}},
			Group:    nil,
			Comments: []*model.Comment{},
			Files:    []*uuid.UUID{},
		}
		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/requests/%s/status", request.ID.String())
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		ctx = c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)

		err = h.Handlers.PutStatus(c)
		require.Error(t, err)
		// FIXME: http.StatusForbiddenだけ判定したい
		require.Equal(
			t,
			echo.NewHTTPError(http.StatusForbidden, "you are not request creator"),
			err)
	})
}

// TODO: TestHandlers_PostComment
