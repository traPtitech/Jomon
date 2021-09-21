package router

import (
	"bytes"
	"encoding/json"
	"errors"
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

// TODO: 直す
func TestHandlers_GetTags(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		date := time.Now()

		tag1 := &model.Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		tag2 := &model.Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
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

		resOverview := []*TagOverview{}
		for _, tag := range tags {
			resOverview = append(resOverview, &TagOverview{
				ID:          tag.ID,
				Name:        tag.Name,
				Description: tag.Description,
				CreatedAt:   tag.CreatedAt,
				UpdatedAt:   tag.UpdatedAt,
			})
		}
		res := TagResponse{resOverview}
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
		res := TagResponse{resOverview}
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

// TODO: 直す
func TestHandlers_PostTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
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

		reqTag := Tag{
			Name:        tag.Name,
			Description: tag.Description,
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
			CreateTag(c.Request().Context(), tag.Name, tag.Description).
			Return(tag, nil)

		res := TagOverview{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
			CreatedAt:   tag.CreatedAt,
			UpdatedAt:   tag.UpdatedAt,
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
			ID:          uuid.New(),
			Name:        "",
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		reqTag := Tag{
			Name:        "",
			Description: tag.Description,
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
			CreateTag(c.Request().Context(), tag.Name, tag.Description).
			Return(nil, mocErr)

		err = h.Handlers.PostTag(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})
}

/*

// TODO: 直す
func TestHandlers_PutTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)
		date := time.Now()

		tag := &model.Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			UpdateTag(ctx, tag.ID, tag.Name, tag.Description).
			Return(tag, nil)

		req := Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}

		var resBody TagOverview
		path := fmt.Sprintf("/api/tags/%s", tag.ID.String())
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.PUT, path, &req, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, tag.ID, resBody.ID)
		assert.Equal(t, tag.Name, resBody.Name)
		assert.Equal(t, tag.Description, resBody.Description)
	})

	t.Run("MissingName", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)
		date := time.Now()

		tag := &model.Tag{
			ID:          uuid.New(),
			Name:        "",
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			UpdateTag(ctx, tag.ID, tag.Name, tag.Description).
			Return(nil, errors.New("Tag name can't be empty."))

		req := Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}

		path := fmt.Sprintf("/api/tags/%s", tag.ID.String())
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.PUT, path, &req, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)
		date := time.Now()

		tag := &model.Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		req := Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}

		path := "/api/tags/hoge" // Invalid UUID
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.PUT, path, &req, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)
		date := time.Now()

		tag := &model.Tag{
			ID:          uuid.Nil,
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		req := Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}

		path := fmt.Sprintf("/api/tags/%s", tag.ID.String())
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.PUT, path, &req, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})
}


// TODO: 直す
func TestHandlers_DeleteTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)

		id := uuid.New()

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			DeleteTag(ctx, id).
			Return(nil)

		path := fmt.Sprintf("/api/tags/%s", id.String())
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.DELETE, path, nil, nil)
		assert.Equal(t, http.StatusOK, statusCode)
	})

	t.Run("UnknownID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)

		id := uuid.New()

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			DeleteTag(ctx, id).
			Return(errors.New("Tag not found"))

		path := fmt.Sprintf("/api/tags/%s", id.String())
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.DELETE, path, nil, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)

		path := "/api/tags/hoge"
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.DELETE, path, nil, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)

		path := fmt.Sprintf("/api/tags/%s", uuid.Nil)
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.DELETE, path, nil, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})
}
*/
