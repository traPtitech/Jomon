package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/testutil"
	"github.com/traPtitech/Jomon/testutil/random"
	"go.uber.org/mock/gomock"
)

// TODO: これ消す userFromModelUserがある
func modelUserToUser(user *model.User) *User {
	return &User{
		ID:             user.ID,
		Name:           user.Name,
		DisplayName:    user.DisplayName,
		AccountManager: user.AccountManager,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		DeletedAt:      user.DeletedAt,
	}
}

func TestHandlers_GetUsers(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		user1 := makeUser(t, random.Numeric(t, 2) == 1)
		user2 := makeUser(t, random.Numeric(t, 2) == 1)
		users := []*model.User{user1, user2}

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUsers(c.Request().Context()).
			Return(users, nil)

		require.NoError(t, h.Handlers.GetUsers(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*User
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := lo.Map(users, func(u *model.User, _ int) *User {
			return modelUserToUser(u)
		})
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		users := []*model.User{}

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockUserRepository.
			EXPECT().
			GetUsers(c.Request().Context()).
			Return(users, nil)

		require.NoError(t, h.Handlers.GetUsers(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*User
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := lo.Map(users, func(u *model.User, _ int) *User {
			return modelUserToUser(u)
		})
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("FailedToGetUsers", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
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
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; resErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
	})
}

func TestHandlers_UpdateUserInfo(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		user := makeUser(t, random.Numeric(t, 2) == 1)
		updateUser := &model.User{
			ID:             user.ID,
			Name:           user.Name,
			DisplayName:    user.DisplayName,
			AccountManager: !user.AccountManager,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      time.Now(),
		}

		reqUser := PutUserRequest{
			Name:           updateUser.Name,
			DisplayName:    updateUser.DisplayName,
			AccountManager: updateUser.AccountManager,
		}
		reqBody, err := json.Marshal(reqUser)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPut, "/api/users", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		h.Repository.MockUserRepository.
			EXPECT().
			GetUserByName(c.Request().Context(), updateUser.Name).
			Return(user, nil)
		h.Repository.MockUserRepository.
			EXPECT().
			UpdateUser(
				c.Request().Context(),
				user.ID, updateUser.Name, updateUser.DisplayName, updateUser.AccountManager).
			Return(updateUser, nil)

		require.NoError(t, h.Handlers.UpdateUserInfo(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got User
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		// FIXME: #835
		opts = append(opts,
			cmpopts.IgnoreFields(User{}, "CreatedAt", "UpdatedAt", "DeletedAt"))
		exp := modelUserToUser(updateUser)
		testutil.RequireEqual(t, exp, &got, opts...)
	})

	t.Run("FailedToUpdateUser", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		user := makeUser(t, random.Numeric(t, 2) == 1)
		updateUser := &model.User{
			ID:             user.ID,
			Name:           user.Name,
			DisplayName:    user.DisplayName,
			AccountManager: !user.AccountManager,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      time.Now(),
		}
		reqUser := PutUserRequest{
			Name:           updateUser.Name,
			DisplayName:    updateUser.DisplayName,
			AccountManager: updateUser.AccountManager,
		}
		bodyReqUser, err := json.Marshal(reqUser)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPut, "/api/users", bytes.NewReader(bodyReqUser))
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
				user.ID, updateUser.Name, updateUser.DisplayName, updateUser.AccountManager).
			Return(nil, mocErr)

		err = h.Handlers.UpdateUserInfo(c)
		require.Error(t, err)
		// FIXME: http.StatusInternalServerErrorだけ判定したい; mocErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusInternalServerError, mocErr), err)
	})

	t.Run("FailedToGetUser", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		user := makeUser(t, random.Numeric(t, 2) == 1)
		updateUser := &model.User{
			ID:             user.ID,
			Name:           user.Name,
			DisplayName:    user.DisplayName,
			AccountManager: !user.AccountManager,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      time.Now(),
		}
		reqUser := PutUserRequest{
			Name:           updateUser.Name,
			DisplayName:    updateUser.DisplayName,
			AccountManager: updateUser.AccountManager,
		}
		bodyReqUser, err := json.Marshal(reqUser)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPut, "/api/users", bytes.NewReader(bodyReqUser))
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
		require.Error(t, err)
		// FIXME: http.StatusBadRequestだけ判定したい; mocErrの内容は関係ない
		require.Equal(t, echo.NewHTTPError(http.StatusBadRequest, mocErr), err)
	})
}

func TestHandlers_GetMe(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		accessUser := makeUser(t, random.Numeric(t, 2) == 1)
		user := userFromModelUser(*accessUser)

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, "/api/users/me", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/users/me")
		c.Set(loginUserKey, user)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		require.NoError(t, h.Handlers.GetMe(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got User
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelUserToUser(accessUser)
		testutil.RequireEqual(t, exp, &got, opts...)
	})

	// TODO: checkLoginMiddlewareのテストを追加する
}
