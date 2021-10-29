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

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/model"
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
			Status:    model.Submitted.String(),
			CreatedBy: uuid.New(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date1,
			UpdatedAt: date1,
		}
		request2 := &model.RequestResponse{
			ID:        uuid.New(),
			Status:    model.Submitted.String(),
			CreatedBy: uuid.New(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date2,
			UpdatedAt: date2,
		}
		requests := []*model.RequestResponse{request2, request1}

		var sort string
		var target string
		var year int
		var since time.Time
		var until time.Time
		var tag string
		var group string

		query := model.RequestQuery{
			Sort:   &sort,
			Target: &target,
			Year:   &year,
			Since:  &since,
			Until:  &until,
			Tag:    &tag,
			Group:  &group,
		}

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
			GetRequests(c.Request().Context(), query).
			Return(requests, nil)

		res := []*RequestResponse{
			{
				ID:        request2.ID,
				Status:    request2.Status,
				CreatedAt: request2.CreatedAt,
				UpdatedAt: request2.UpdatedAt,
				CreatedBy: request2.CreatedBy,
				Amount:    request2.Amount,
				Title:     request2.Title,
				Content:   request2.Content,
			},
			{
				ID:        request1.ID,
				Status:    request1.Status,
				CreatedAt: request1.CreatedAt,
				UpdatedAt: request1.UpdatedAt,
				CreatedBy: request1.CreatedBy,
				Amount:    request1.Amount,
				Title:     request1.Title,
				Content:   request1.Content,
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

		var title string
		var target string
		var year int
		var since time.Time
		var until time.Time
		var tag string
		var group string

		query := model.RequestQuery{
			Sort:   &title,
			Target: &target,
			Year:   &year,
			Since:  &since,
			Until:  &until,
			Tag:    &tag,
			Group:  &group,
		}

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
			GetRequests(c.Request().Context(), query).
			Return(requests, nil)

		var res []*RequestResponse
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetRequests(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("FailedToGetRequests", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		var title string
		var target string
		var year int
		var since time.Time
		var until time.Time
		var tag string
		var group string

		query := model.RequestQuery{
			Sort:   &title,
			Target: &target,
			Year:   &year,
			Since:  &since,
			Until:  &until,
			Tag:    &tag,
			Group:  &group,
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/requests", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		mocErr := errors.New("Failed to get requests.")
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), query).
			Return(nil, mocErr)
		
	  err = h.Handlers.GetRequests(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})
}

//group, file, comment考えてない
func TestHandlers_PostRequest(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		date := time.Now()

		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}

		reqRequest := Request{
			CreatedBy: request.CreatedBy,
			Amount:    request.Amount,
			Title:     request.Title,
			Content:   request.Content,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		var tags []*model.Tag
		var group *model.Group

		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		h.Repository.MockRequestRepository.
			EXPECT().
			CreateRequest(c.Request().Context(), reqRequest.Amount, reqRequest.Title, reqRequest.Content, tags, group, reqRequest.CreatedBy).
			Return(request, nil)

		res := &RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Amount:    request.Amount,
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
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		tags := []*model.Tag{tag}

		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			Tags:      tags,
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}

		reqRequest := Request{
			CreatedBy: request.CreatedBy,
			Amount:    request.Amount,
			Title:     request.Title,
			Content:   request.Content,
			Tags:      []*uuid.UUID{&tag.ID},
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		var group *model.Group

		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, "/api/requests", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		h.Repository.MockTagRepository.
			EXPECT().
			GetTag(c.Request().Context(), tag.ID).
			Return(tag, nil)
		h.Repository.MockRequestRepository.
			EXPECT().
			CreateRequest(c.Request().Context(), reqRequest.Amount, reqRequest.Title, reqRequest.Content, tags, group, reqRequest.CreatedBy).
			Return(request, nil)

		res := &RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Amount:    request.Amount,
			Title:     request.Title,
			Content:   request.Content,
			Tags: []*TagOverview{{
				ID:          tag.ID,
				Name:        tag.Name,
				Description: tag.Description,
				CreatedAt:   tag.CreatedAt,
				UpdatedAt:   tag.UpdatedAt,
			}},
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PostRequest(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
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
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		reqRequest := PutRequest{
			Amount:  random.Numeric(t, 1000000),
			Title:   request.Title,
			Content: random.AlphaNumeric(t, 50),
			Group:   &group.ID,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		var tags []*model.Tag

		updateRequest := &model.RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: request.CreatedBy,
			Amount:    reqRequest.Amount,
			Title:     request.Title,
			Content:   reqRequest.Content,
			Group:     group,
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

		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), group.ID).
			Return(group, nil)
		h.Repository.MockRequestRepository.
			EXPECT().
			UpdateRequest(c.Request().Context(), request.ID, reqRequest.Amount, request.Title, reqRequest.Content, tags, group).
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
			Amount:    updateRequest.Amount,
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

	t.Run("SuccessWithTag", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		date := time.Now()

		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		tag1 := &model.Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 30),
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		tag2 := &model.Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 30),
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		tags := []*model.Tag{tag1, tag2}
		reqRequest := PutRequest{
			Amount:  random.Numeric(t, 1000000),
			Title:   request.Title,
			Content: random.AlphaNumeric(t, 50),
			Tags:    []*uuid.UUID{&tag1.ID, &tag2.ID},
			Group:   &group.ID,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)

		updateRequest := &model.RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: request.CreatedBy,
			Amount:    reqRequest.Amount,
			Title:     request.Title,
			Content:   reqRequest.Content,
			Tags:      tags,
			Group:     group,
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

		h.Repository.MockTagRepository.
			EXPECT().
			GetTag(c.Request().Context(), tag1.ID).
			Return(tag1, nil)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTag(c.Request().Context(), tag2.ID).
			Return(tag2, nil)
		h.Repository.MockGroupRepository.
			EXPECT().
			GetGroup(c.Request().Context(), group.ID).
			Return(group, nil)
		h.Repository.MockRequestRepository.
			EXPECT().
			UpdateRequest(c.Request().Context(), request.ID, reqRequest.Amount, request.Title, reqRequest.Content, tags, group).
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
			Amount:    updateRequest.Amount,
			Title:     updateRequest.Title,
			Content:   updateRequest.Content,
			Tags: []*TagOverview{
				{
					ID:          tag1.ID,
					Name:        tag1.Name,
					Description: tag1.Description,
					CreatedAt:   tag1.CreatedAt,
					UpdatedAt:   tag1.UpdatedAt,
				},
				{
					ID:          tag2.ID,
					Name:        tag2.Name,
					Description: tag2.Description,
					CreatedAt:   tag2.CreatedAt,
					UpdatedAt:   tag2.UpdatedAt,
				},
			},
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
}

// 	func TestHandlers_PostComment(t *testing.T) {
// 		t.Parallel()

// 		t.Run("Success")
// }
