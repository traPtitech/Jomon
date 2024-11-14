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

	"github.com/gorilla/sessions"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestHandlers_GetUsers(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
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
			CreatedAt:   user1.CreatedAt,
			UpdatedAt:   user1.UpdatedAt,
		}
		resUser2 := &User{
			ID:          user2.ID,
			Name:        user2.Name,
			DisplayName: user2.DisplayName,
			Admin:       user2.Admin,
			CreatedAt:   user2.CreatedAt,
			UpdatedAt:   user2.UpdatedAt,
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

	t.Run("Success2", func(t *testing.T) {
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

	t.Run("FailedToGetUsers", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/users")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		mocErr := errors.New("failed to get users")
		h.Repository.MockUserRepository.
			EXPECT().
			GetUsers(c.Request().Context()).
			Return(nil, mocErr)

		err = h.Handlers.GetUsers(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})
}

func TestHandlers_UpdateUserInfo(t *testing.T) {
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
			UpdateUser(
				c.Request().Context(),
				user.ID, updateUser.Name, updateUser.DisplayName, updateUser.Admin).
			Return(updateUser, nil)

		if assert.NoError(t, h.Handlers.UpdateUserInfo(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(bodyResUser), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})
	t.Run("FailedToUpdateUser", func(t *testing.T) {
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
		require.NoError(t, err)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByName(c.Request().Context(), updateUser.Name).
			Return(user, nil)
		mocErr := errors.New("failed to get users.")
		h.Repository.MockUserRepository.
			EXPECT().
			UpdateUser(
				c.Request().Context(),
				user.ID, updateUser.Name, updateUser.DisplayName, updateUser.Admin).
			Return(nil, mocErr)

		err = h.Handlers.UpdateUserInfo(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
		}
	})

	t.Run("FailedToGetUser", func(t *testing.T) {
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
		require.NoError(t, err)
		mocErr := errors.New("user not found.")
		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByName(c.Request().Context(), updateUser.Name).
			Return(nil, mocErr)

		err = h.Handlers.UpdateUserInfo(c)
		if assert.Error(t, err) {
			assert.Equal(t, echo.NewHTTPError(http.StatusBadRequest, mocErr), err)
		}
	})
}

func TestHandlers_GetMe(t *testing.T) {
	t.Parallel()
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, random.Numeric(t, 2) == 1)
		user := User{
			ID:          accessUser.ID,
			Name:        accessUser.Name,
			DisplayName: accessUser.DisplayName,
			Admin:       accessUser.Admin,
			CreatedAt:   accessUser.CreatedAt,
			UpdatedAt:   accessUser.UpdatedAt,
		}
		bodyAccessUser, err := json.Marshal(user)
		assert.NoError(t, err)

		e := echo.New()
		req, err := http.NewRequest(http.MethodPut, "/api/users/me", nil)
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		sess, err := session.Get(h.Handlers.SessionName, c)
		require.NoError(t, err)
		sess.Values[sessionUserKey] = user
		require.NoError(t, sess.Save(c.Request(), c.Response()))

		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByID(c.Request().Context(), user.ID).
			Return(accessUser, nil)

		err = h.Handlers.GetMe(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(bodyAccessUser), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("FailedToGetUser", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		e := echo.New()
		e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
		req, err := http.NewRequest(http.MethodPut, "/api/users/me", nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		mw := session.Middleware(sessions.NewCookieStore([]byte("secret")))
		hn := mw(echo.HandlerFunc(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		}))
		err = hn(c)
		require.NoError(t, err)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		err = h.Handlers.GetMe(c)
		if assert.Error(t, err) {
			assert.Equal(
				t,
				echo.NewHTTPError(http.StatusInternalServerError, "failed to get user info"),
				err)
		}
	})
}
