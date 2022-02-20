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
	"github.com/traPtitech/Jomon/ent"
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
			CreatedAt: date1,
			UpdatedAt: date1,
		}
		request2 := &model.RequestResponse{
			ID:        uuid.New(),
			Status:    model.Submitted.String(),
			CreatedBy: uuid.New(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
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
		require.NoError(t, err)
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
			},
			{
				ID:        request1.ID,
				Status:    request1.Status,
				CreatedAt: request1.CreatedAt,
				UpdatedAt: request1.UpdatedAt,
				CreatedBy: request1.CreatedBy,
				Amount:    request1.Amount,
				Title:     request1.Title,
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
		resErr := errors.New("Failed to get requests.")
		h.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(c.Request().Context(), query).
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

		userID := uuid.New()
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      userID,
			Comment:   random.AlphaNumeric(t, 30),
			CreatedAt: date,
			UpdatedAt: date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
			Comments:  []*model.Comment{comment},
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: userID,
		}

		reqRequest := Request{
			CreatedBy: request.CreatedBy,
			Amount:    request.Amount,
			Title:     request.Title,
			Comment:   comment.Comment,
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
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			CreateRequest(c.Request().Context(), reqRequest.Amount, reqRequest.Title, tags, group, reqRequest.CreatedBy).
			Return(request, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(c.Request().Context(), reqRequest.Comment, request.ID, reqRequest.CreatedBy).
			Return(comment, nil)

		resComment := &CommentDetail{
			ID:        comment.ID,
			User:      comment.User,
			Comment:   comment.Comment,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}
		res := &RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Amount:    request.Amount,
			Title:     request.Title,
			Comments:  []*CommentDetail{resComment},
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

		userID := uuid.New()
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      userID,
			Comment:   random.AlphaNumeric(t, 30),
			CreatedAt: date,
			UpdatedAt: date,
		}
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
			Comments:  []*model.Comment{comment},
			Tags:      tags,
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: userID,
		}

		reqRequest := Request{
			CreatedBy: request.CreatedBy,
			Amount:    request.Amount,
			Title:     request.Title,
			Comment:   comment.Comment,
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
		require.NoError(t, err)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTag(c.Request().Context(), tag.ID).
			Return(tag, nil)
		h.Repository.MockRequestRepository.
			EXPECT().
			CreateRequest(c.Request().Context(), reqRequest.Amount, reqRequest.Title, tags, group, reqRequest.CreatedBy).
			Return(request, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(c.Request().Context(), reqRequest.Comment, request.ID, reqRequest.CreatedBy).
			Return(comment, nil)

		resComment := &CommentDetail{
			ID:        comment.ID,
			User:      comment.User,
			Comment:   comment.Comment,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}
		res := &RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Amount:    request.Amount,
			Title:     request.Title,
			Comments:  []*CommentDetail{resComment},
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

	t.Run("SuccessWithGroup", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		userID := uuid.New()
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      userID,
			Comment:   random.AlphaNumeric(t, 30),
			CreatedAt: date,
			UpdatedAt: date,
		}
		budget := random.Numeric(t, 100000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
			Comments:  []*model.Comment{comment},
			Group:     group,
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: userID,
		}

		reqRequest := Request{
			CreatedBy: request.CreatedBy,
			Amount:    request.Amount,
			Title:     request.Title,
			Group:     &group.ID,
			Comment:   comment.Comment,
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		var tags []*model.Tag

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
			CreateRequest(c.Request().Context(), reqRequest.Amount, reqRequest.Title, tags, group, reqRequest.CreatedBy).
			Return(request, nil)
		h.Repository.MockCommentRepository.
			EXPECT().
			CreateComment(c.Request().Context(), reqRequest.Comment, request.ID, reqRequest.CreatedBy).
			Return(comment, nil)

		resComment := &CommentDetail{
			ID:        comment.ID,
			User:      comment.User,
			Comment:   comment.Comment,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}
		res := &RequestResponse{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: request.UpdatedAt,
			CreatedBy: request.CreatedBy,
			Amount:    request.Amount,
			Title:     request.Title,
			Comments:  []*CommentDetail{resComment},
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

	t.Run("UnknownTagID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		userID := uuid.New()
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      userID,
			Comment:   random.AlphaNumeric(t, 30),
			CreatedAt: date,
			UpdatedAt: date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
			Comments:  []*model.Comment{comment},
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: userID,
		}

		unknownTagID := uuid.New()

		reqRequest := Request{
			CreatedBy: request.CreatedBy,
			Amount:    request.Amount,
			Title:     request.Title,
			Tags:      []*uuid.UUID{&unknownTagID},
			Comment:   comment.Comment,
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

		userID := uuid.New()
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      userID,
			Comment:   random.AlphaNumeric(t, 30),
			CreatedAt: date,
			UpdatedAt: date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
			Comments:  []*model.Comment{comment},
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: userID,
		}

		unknownGroupID := uuid.New()

		reqRequest := Request{
			CreatedBy: request.CreatedBy,
			Amount:    request.Amount,
			Title:     request.Title,
			Comment:   comment.Comment,
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

		userID := uuid.New()
		comment := &model.Comment{
			ID:        uuid.New(),
			User:      userID,
			Comment:   random.AlphaNumeric(t, 30),
			CreatedAt: date,
			UpdatedAt: date,
		}
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
			Comments:  []*model.Comment{comment},
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: userID,
		}

		reqRequest := Request{
			CreatedBy: request.CreatedBy,
			Amount:    request.Amount,
			Title:     request.Title,
			Comment:   comment.Comment,
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

		var resErr *ent.NotFoundError
		errors.As(errors.New("unknown id"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockRequestRepository.
			EXPECT().
			CreateRequest(c.Request().Context(), reqRequest.Amount, reqRequest.Title, tags, group, reqRequest.CreatedBy).
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
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
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
			Amount:    request.Amount,
			Title:     request.Title,
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
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
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
			Amount:    request.Amount,
			Title:     request.Title,
			Comments:  resComments,
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
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 100000),
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}

		reqRequest := PutRequest{
			Amount: random.Numeric(t, 100000),
			Title:  random.AlphaNumeric(t, 30),
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		var tags []*model.Tag
		var group *model.Group

		updateRequest := &model.RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: request.CreatedBy,
			Amount:    reqRequest.Amount,
			Title:     reqRequest.Title,
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
			UpdateRequest(c.Request().Context(), request.ID, reqRequest.Amount, reqRequest.Title, tags, group).
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
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
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
			Amount: random.Numeric(t, 1000000),
			Title:  random.AlphaNumeric(t, 30),
			Tags:   []*uuid.UUID{&tag1.ID, &tag2.ID},
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		var group *model.Group

		updateRequest := &model.RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: request.CreatedBy,
			Amount:    reqRequest.Amount,
			Title:     request.Title,
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
			UpdateRequest(c.Request().Context(), request.ID, reqRequest.Amount, reqRequest.Title, tags, group).
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
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
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
			Amount: random.Numeric(t, 1000000),
			Title:  random.AlphaNumeric(t, 30),
			Group:  &group.ID,
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
			UpdateRequest(c.Request().Context(), request.ID, reqRequest.Amount, reqRequest.Title, tags, group).
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
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 100000),
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}

		reqRequest := PutRequest{
			Amount: random.Numeric(t, 100000),
			Title:  random.AlphaNumeric(t, 30),
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		var tags []*model.Tag
		var group *model.Group

		updateRequest := &model.RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: request.CreatedBy,
			Amount:    reqRequest.Amount,
			Title:     reqRequest.Title,
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
			UpdateRequest(c.Request().Context(), request.ID, reqRequest.Amount, reqRequest.Title, tags, group).
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
			Amount:    updateRequest.Amount,
			Title:     updateRequest.Title,
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
			Amount: random.Numeric(t, 100000),
			Title:  random.AlphaNumeric(t, 30),
		}
		reqBody, err := json.Marshal(reqRequest)
		require.NoError(t, err)
		var tags []*model.Tag
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
			UpdateRequest(c.Request().Context(), unknownID, reqRequest.Amount, reqRequest.Title, tags, group).
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
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}

		tag := &model.Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 30),
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		reqRequest := PutRequest{
			Amount: random.Numeric(t, 100000),
			Title:  random.AlphaNumeric(t, 30),
			Tags:   []*uuid.UUID{&tag.ID},
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
			Status:    model.Submitted.String(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
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
			Amount: random.Numeric(t, 100000),
			Title:  random.AlphaNumeric(t, 30),
			Group:  &group.ID,
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
