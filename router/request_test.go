package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/sessions"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
	"github.com/traPtitech/Jomon/testutil/random"
)

// To do
func TestHandlers_GetRequests(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
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
		}
		request2 := &model.RequestResponse{
			ID:        uuid.New(),
			Status:    model.Submitted,
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date2,
			UpdatedAt: date2,
		}
		requests := []*model.RequestResponse{request2, request1}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/requests", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), model.RequestQuery{}).
			Return(requests, nil)

		res := []*RequestResponse{
			{
				ID:        request2.ID,
				Status:    request2.Status,
				CreatedAt: request2.CreatedAt,
				UpdatedAt: request2.UpdatedAt,
				CreatedBy: request2.CreatedBy,
				Title:     request2.Title,
				Content:   request2.Content,
				Targets:   []*TargetOverview{},
				Comments:  []*CommentDetail{},
				Tags:      []*TagOverview{},
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
				Comments:  []*CommentDetail{},
				Tags:      []*TagOverview{},
			},
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetRequests(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		requests := []*model.RequestResponse{}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/requests", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), model.RequestQuery{}).
			Return(requests, nil)

		res := []*RequestResponse{}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetRequests(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("Success3", func(t *testing.T) {
		t.Parallel()
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
		}
		requests := []*model.RequestResponse{request1}

		status := "submitted"

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/requests?status=submitted", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), model.RequestQuery{
				Status: &status,
			}).
			Return(requests, nil)

		res := []*RequestResponse{
			{
				ID:        request1.ID,
				Status:    request1.Status,
				CreatedAt: request1.CreatedAt,
				UpdatedAt: request1.UpdatedAt,
				CreatedBy: request1.CreatedBy,
				Title:     request1.Title,
				Content:   request1.Content,
				Targets:   []*TargetOverview{},
				Comments:  []*CommentDetail{},
				Tags:      []*TagOverview{},
			},
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetRequests(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("Success4", func(t *testing.T) {
		t.Parallel()
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
		}
		requests := []*model.RequestResponse{request1}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/requests?until=%s", date2str), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), model.RequestQuery{
				Until: &date2,
			}).
			Return(requests, nil)

		res := []*RequestResponse{
			{
				ID:        request1.ID,
				Status:    request1.Status,
				CreatedAt: request1.CreatedAt,
				UpdatedAt: request1.UpdatedAt,
				CreatedBy: request1.CreatedBy,
				Title:     request1.Title,
				Content:   request1.Content,
				Targets:   []*TargetOverview{},
				Comments:  []*CommentDetail{},
				Tags:      []*TagOverview{},
			},
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetRequests(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("Success5", func(t *testing.T) {
		t.Parallel()
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
		}
		requests := []*model.RequestResponse{request1}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/requests?since=%s", date2str), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), model.RequestQuery{
				Since: &date2,
			}).
			Return(requests, nil)

		res := []*RequestResponse{
			{
				ID:        request1.ID,
				Status:    request1.Status,
				CreatedAt: request1.CreatedAt,
				UpdatedAt: request1.UpdatedAt,
				CreatedBy: request1.CreatedBy,
				Title:     request1.Title,
				Content:   request1.Content,
				Targets:   []*TargetOverview{},
				Comments:  []*CommentDetail{},
				Tags:      []*TagOverview{},
			},
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetRequests(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("Success6", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		date1 := time.Now()

		tag1 := model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 10),
			CreatedAt: date1,
			UpdatedAt: date1,
		}

		tag1ov := TagOverview{
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
			Tags:      []*model.Tag{&tag1},
			CreatedAt: date1,
			UpdatedAt: date1,
		}
		requests := []*model.RequestResponse{request1}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/requests?tag=%s", tag1.Name), nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), model.RequestQuery{
				Tag: &tag1.Name,
			}).
			Return(requests, nil)

		res := []*RequestResponse{
			{
				ID:        request1.ID,
				Status:    request1.Status,
				CreatedAt: request1.CreatedAt,
				UpdatedAt: request1.UpdatedAt,
				CreatedBy: request1.CreatedBy,
				Title:     request1.Title,
				Tags:      []*TagOverview{&tag1ov},
				Content:   request1.Content,
				Targets:   []*TargetOverview{},
				Comments:  []*CommentDetail{},
			},
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetRequests(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("InvaildStatus", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/requests?status=po", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.GetRequests(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, "invalid status"), err)
		}
	})

	t.Run("FailedToGetRequests", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/requests", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		resErr := errors.New("Failed to get requests.")
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), model.RequestQuery{}).
			Return(nil, resErr)

		err = h.Handlers.GetRequests(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})
}

