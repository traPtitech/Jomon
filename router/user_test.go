package router

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestHandlers_GetUsers(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
		assert.NoError(t, err)
		date := time.Now()

		user1 := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       random.Numeric(t, 2) == 1,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		user2 := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 20),
			Admin:       random.Numeric(t, 2) == 1,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		users := []*model.User{user1, user2}

		ctx := context.Background()
		th.Repository.MockUserRepository.
			EXPECT().
			GetUsers(ctx).
			Return(users, nil)

		var resBody UserResponse
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.GET, "/api/users", nil, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		if assert.Len(t, resBody.Users, 2) {
			if resBody.Users[0].Name == user1.Name {
				assert.Equal(t, resBody.Users[0].Name, user1.Name)
				assert.Equal(t, resBody.Users[0].DisplayName, user1.DisplayName)
				assert.Equal(t, resBody.Users[0].Admin, user1.Admin)
				assert.Equal(t, resBody.Users[1].Name, user2.Name)
				assert.Equal(t, resBody.Users[1].DisplayName, user2.DisplayName)
				assert.Equal(t, resBody.Users[1].Admin, user2.Admin)
			} else {
				assert.Equal(t, resBody.Users[1].Name, user1.Name)
				assert.Equal(t, resBody.Users[1].DisplayName, user1.DisplayName)
				assert.Equal(t, resBody.Users[1].Admin, user1.Admin)
				assert.Equal(t, resBody.Users[0].Name, user2.Name)
				assert.Equal(t, resBody.Users[0].DisplayName, user2.DisplayName)
				assert.Equal(t, resBody.Users[0].Admin, user2.Admin)
			}
		}
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
		assert.NoError(t, err)

		users := []*model.User{}

		ctx := context.Background()
		th.Repository.MockUserRepository.
			EXPECT().
			GetUsers(ctx).
			Return(users, nil)

		var resBody UserResponse
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.GET, "/api/users", nil, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Len(t, resBody.Users, 0)
	})

	t.Run("FailedToGetUsers", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
		assert.NoError(t, err)

		ctx := context.Background()
		th.Repository.MockUserRepository.
			EXPECT().
			GetUsers(ctx).
			Return(nil, errors.New("failed to get users."))

		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.GET, "/api/users", nil, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}

func TestHandlers_PutUser(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		user := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
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

		req := &PostUser{
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
			Return(updateUser, nil)

		var resBody UserOverview
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.PUT, "/api/users", &req, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, updateUser.Name, resBody.Name)
		assert.Equal(t, updateUser.DisplayName, resBody.DisplayName)
		assert.Equal(t, updateUser.Admin, resBody.Admin)
	})

	t.Run("FailedToUpdateUser", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		user := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
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

		req := &PostUser{
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
		accessUser := mustMakeUser(t, false)
		user := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
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

		req := &PostUser{
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

func TestHandlers_GetMe(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
		assert.NoError(t, err)

		ctx := context.Background()
		th.Repository.MockUserRepository.
			EXPECT().
			GetUserByName(ctx, accessUser.Name).
			Return(accessUser, nil)

		var resBody UserOverview
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.GET, "/api/users/me", nil, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, accessUser.Name, resBody.Name)
		assert.Equal(t, accessUser.DisplayName, resBody.DisplayName)
		assert.Equal(t, accessUser.Admin, resBody.Admin)
	})

	t.Run("FailedToGetUser", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := mustMakeUser(t, false)
		th, err := SetupTestHandlers(t, ctrl, accessUser)
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
