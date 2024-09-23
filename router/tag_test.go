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

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/testutil/random"
	"go.uber.org/mock/gomock"
)

func TestHandlers_GetTags(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
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
		req, err := http.NewRequest(http.MethodGet, "/api/tags", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTags(c.Request().Context()).
			Return(tags, nil)

		resOverview := lo.Map(tags, func(tag *model.Tag, index int) *TagOverview {
			return &TagOverview{
				ID:        tag.ID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt,
				UpdatedAt: tag.UpdatedAt,
			}
		})

		res := resOverview
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetTags(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		tags := []*model.Tag{}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/tags", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockTagRepository.
			EXPECT().
			GetTags(c.Request().Context()).
			Return(tags, nil)

		resOverview := []*TagOverview{}
		res := resOverview
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetTags(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("FailedToGetTags", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/tags", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		mocErr := errors.New("failed to get tags")
		h.Repository.MockTagRepository.
			EXPECT().
			GetTags(c.Request().Context()).
			Return(nil, mocErr)

		err = h.Handlers.GetTags(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})
}

func TestHandlers_PostTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}

		reqTag := Tag{
			Name: tag.Name,
		}
		reqBody, err := json.Marshal(reqTag)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, "/api/tags", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		h.Repository.MockTagRepository.
			EXPECT().
			CreateTag(c.Request().Context(), tag.Name).
			Return(tag, nil)

		res := TagOverview{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PostTag(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("MissingName", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      "",
			CreatedAt: date,
			UpdatedAt: date,
		}

		reqTag := Tag{
			Name: "",
		}
		reqBody, err := json.Marshal(reqTag)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPost, "/api/tags", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		mocErr := errors.New("Tag name can't be empty.")
		h.Repository.MockTagRepository.
			EXPECT().
			CreateTag(c.Request().Context(), tag.Name).
			Return(nil, mocErr)

		err = h.Handlers.PostTag(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})
}

func TestHandlers_PutTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}

		reqTag := Tag{
			Name:        tag.Name,
			Description: random.AlphaNumeric(t, 50),
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
		req, err := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("/api/tags/%s", tag.ID),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues(tag.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		h.Repository.MockTagRepository.
			EXPECT().
			UpdateTag(c.Request().Context(), tag.ID, reqTag.Name).
			Return(updateTag, nil)

		res := TagOverview{
			ID:        updateTag.ID,
			Name:      updateTag.Name,
			CreatedAt: updateTag.CreatedAt,
			UpdatedAt: updateTag.UpdatedAt,
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PutTag(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("MissingName", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}

		reqTag := Tag{
			Name: "",
		}
		reqBody, err := json.Marshal(reqTag)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("/api/tags/%s", tag.ID),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues(tag.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		mocErr := errors.New("Tag name can't be empty.")
		h.Repository.MockTagRepository.
			EXPECT().
			UpdateTag(c.Request().Context(), tag.ID, reqTag.Name).
			Return(nil, mocErr)

		err = h.Handlers.PutTag(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}

		reqTag := Tag{
			Name:        tag.Name,
			Description: random.AlphaNumeric(t, 50),
		}
		reqBody, err := json.Marshal(reqTag)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, "/api/tags/hoge", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues("hoge")

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		_, resErr := uuid.Parse("hoge")

		err = h.Handlers.PutTag(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		tag := &model.Tag{
			ID:        uuid.Nil,
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}

		reqTag := Tag{
			Name:        tag.Name,
			Description: random.AlphaNumeric(t, 50),
		}
		reqBody, err := json.Marshal(reqTag)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("/api/tags/%s", tag.ID),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues(tag.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		err = h.Handlers.PutTag(c)
		if assert.Error(t, err) {
			assert.Equal(
				t,
				echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid tag ID")),
				err)
		}
	})
}

func TestHandlers_DeleteTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}

		reqTag := Tag{
			Name:        tag.Name,
			Description: random.AlphaNumeric(t, 50),
		}
		reqBody, err := json.Marshal(reqTag)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/tags/%s", tag.ID),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues(tag.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		h.Repository.MockTagRepository.
			EXPECT().
			DeleteTag(c.Request().Context(), tag.ID).
			Return(nil)

		if assert.NoError(t, h.Handlers.DeleteTag(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("UnknownID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}

		reqTag := Tag{
			Name:        tag.Name,
			Description: random.AlphaNumeric(t, 50),
		}
		reqBody, err := json.Marshal(reqTag)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/tags/%s", tag.ID),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues(tag.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		mocErr := errors.New("Unknown Tag ID")
		h.Repository.MockTagRepository.
			EXPECT().
			DeleteTag(c.Request().Context(), tag.ID).
			Return(mocErr)

		err = h.Handlers.DeleteTag(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}

		reqTag := Tag{
			Name:        tag.Name,
			Description: random.AlphaNumeric(t, 50),
		}
		reqBody, err := json.Marshal(reqTag)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodDelete, "/api/tags/hoge", bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues("hoge")

		_, resErr := uuid.Parse("hoge")

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		err = h.Handlers.DeleteTag(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		date := time.Now()

		tag := &model.Tag{
			ID:        uuid.Nil,
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}

		reqTag := Tag{
			Name:        tag.Name,
			Description: random.AlphaNumeric(t, 50),
		}
		reqBody, err := json.Marshal(reqTag)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/api/tags/%s", tag.ID),
			bytes.NewReader(reqBody))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/tags/:tagID")
		c.SetParamNames("tagID")
		c.SetParamValues(tag.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		err = h.Handlers.DeleteTag(c)
		if assert.Error(t, err) {
			assert.Equal(
				t,
				echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid tag ID")),
				err)
		}
	})
}
