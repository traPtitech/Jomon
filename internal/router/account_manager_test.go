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
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/internal/model"
	"github.com/traPtitech/Jomon/internal/testutil"
	"go.uber.org/mock/gomock"
)

func TestHandler_GetAccountManagers(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accountManager := &model.AccountManager{
			ID: uuid.New(),
		}

		accountManagers := []*model.AccountManager{
			accountManager,
		}

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/account-managers", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockAccountManagerRepository.
			EXPECT().
			GetAccountManagers(c.Request().Context()).
			Return(accountManagers, nil)

		require.NoError(t, err)

		require.NoError(t, h.Handlers.GetAccountManagers(c))
		testutil.AssertEqual(t, http.StatusOK, rec.Code)
		var res []uuid.UUID
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)
		testutil.RequireEqual(t, []uuid.UUID{accountManager.ID}, res)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)
		var accountManagers []*model.AccountManager

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/account-managers", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockAccountManagerRepository.
			EXPECT().
			GetAccountManagers(c.Request().Context()).
			Return(accountManagers, nil)

		require.NoError(t, err)

		require.NoError(t, h.Handlers.GetAccountManagers(c))
		testutil.AssertEqual(t, http.StatusOK, rec.Code)
		var res []uuid.UUID
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Empty(t, res)
	})

	t.Run("FailedWithError", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/account-managers", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		resErr := errors.New("failed to get accountManagers")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockAccountManagerRepository.
			EXPECT().
			GetAccountManagers(c.Request().Context()).
			Return(nil, resErr)

		err = h.Handlers.GetAccountManagers(c)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})
}

func TestHandler_PostAccountManager(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accountManager := uuid.New()
		accountManagers := []uuid.UUID{accountManager}
		reqBody, err := json.Marshal(accountManagers)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/account-managers", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockAccountManagerRepository.
			EXPECT().
			AddAccountManagers(c.Request().Context(), accountManagers).
			Return(nil)

		require.NoError(t, h.Handlers.PostAccountManagers(c))
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("FailedWithError", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accountManager := uuid.New()
		accountManagers := []uuid.UUID{accountManager}
		reqBody, err := json.Marshal(accountManagers)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/account-managers", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		resErr := errors.New("failed to create accountManager")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockAccountManagerRepository.
			EXPECT().
			AddAccountManagers(c.Request().Context(), accountManagers).
			Return(resErr)

		err = h.Handlers.PostAccountManagers(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})

	t.Run("FailedWithEntConstraintError", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accountManager := uuid.New()
		accountManagers := []uuid.UUID{accountManager}
		reqBody, err := json.Marshal(accountManagers)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/account-managers", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var resErr *ent.ConstraintError
		errors.As(errors.New("failed to create accountManager"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockAccountManagerRepository.
			EXPECT().
			AddAccountManagers(c.Request().Context(), accountManagers).
			Return(resErr)

		err = h.Handlers.PostAccountManagers(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, resErr), err)
	})
}

func TestHandler_DeleteAccountManager(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accountManager := uuid.New()
		accountManagers := []uuid.UUID{accountManager}
		reqBody, err := json.Marshal(accountManagers)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodDelete, "/api/account-managers", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockAccountManagerRepository.
			EXPECT().
			DeleteAccountManagers(c.Request().Context(), accountManagers).
			Return(nil)

		require.NoError(t, h.Handlers.DeleteAccountManagers(c))
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("FailedWithError", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accountManager := uuid.New()
		accountManagers := []uuid.UUID{accountManager}
		reqBody, err := json.Marshal(accountManagers)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodDelete, "/api/account-managers", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		resErr := errors.New("failed to delete accountManager")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockAccountManagerRepository.
			EXPECT().
			DeleteAccountManagers(c.Request().Context(), accountManagers).
			Return(resErr)

		err = h.Handlers.DeleteAccountManagers(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})

	t.Run("InvalidAccountManagerID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		invalidUUID := "invalid-uuid"
		reqBody, err := json.Marshal([]string{invalidUUID})
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodDelete, "/api/account-managers", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.DeleteAccountManagers(c)

		require.Error(t, err)
		// FIXME: http.StatusBadRequestの判定をしたい
	})

	t.Run("FailedWithEntConstraintError", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accountManager := uuid.New()
		accountManagers := []uuid.UUID{accountManager}
		reqBody, err := json.Marshal(accountManagers)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodDelete, "/api/account-managers", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		var resErr *ent.ConstraintError
		errors.As(errors.New("failed to delete accountManager"), &resErr)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockAccountManagerRepository.
			EXPECT().
			DeleteAccountManagers(c.Request().Context(), accountManagers).
			Return(resErr)

		err = h.Handlers.DeleteAccountManagers(c)
		require.Error(t, err)
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, resErr), err)
	})
}
