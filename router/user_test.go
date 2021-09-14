package router

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestGetUsers(t *testing.T) {
	t.Parallel()

	t.Run("Success 1", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		user1 := makeUser(t, random.Numeric(t, 2) == 1)
		user2 := makeUser(t, random.Numeric(t, 2) == 1)
		users := []*model.User{user1, user2}

		resUser1 := &User{
			ID:          user1.ID,
			Name:        user1.Name,
			DisplayName: user1.DisplayName,
			Admin:       user1.Admin,
		}
		resUser2 := &User{
			ID:          user2.ID,
			Name:        user2.Name,
			DisplayName: user2.DisplayName,
			Admin:       user2.Admin,
		}
		resUsers := []*User{resUser1, resUser2}
		resBody, err := json.Marshal(resUsers)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/api/users", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUsers(c.Request().Context()).
			Return(users, nil)

		if assert.NoError(t, h.Handlers.GetUsers(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("Success 2", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		resUsers := []*User{}
		body, err := json.Marshal(resUsers)
		assert.NoError(t, err)

		users := []*model.User{}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUsers(c.Request().Context()).
			Return(users, nil)

		if assert.NoError(t, h.Handlers.GetUsers(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(body), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("Fail", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUsers(c.Request().Context()).
			Return(nil, errors.New("failed to get users"))

		if assert.Error(t, h.Handlers.GetUsers(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

// TODO: 直す
func TestPutUser(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		user := makeUser(t, random.Numeric(t, 2) == 1)

		updateUser := &model.User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: user.DisplayName,
			Admin:       !user.Admin,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   time.Now(),
		}

		resUser := User{
			ID:          updateUser.ID,
			Name:        updateUser.Name,
			DisplayName: updateUser.DisplayName,
			Admin:       updateUser.Admin,
		}
		bodyResUser, err := json.Marshal(resUser)
		assert.NoError(t, err)

		reqUser := PutUserRequest{
			Name:        updateUser.Name,
			DisplayName: updateUser.DisplayName,
			Admin:       updateUser.Admin,
		}
		bodyReqUser, err := json.Marshal(reqUser)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, "/api/users", bytes.NewReader(bodyReqUser))
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)

		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByName(c.Request().Context(), updateUser.Name).
			Return(user, nil)
		h.Repository.MockUserRepository.
			EXPECT().
			UpdateUser(c.Request().Context(), user.ID, updateUser.Name, updateUser.DisplayName, updateUser.Admin).
			Return(updateUser, nil)

		if assert.NoError(t, h.Handlers.UpdateUserInfo(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(bodyResUser), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("FailedToUpdateUser", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		user := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)
		date := time.Now()

		updateUser := &model.User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       true,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   date,
		}

		req := &PutUserRequest{
			Name:        updateUser.Name,
			DisplayName: updateUser.DisplayName,
			Admin:       updateUser.Admin,
		}

		ctx := context.Background()
		th.Repository.MockUserRepository.
			EXPECT().
			GetUserByName(ctx, updateUser.Name).
			Return(user, nil)
		th.Repository.MockUserRepository.
			EXPECT().
			UpdateUser(ctx, user.ID, updateUser.Name, updateUser.DisplayName, updateUser.Admin).
			Return(nil, errors.New("failed to get users."))

		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.PUT, "/api/users", &req, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})

	t.Run("FailedToGetUser", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		user := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)
		date := time.Now()

		updateUser := &model.User{
			ID:          user.ID,
			Name:        user.Name,
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       true,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   date,
		}

		req := &PutUserRequest{
			Name:        updateUser.Name,
			DisplayName: updateUser.DisplayName,
			Admin:       updateUser.Admin,
		}

		ctx := context.Background()
		th.Repository.MockUserRepository.
			EXPECT().
			GetUserByName(ctx, updateUser.Name).
			Return(nil, errors.New("user not found."))

		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.PUT, "/api/users", &req, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}

// TODO: 直す
func TestHandlers_GetMe(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)

		ctx := context.Background()
		th.Repository.MockUserRepository.
			EXPECT().
			GetUserByName(ctx, accessUser.Name).
			Return(accessUser, nil)

		var resBody User
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.GET, "/api/users/me", nil, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, accessUser.Name, resBody.Name)
		assert.Equal(t, accessUser.DisplayName, resBody.DisplayName)
		assert.Equal(t, accessUser.Admin, resBody.Admin)
	})

	t.Run("FailedToGetUser", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)

		ctx := context.Background()
		th.Repository.MockUserRepository.
			EXPECT().
			GetUserByName(ctx, accessUser.Name).
			Return(nil, errors.New("failed to get user."))

		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.GET, "/api/users/me", nil, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
