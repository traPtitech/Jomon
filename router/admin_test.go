package router

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/model"
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

		res := []*uuid.UUID{
			&admin.ID,
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetAdmins(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		admins := []*model.Admin{}

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

		res := []*uuid.UUID{}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetAdmins(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
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
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})
}

func TestHandler_PostAdmin(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		admin := []uuid.UUID{uuid.New()}

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			"/api/admins",
			strings.NewReader(`["`+admin[0].String()+`"]`))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockAdminRepository.
			EXPECT().
			AddAdmins(c.Request().Context(), admin).
			Return(nil)

		if assert.NoError(t, h.Handlers.PostAdmins(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedWithError", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		adminID := uuid.New()

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			"/api/admins",
			strings.NewReader(`["`+adminID.String()+`"]`))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		resErr := errors.New("failed to create admin")

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockAdminRepository.
			EXPECT().
			AddAdmins(c.Request().Context(), []uuid.UUID{adminID}).
			Return(resErr)

		err = h.Handlers.PostAdmins(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})

	t.Run("FailedWithEntConstraintError", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		adminID := uuid.New()

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodPost,
			"/api/admins",
			strings.NewReader(`["`+adminID.String()+`"]`))
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
			AddAdmins(c.Request().Context(), []uuid.UUID{adminID}).
			Return(resErr)

		err = h.Handlers.PostAdmins(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
		}
	})
}

func TestHandler_DeleteAdmin(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		admin := []uuid.UUID{uuid.New()}

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			"/api/admins",
			strings.NewReader(`["`+admin[0].String()+`"]`))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockAdminRepository.
			EXPECT().
			DeleteAdmins(c.Request().Context(), admin).
			Return(nil)

		if assert.NoError(t, h.Handlers.DeleteAdmins(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("FailedWithError", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		adminID := uuid.New()

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			"/api/admins",
			strings.NewReader(`["`+adminID.String()+`"]`))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		resErr := errors.New("failed to delete admin")

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockAdminRepository.
			EXPECT().
			DeleteAdmins(c.Request().Context(), []uuid.UUID{adminID}).
			Return(resErr)

		err = h.Handlers.DeleteAdmins(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})

	t.Run("InvalidAdminID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			"/api/admins",
			strings.NewReader(`["invalid"]`))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		err = h.Handlers.DeleteAdmins(c)
		assert.Error(t, err)
	})

	t.Run("FailedWithEntConstraintError", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		adminID := uuid.New()

		e := echo.New()
		req, err := http.NewRequest(
			http.MethodDelete,
			"/api/admins",
			strings.NewReader(`["`+adminID.String()+`"]`))
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
			DeleteAdmins(c.Request().Context(), []uuid.UUID{adminID}).
			Return(resErr)

		err = h.Handlers.DeleteAdmins(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
		}
	})
}