func TestHandlers_PostRequest(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
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
		req, err := http.NewRequest(http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			CreateRequest(c.Request().Context(), reqRequest.Title, reqRequest.Content, tags, targets, group, reqRequest.CreatedBy).
			Return(request, nil)

		res := &RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PostRequest(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessWithTags", func(t *testing.T) {
		t.Parallel()
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
			Tags:      tags,
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
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
		req, err := http.NewRequest(http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		assert.NoError(t, err)
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
			CreateRequest(c.Request().Context(), reqRequest.Title, reqRequest.Content, tags, targets, group, reqRequest.CreatedBy).
			Return(request, nil)

		res := &RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
			Tags: []*TagOverview{{
				ID:        tag.ID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt,
				UpdatedAt: tag.UpdatedAt,
			}},
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PostRequest(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessWithGroup", func(t *testing.T) {
		t.Parallel()
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
			Group:     group,
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
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
		req, err := http.NewRequest(http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		assert.NoError(t, err)
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
			CreateRequest(c.Request().Context(), reqRequest.Title, reqRequest.Content, tags, targets, group, reqRequest.CreatedBy).
			Return(request, nil)

		res := &RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
			Group: &GroupOverview{
				ID:          group.ID,
				Name:        group.Name,
				Description: group.Description,
				Budget:      group.Budget,
				CreatedAt:   group.CreatedAt,
				UpdatedAt:   group.UpdatedAt,
			},
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PostRequest(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessWithTarget", func(t *testing.T) {
		t.Parallel()
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
			Targets:   []*model.RequestTargetDetail{tgd},
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
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
		req, err := http.NewRequest(http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			CreateRequest(c.Request().Context(), reqRequest.Title, reqRequest.Content, tags, []*model.RequestTarget{target}, group, reqRequest.CreatedBy).
			Return(request, nil)

		tgov := &TargetOverview{
			ID:        request.Targets[0].ID,
			Target:    request.Targets[0].Target,
			Amount:    request.Targets[0].Amount,
			CreatedAt: request.Targets[0].CreatedAt,
		}

		res := &RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
			Targets:   []*TargetOverview{tgov},
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PostRequest(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("UnknownTagID", func(t *testing.T) {
		t.Parallel()
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
		req, err := http.NewRequest(http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		assert.NoError(t, err)
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
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
		}
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
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
		req, err := http.NewRequest(http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		assert.NoError(t, err)
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
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
		}
	})

	t.Run("UnknownUserID", func(t *testing.T) {
		t.Parallel()
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
		req, err := http.NewRequest(http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			CreateRequest(c.Request().Context(), reqRequest.Title, reqRequest.Content, tags, targets, group, reqRequest.CreatedBy).
			Return(nil, resErr)

		err = h.Handlers.PostRequest(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
		}
	})
}

func TestHandlers_GetRequest(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			Comments:  []*model.Comment{},
			Files:     []*uuid.UUID{},
			Statuses:  []*model.RequestStatus{},
			Tags:      []*model.Tag{},
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/requests/%s", request.ID.String()), nil)
		assert.NoError(t, err)
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

		res := &RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
			Tags:      []*TagOverview{},
			Comments:  []*CommentDetail{},
			Targets:   []*TargetOverview{},
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetRequest(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessWithComments", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			Comments:  []*model.Comment{},
			Files:     []*uuid.UUID{},
			Statuses:  []*model.RequestStatus{},
			Tags:      []*model.Tag{},
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
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
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/requests/%s", request.ID.String()), nil)
		assert.NoError(t, err)
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

		resComments := []*CommentDetail{
			{
				ID:        comment1.ID,
				User:      comment1.User,
				Comment:   comment1.Comment,
				CreatedAt: comment1.CreatedAt,
				UpdatedAt: comment1.UpdatedAt,
			},
			{
				ID:        comment2.ID,
				User:      comment2.User,
				Comment:   comment2.Comment,
				CreatedAt: comment2.CreatedAt,
				UpdatedAt: comment2.UpdatedAt,
			},
		}

		res := &RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
			Comments:  resComments,
			Tags:      []*TagOverview{},
			Targets:   []*TargetOverview{},
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetRequest(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessWithTarget", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		target := &TargetOverview{
			ID:        uuid.New(),
			Target:    uuid.New(),
			Amount:    random.Numeric(t, 1000000),
			PaidAt:    nil,
			CreatedAt: date,
		}

		modeltarget := &model.RequestTargetDetail{
			ID:        target.ID,
			Target:    target.Target,
			Amount:    target.Amount,
			PaidAt:    target.PaidAt,
			CreatedAt: target.CreatedAt,
		}

		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			Comments:  []*model.Comment{},
			Files:     []*uuid.UUID{},
			Statuses:  []*model.RequestStatus{},
			Tags:      []*model.Tag{},
			Content:   random.AlphaNumeric(t, 50),
			Targets:   []*model.RequestTargetDetail{modeltarget},
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}
		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/requests/%s", request.ID.String()), nil)
		assert.NoError(t, err)
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

		res := &RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   request.Content,
			Targets:   []*TargetOverview{target},
			Tags:      []*TagOverview{},
			Comments:  []*CommentDetail{},
		}

		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetRequest(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/requests/hoge", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues("hoge")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		_, resErr := uuid.Parse(c.Param("requestID"))

		err = h.Handlers.GetRequest(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/requests/%s", uuid.Nil), nil)
		assert.NoError(t, err)
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
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("UnknownID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		unknownID := uuid.New()

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/requests/%s", unknownID.String()), nil)
		assert.NoError(t, err)
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
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
		}
	})
}

func TestHandlers_PutRequest(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
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
		}

		reqRequest := PutRequest{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		tags := []*model.Tag{}
		targets := []*model.RequestTarget{}
		var group *model.Group

		updateRequest := &model.RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: request.CreatedBy,
			Title:     reqRequest.Title,
			Content:   reqRequest.Content,
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s", request.ID), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			UpdateRequest(c.Request().Context(), request.ID, reqRequest.Title, reqRequest.Content, tags, targets, group).
			Return(updateRequest, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			GetComments(c.Request().Context(), request.ID).
			Return([]*model.Comment{}, nil)

		res := &RequestResponse{
			ID:        updateRequest.ID,
			Status:    updateRequest.Status,
			CreatedAt: updateRequest.CreatedAt,
			UpdatedAt: updateRequest.UpdatedAt,
			CreatedBy: updateRequest.CreatedBy,
			Title:     updateRequest.Title,
			Content:   updateRequest.Content,
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PutRequest(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessWithTag", func(t *testing.T) {
		t.Parallel()
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
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		targets := []*model.RequestTarget{}
		var group *model.Group

		updateRequest := &model.RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   reqRequest.Content,
			Tags:      tags,
			Group:     group,
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

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
			UpdateRequest(c.Request().Context(), request.ID, reqRequest.Title, reqRequest.Content, tags, targets, group).
			Return(updateRequest, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			GetComments(c.Request().Context(), request.ID).
			Return([]*model.Comment{}, nil)

		res := &RequestResponse{
			ID:        updateRequest.ID,
			Status:    updateRequest.Status,
			CreatedAt: updateRequest.CreatedAt,
			UpdatedAt: updateRequest.UpdatedAt,
			CreatedBy: updateRequest.CreatedBy,
			Title:     updateRequest.Title,
			Content:   updateRequest.Content,
			Tags: []*TagOverview{
				{
					ID:        tag1.ID,
					Name:      tag1.Name,
					CreatedAt: tag1.CreatedAt,
					UpdatedAt: tag1.UpdatedAt,
				},
				{
					ID:        tag2.ID,
					Name:      tag2.Name,
					CreatedAt: tag2.CreatedAt,
					UpdatedAt: tag2.UpdatedAt,
				},
			},
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PutRequest(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessWithTarget", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
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
		targets := []*model.RequestTarget{
			{
				Target: target1.Target,
				Amount: target1.Amount,
			},
			{
				Target: target2.Target,
				Amount: target2.Amount,
			},
		}
		targetDetails := []*model.RequestTargetDetail{target1, target2}

		reqRequest := PutRequest{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
			Targets: []*Target{
				{
					Target: target1.Target,
					Amount: target1.Amount,
				},
				{
					Target: target2.Target,
					Amount: target2.Amount,
				},
			},
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		tags := []*model.Tag{}
		var group *model.Group

		updateRequest := &model.RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Tags:      tags,
			Targets:   targetDetails,
			Group:     group,
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		h.Repository.MockRequestRepository.
			EXPECT().
			UpdateRequest(c.Request().Context(), request.ID, reqRequest.Title, reqRequest.Content, tags, targets, group).
			Return(updateRequest, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			GetComments(c.Request().Context(), request.ID).
			Return([]*model.Comment{}, nil)

		res := &RequestResponse{
			ID:        updateRequest.ID,
			Status:    updateRequest.Status,
			CreatedAt: updateRequest.CreatedAt,
			UpdatedAt: updateRequest.UpdatedAt,
			CreatedBy: updateRequest.CreatedBy,
			Title:     updateRequest.Title,
			Targets: []*TargetOverview{
				{
					ID:        target1.ID,
					Target:    target1.Target,
					Amount:    target1.Amount,
					PaidAt:    target1.PaidAt,
					CreatedAt: target1.CreatedAt,
				},
				{
					ID:        target2.ID,
					Target:    target2.Target,
					Amount:    target2.Amount,
					PaidAt:    target2.PaidAt,
					CreatedAt: target2.CreatedAt,
				},
			},
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PutRequest(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessWithGroup", func(t *testing.T) {
		t.Parallel()
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
			Group:   &group.ID,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		tags := []*model.Tag{}
		targets := []*model.RequestTarget{}

		updateRequest := &model.RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: request.CreatedBy,
			Title:     request.Title,
			Content:   reqRequest.Content,
			Tags:      tags,
			Group:     group,
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), group.ID).
			Return(group, nil)
		h.Repository.MockRequestRepository.
			EXPECT().
			UpdateRequest(c.Request().Context(), request.ID, reqRequest.Title, reqRequest.Content, tags, targets, group).
			Return(updateRequest, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			GetComments(c.Request().Context(), request.ID).
			Return([]*model.Comment{}, nil)

		res := &RequestResponse{
			ID:        updateRequest.ID,
			Status:    updateRequest.Status,
			CreatedAt: updateRequest.CreatedAt,
			UpdatedAt: updateRequest.UpdatedAt,
			CreatedBy: updateRequest.CreatedBy,
			Title:     updateRequest.Title,
			Content:   updateRequest.Content,
			Group: &GroupOverview{
				ID:          group.ID,
				Name:        group.Name,
				Description: group.Description,
				Budget:      group.Budget,
				CreatedAt:   group.CreatedAt,
				UpdatedAt:   group.UpdatedAt,
			},
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PutRequest(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessWithComment", func(t *testing.T) {
		t.Parallel()
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

		reqRequest := PutRequest{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		tags := []*model.Tag{}
		targets := []*model.RequestTarget{}
		var group *model.Group

		updateRequest := &model.RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: request.CreatedBy,
			Title:     reqRequest.Title,
			Content:   reqRequest.Content,
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
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s", request.ID), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			UpdateRequest(c.Request().Context(), request.ID, reqRequest.Title, reqRequest.Content, tags, targets, group).
			Return(updateRequest, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			GetComments(c.Request().Context(), request.ID).
			Return(comments, nil)

		resComments := []*CommentDetail{
			{
				ID:        comment1.ID,
				User:      comment1.User,
				Comment:   comment1.Comment,
				CreatedAt: comment1.CreatedAt,
				UpdatedAt: comment1.UpdatedAt,
			},
			{
				ID:        comment2.ID,
				User:      comment2.User,
				Comment:   comment2.Comment,
				CreatedAt: comment2.CreatedAt,
				UpdatedAt: comment2.UpdatedAt,
			},
		}

		res := &RequestResponse{
			ID:        updateRequest.ID,
			Status:    updateRequest.Status,
			CreatedAt: updateRequest.CreatedAt,
			UpdatedAt: updateRequest.UpdatedAt,
			CreatedBy: updateRequest.CreatedBy,
			Title:     updateRequest.Title,
			Content:   updateRequest.Content,
			Comments:  resComments,
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PutRequest(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, "/api/requests/hoge", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues("hoge")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		_, resErr := uuid.Parse(c.Param("requestID"))

		err = h.Handlers.PutRequest(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s", uuid.Nil), nil)
		assert.NoError(t, err)
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
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("UnknownID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		unknownID := uuid.New()
		reqRequest := PutRequest{
			Title:   random.AlphaNumeric(t, 30),
			Content: random.AlphaNumeric(t, 50),
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		tags := []*model.Tag{}
		targets := []*model.RequestTarget{}
		var group *model.Group

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s", unknownID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
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
			UpdateRequest(c.Request().Context(), unknownID, reqRequest.Title, reqRequest.Content, tags, targets, group).
			Return(nil, resErr)

		err = h.Handlers.PutRequest(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
		}
	})

	t.Run("UnknownTagID", func(t *testing.T) {
		t.Parallel()
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
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTag(c.Request().Context(), tag.ID).
			Return(nil, resErr)

		err = h.Handlers.PutRequest(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
		}
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
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
			Group:   &group.ID,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/requests/:requestID")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), group.ID).
			Return(nil, resErr)

		err = h.Handlers.PutRequest(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
		}
	})
}

func TestHandlers_PutStatus(t *testing.T) {
	t.Parallel()

	t.Run("SuccessByCreator", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       false,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: user.ID,
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
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		ctx := c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil)
		h.Repository.MockRequestStatusRepository.
			EXPECT().
			CreateStatus(ctx, request.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, request.ID, user.ID).
			Return(comment, nil)

		res := &StatusResponse{
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
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PutStatus(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessByAdminFromSubmittedToFixRequired", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
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
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		ctx := c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil)
		h.Repository.MockRequestStatusRepository.
			EXPECT().
			CreateStatus(ctx, request.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, request.ID, user.ID).
			Return(comment, nil)

		res := &StatusResponse{
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
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PutStatus(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessByAdminFromSubmittedToAccepted", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
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
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		ctx := c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil)
		h.Repository.MockRequestStatusRepository.
			EXPECT().
			CreateStatus(ctx, request.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, request.ID, user.ID).
			Return(comment, nil)

		res := &StatusResponse{
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
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PutStatus(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessByAdminFromSubmittedToFixRequired", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
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
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		ctx := c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil)
		h.Repository.MockRequestStatusRepository.
			EXPECT().
			CreateStatus(ctx, request.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, request.ID, user.ID).
			Return(comment, nil)

		res := &StatusResponse{
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
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PutStatus(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessByAdminFromFixRequiredToSubmitted", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
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
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		ctx := c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil)
		h.Repository.MockRequestStatusRepository.
			EXPECT().
			CreateStatus(ctx, request.ID, user.ID, reqStatus.Status).
			Return(status, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(ctx, reqStatus.Comment, request.ID, user.ID).
			Return(comment, nil)

		res := &StatusResponse{
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
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PutStatus(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessByAdminFromAcceptedToSubmitted", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Accepted,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
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
			Status:    reqStatus.Status,
			CreatedAt: date,
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		ctx := c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil)
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

		res := &StatusResponse{
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
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PutStatus(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("InvalidStatus", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       false,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: user.ID,
		}

		invalidStatus := random.AlphaNumeric(t, 20)
		reqStatus := fmt.Sprintf(`
		{
			"status": "%s",
			"comment": "%s"
		}`, invalidStatus, random.AlphaNumeric(t, 20))

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), strings.NewReader(reqStatus))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		resErr := echo.NewHTTPError(http.StatusBadRequest)
		resErrMessage := echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid Status %s", invalidStatus))
		resErrMessage.Internal = fmt.Errorf("invalid Status %s", invalidStatus)
		resErr.Message = resErrMessage

		err = h.Handlers.PutStatus(c)
		if assert.Error(t, err) {
			assert.Equal(t, resErr, err)
		}
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       false,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", "hoge"), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues("hoge")
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		_, resErr := uuid.Parse(c.Param("requestID"))

		err = h.Handlers.PutStatus(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("NillUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       false,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", uuid.Nil), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(uuid.Nil.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		_, resErr := uuid.Parse(c.Param("requestID"))

		err = h.Handlers.PutStatus(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("SissionNotFound", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       false,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: user.ID,
		}

		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		resErr := errors.New("sessionUser not found")

		err = h.Handlers.PutStatus(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusForbidden, resErr), err)
		}
	})

	t.Run("SameStatusError", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       false,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Status(random.Numeric(t, 5) + 1),
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: user.ID,
		}

		reqStatus := PutStatus{
			Status:  request.Status,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		ctx := c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)

		resErr := errors.New("invalid request: same status")

		err = h.Handlers.PutStatus(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("CommentRequiredErrorFromSubmittedToFixRequired", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       false,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: user.ID,
		}

		reqStatus := PutStatus{
			Status: model.FixRequired,
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		ctx := c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)

		resErr := fmt.Errorf("unable to change %v to %v without comment", request.Status.String(), reqStatus.Status.String())

		err = h.Handlers.PutStatus(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("CommentRequiredErrorFromSubmittedToRejected", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       false,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: user.ID,
		}

		reqStatus := PutStatus{
			Status: model.Rejected,
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		ctx := c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)

		resErr := fmt.Errorf("unable to change %v to %v without comment", request.Status.String(), reqStatus.Status.String())

		err = h.Handlers.PutStatus(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("CommentRequiredErrorFromAcceptedToSubmitted", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       false,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Accepted,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: user.ID,
		}

		reqStatus := PutStatus{
			Status: model.Submitted,
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		ctx := c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)

		resErr := fmt.Errorf("unable to change %v to %v without comment", request.Status.String(), reqStatus.Status.String())

		err = h.Handlers.PutStatus(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("UnknownUser", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       false,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: user.ID,
		}

		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		ctx := c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByID(ctx, user.ID).
			Return(nil, resErr)

		err = h.Handlers.PutStatus(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusNotFound, resErr), err)
		}
	})

	t.Run("AdminNoPrivilege", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}

		reqStatus := PutStatus{
			Status:  model.Accepted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		ctx := c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil)

		resErr := fmt.Errorf("admin unable to change %v to %v", request.Status.String(), reqStatus.Status.String())

		err = h.Handlers.PutStatus(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("AlreadyPaid", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Accepted,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
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
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		ctx := c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil)
		h.Repository.MockRequestTargetRepository.
			EXPECT().
			GetRequestTargets(ctx, request.ID).
			Return(targets, nil)

		resErr := errors.New("someone already paid")

		err = h.Handlers.PutStatus(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("CreatorNoPrivilege", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       false,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: user.ID,
		}

		reqStatus := PutStatus{
			Status:  model.Accepted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		ctx := c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil)

		resErr := fmt.Errorf("creator unable to change %v to %v", request.Status.String(), reqStatus.Status.String())

		err = h.Handlers.PutStatus(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("NoPrivilege", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       false,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.FixRequired,
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}

		reqStatus := PutStatus{
			Status:  model.Submitted,
			Comment: random.AlphaNumeric(t, 20),
		}
		reqBody, err := json.Marshal(reqStatus)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/requests/%s/status", request.ID.String()), bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/requests/:requestID/status")
		c.SetParamNames("requestID")
		c.SetParamValues(request.ID.String())
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		assert.NoError(t, err)
		sess.Values[sessionUserKey] = User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.Name,
			Admin:       user.Admin,
		}

		ctx := c.Request().Context()
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequest(ctx, request.ID).
			Return(request, nil)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByID(ctx, user.ID).
			Return(user, nil)

		err = h.Handlers.PutStatus(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusForbidden), err)
		}
	})
}

// func TestHandlers_PostComment(t *testing.T) {
// 	t.Parallel()

// 	t.Run("Success", func(t *testing.T) {
// 		t.Parallel()
// 		ctrl := gomock.NewController(t)

// 		date := time.Now()

// 		comment := &model.Comment{
// 			ID:        uuid.New(),
// 			User:      uuid.New(),
// 			Comment:   random.AlphaNumeric(t, 20),
// 			CreatedAt: date,
// 			UpdatedAt: date,
// 		}

// 		requestID := uuid.New()
// 		reqComment := Comment{
// 			Comment: comment.Comment,
// 		}
// 		reqBody, err := json.Marshal(reqComment)
// 		require.NoError(t, err)

// 		e := echo.New()
// 		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/requests/%s/comments", requestID), bytes.NewReader(reqBody))
// 		assert.NoError(t, err)
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		c.SetPath("api/requests/:requestID/comments")
// 		c.SetParamNames("requestID")
// 		c.SetParamValues(requestID.String())

// 		h, err := NewTestHandlers(t, ctrl)
// 		assert.NoError(t, err)

// 		h.Repository.MockCommentRepository.
// 			EXPECT().
// 			CreateComment(c.Request().Context(), comment.Comment, requestID, comment.User).
// 			Return(comment, nil)

// 		res := &CommentDetail{
// 			ID:        comment.ID,
// 			User:      comment.User,
// 			Comment:   comment.Comment,
// 			CreatedAt: comment.CreatedAt,
// 			UpdatedAt: comment.UpdatedAt,
// 		}
// 		resBody, err := json.Marshal(res)
// 		require.NoError(t, err)

// 		if assert.NoError(t, h.Handlers.PostComment(c)) {
// 			assert.Equal(t, http.StatusOK, rec.Code)
// 			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
// 		}
// 	})
// }
