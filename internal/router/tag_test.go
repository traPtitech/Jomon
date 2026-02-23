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
	"github.com/traPtitech/Jomon/internal/model"
	"github.com/traPtitech/Jomon/internal/service"
	"github.com/traPtitech/Jomon/internal/testutil"
	"github.com/traPtitech/Jomon/internal/testutil/random"
	"go.uber.org/mock/gomock"
)

func TestHandlers_GetTags(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
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

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/tags", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTags(c.Request().Context()).
			Return(tags, nil)

		require.NoError(t, h.Handlers.GetTags(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*TagResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := lo.Map(tags, func(tag *model.Tag, _ int) *TagResponse {
			return &TagResponse{
				ID:        tag.ID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt,
				UpdatedAt: tag.UpdatedAt,
			}
		})
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		tags := []*model.Tag{}

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/tags", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTags(c.Request().Context()).
			Return(tags, nil)

		require.NoError(t, h.Handlers.GetTags(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*TagResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		var exp []*TagResponse
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("FailedToGetTags", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/tags", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		resErr := service.NewUnexpectedError(errors.New("failed to get tags"))
		h.Repository.MockTagRepository.
			EXPECT().
			GetTags(c.Request().Context()).
			Return(nil, resErr)

		err = h.Handlers.GetTags(c)
		require.Error(t, err)
		require.Equal(t, http.StatusInternalServerError, HTTPErrorHandlerInner(err).Code)
	})
}

func TestHandlers_PostTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
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
		reqTag := PostTagRequest{
			Name: tag.Name,
		}
		reqBody, err := json.Marshal(reqTag)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/tags", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		h.Repository.MockTagRepository.
			EXPECT().
			CreateTag(c.Request().Context(), tag.Name).
			Return(tag, nil)

		require.NoError(t, h.Handlers.PostTag(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got TagResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &TagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}
		testutil.RequireEqual(t, exp, &got, opts...)
	})

	t.Run("MissingName", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      "",
			CreatedAt: date,
			UpdatedAt: date,
		}
		reqTag := PostTagRequest{
			Name: "",
		}
		reqBody, err := json.Marshal(reqTag)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/tags", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		resErr := service.NewBadInputError("Tag name can't be empty.")
		h.Repository.MockTagRepository.
			EXPECT().
			CreateTag(c.Request().Context(), tag.Name).
			Return(nil, resErr)

		err = h.Handlers.PostTag(c)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, HTTPErrorHandlerInner(err).Code)
	})
}

func TestHandlers_PutTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
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
		reqTag := PutTagRequest{
			Name: tag.Name,
		}
		reqBody, err := json.Marshal(reqTag)
		require.NoError(t, err)
		updateTag := &model.Tag{
			ID:        tag.ID,
			Name:      reqTag.Name,
			CreatedAt: date,
			UpdatedAt: time.Now(),
		}

		e := echo.New()
		path := fmt.Sprintf("/api/tags/%s", tag.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues(tag.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		h.Repository.MockTagRepository.
			EXPECT().
			UpdateTag(c.Request().Context(), tag.ID, reqTag.Name).
			Return(updateTag, nil)

		require.NoError(t, h.Handlers.PutTag(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got TagResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &TagResponse{
			ID:        updateTag.ID,
			Name:      updateTag.Name,
			CreatedAt: updateTag.CreatedAt,
			UpdatedAt: updateTag.UpdatedAt,
		}
		testutil.RequireEqual(t, exp, &got, opts...)
	})

	t.Run("MissingName", func(t *testing.T) {
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
		reqTag := PutTagRequest{
			Name: "",
		}
		reqBody, err := json.Marshal(reqTag)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/tags/%s", tag.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues(tag.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		resErr := service.NewBadInputError("Tag name can't be empty.")
		h.Repository.MockTagRepository.
			EXPECT().
			UpdateTag(c.Request().Context(), tag.ID, reqTag.Name).
			Return(nil, resErr)

		err = h.Handlers.PutTag(c)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, HTTPErrorHandlerInner(err).Code)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		invalidUUID := "invalid-uuid"
		date := time.Now()
		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		reqTag := PutTagRequest{
			Name: tag.Name,
		}
		reqBody, err := json.Marshal(reqTag)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/tags/%s", invalidUUID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues(invalidUUID)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.PutTag(c)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, HTTPErrorHandlerInner(err).Code)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		tag := &model.Tag{
			ID:        uuid.Nil,
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		reqTag := PutTagRequest{
			Name: tag.Name,
		}
		reqBody, err := json.Marshal(reqTag)
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/tags/%s", tag.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues(tag.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.PutTag(c)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, HTTPErrorHandlerInner(err).Code)
	})
}

func TestHandlers_DeleteTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
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

		e := echo.New()
		path := fmt.Sprintf("/api/tags/%s", tag.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues(tag.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		h.Repository.MockTagRepository.
			EXPECT().
			DeleteTag(c.Request().Context(), tag.ID).
			Return(nil)

		require.NoError(t, h.Handlers.DeleteTag(c))
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("UnknownID", func(t *testing.T) {
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

		e := echo.New()
		path := fmt.Sprintf("/api/tags/%s", tag.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues(tag.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		resErr := service.NewUnexpectedError(errors.New("Unknown Tag ID"))
		h.Repository.MockTagRepository.
			EXPECT().
			DeleteTag(c.Request().Context(), tag.ID).
			Return(resErr)

		err = h.Handlers.DeleteTag(c)
		require.Error(t, err)
		require.Equal(t, http.StatusInternalServerError, HTTPErrorHandlerInner(err).Code)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		invalidUUID := "invalid-uuid"
		e := echo.New()
		path := fmt.Sprintf("/api/tags/%s", invalidUUID)
		req := httptest.NewRequestWithContext(
			ctx, http.MethodDelete, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues(invalidUUID)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.DeleteTag(c)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, HTTPErrorHandlerInner(err).Code)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		tag := &model.Tag{
			ID:        uuid.Nil,
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}

		e := echo.New()
		path := fmt.Sprintf("/api/tags/%s", tag.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues(tag.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.DeleteTag(c)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, HTTPErrorHandlerInner(err).Code)
	})
}
