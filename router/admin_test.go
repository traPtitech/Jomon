package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetAdmins(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		admin := uuid.New()

		admins := []*uuid.UUID{
			&admin,
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

		res := admins
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetAdmins(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})
}
