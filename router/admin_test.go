package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/testutil"
	"go.uber.org/mock/gomock"
)

func TestHandler_GetAdmins(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		admin := &model.Admin{
			ID: uuid.New(),
		}

		admins := []*model.Admin{
			admin,
		}

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/admins", nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockAdminRepository.
			EXPECT().
			GetAdmins(c.Request().Context()).
			Return(admins, nil)

		require.NoError(t, err)

		assert.NoError(t, h.Handlers.GetAdmins(c))
		testutil.AssertEqual(t, http.StatusOK, rec.Code)
		var res []uuid.UUID
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)
		testutil.RequireEqual(t, []uuid.UUID{admin.ID}, res)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		var admins []*model.Admin

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/admins", nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockAdminRepository.
			EXPECT().
			GetAdmins(c.Request().Context()).
			Return(admins, nil)

		require.NoError(t, err)

		assert.NoError(t, h.Handlers.GetAdmins(c))
		testutil.AssertEqual(t, http.StatusOK, rec.Code)
		var res []uuid.UUID
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("FailedWithError", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/admins", nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		resErr := errors.New("failed to get admins")

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockAdminRepository.
			EXPECT().
			GetAdmins(c.Request().Context()).
			Return(nil, resErr)

		err = h.Handlers.GetAdmins(c)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})
}

func TestHandler_PostAdmin(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		admin := uuid.New()
		admins := []uuid.UUID{admin}
		reqBody, err := json.Marshal(admins)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			"/api/admins",
			bytes.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockAdminRepository.
			EXPECT().
			AddAdmins(c.Request().Context(), admins).
			Return(nil)

		assert.NoError(t, h.Handlers.PostAdmins(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("FailedWithError", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		admin := uuid.New()
		admins := []uuid.UUID{admin}
		reqBody, err := json.Marshal(admins)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			"/api/admins",
			bytes.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		resErr := errors.New("failed to create admin")

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockAdminRepository.
			EXPECT().
			AddAdmins(c.Request().Context(), admins).
			Return(resErr)

		err = h.Handlers.PostAdmins(c)
		assert.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})

	t.Run("FailedWithEntConstraintError", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		admin := uuid.New()
		admins := []uuid.UUID{admin}
		reqBody, err := json.Marshal(admins)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			"/api/admins",
			bytes.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var resErr *ent.ConstraintError
		errors.As(errors.New("failed to create admin"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockAdminRepository.
			EXPECT().
			AddAdmins(c.Request().Context(), admins).
			Return(resErr)

		err = h.Handlers.PostAdmins(c)
		assert.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})
}

func TestHandler_DeleteAdmin(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		admin := uuid.New()
		admins := []uuid.UUID{admin}
		reqBody, err := json.Marshal(admins)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			"/api/admins",
			bytes.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockAdminRepository.
			EXPECT().
			DeleteAdmins(c.Request().Context(), admins).
			Return(nil)

		assert.NoError(t, h.Handlers.DeleteAdmins(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("FailedWithError", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		admin := uuid.New()
		admins := []uuid.UUID{admin}
		reqBody, err := json.Marshal(admins)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			"/api/admins",
			bytes.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		resErr := errors.New("failed to delete admin")

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockAdminRepository.
			EXPECT().
			DeleteAdmins(c.Request().Context(), admins).
			Return(resErr)

		err = h.Handlers.DeleteAdmins(c)
		assert.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})

	t.Run("InvalidAdminID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		body := `["invalid"]`

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			"/api/admins",
			bytes.NewReader([]byte(body)))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		err = h.Handlers.DeleteAdmins(c)
		assert.Error(t, err)
		// FIXME: http.StatusBadRequestの判定をしたい
	})

	t.Run("FailedWithEntConstraintError", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		admin := uuid.New()
		admins := []uuid.UUID{admin}
		reqBody, err := json.Marshal(admins)
		require.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			"/api/admins",
			bytes.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var resErr *ent.ConstraintError
		errors.As(errors.New("failed to delete admin"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockAdminRepository.
			EXPECT().
			DeleteAdmins(c.Request().Context(), admins).
			Return(resErr)

		err = h.Handlers.DeleteAdmins(c)
		assert.Error(t, err)
		assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})
}
