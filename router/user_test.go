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
		th, err := SetupTestHandlers(t, ctrl)
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
		statusCode, _ := th.doRequest(t, echo.GET, "/api/users", nil, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Len(t, resBody.Users, 2)
		if resBody.Users[0].ID == user1.ID {
			assert.Equal(t, resBody.Users[0].ID, user1.ID)
			assert.Equal(t, resBody.Users[0].Name, user1.Name)
			assert.Equal(t, resBody.Users[0].DisplayName, user1.DisplayName)
			assert.Equal(t, resBody.Users[0].Admin, user1.Admin)
			assert.Equal(t, resBody.Users[1].ID, user2.ID)
			assert.Equal(t, resBody.Users[1].Name, user2.Name)
			assert.Equal(t, resBody.Users[1].DisplayName, user2.DisplayName)
			assert.Equal(t, resBody.Users[1].Admin, user2.Admin)
		} else {
			assert.Equal(t, resBody.Users[1].ID, user1.ID)
			assert.Equal(t, resBody.Users[1].Name, user1.Name)
			assert.Equal(t, resBody.Users[1].DisplayName, user1.DisplayName)
			assert.Equal(t, resBody.Users[1].Admin, user1.Admin)
			assert.Equal(t, resBody.Users[0].ID, user2.ID)
			assert.Equal(t, resBody.Users[0].Name, user2.Name)
			assert.Equal(t, resBody.Users[0].DisplayName, user2.DisplayName)
			assert.Equal(t, resBody.Users[0].Admin, user2.Admin)
		}
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		assert.NoError(t, err)

		users := []*model.User{}

		ctx := context.Background()
		th.Repository.MockUserRepository.
			EXPECT().
			GetUsers(ctx).
			Return(users, nil)

		var resBody TagResponse
		statusCode, _ := th.doRequest(t, echo.GET, "/api/tags", nil, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Len(t, resBody.Tags, 0)
	})

	t.Run("FailedToGetTags", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		assert.NoError(t, err)

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			GetTags(ctx).
			Return(nil, errors.New("Failed to get tags."))

		statusCode, _ := th.doRequest(t, echo.GET, "/api/tags", nil, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
