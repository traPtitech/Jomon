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

func TestHandlers_GetGroups(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		assert.NoError(t, err)
		date := time.Now()

		budget1 := random.Numeric(t, 1000000)
		budget2 := random.Numeric(t, 1000000)

		group1 := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget1,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		group2 := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget2,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		groups := []*model.Group{group1, group2}

		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			GetGroups(ctx).
			Return(groups, nil)

		var resBody []*GroupOverview
		statusCode, _ := th.doRequest(t, echo.GET, "/api/groups", nil, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Len(t, resBody, 2)
		if resBody[0].ID == group1.ID {
			assert.Equal(t, group1.ID, resBody[0].ID)
			assert.Equal(t, group1.Name, resBody[0].Name)
			assert.Equal(t, group1.Description, resBody[0].Description)
			assert.Equal(t, group1.Budget, resBody[0].Budget)
			assert.Equal(t, group2.ID, resBody[1].ID)
			assert.Equal(t, group2.Name, resBody[1].Name)
			assert.Equal(t, group2.Description, resBody[1].Description)
			assert.Equal(t, group2.Budget, resBody[1].Budget)
		} else {
			assert.Equal(t, group2.ID, resBody[0].ID)
			assert.Equal(t, group2.Name, resBody[0].Name)
			assert.Equal(t, group2.Description, resBody[0].Description)
			assert.Equal(t, group2.Budget, resBody[0].Budget)
			assert.Equal(t, group1.ID, resBody[1].ID)
			assert.Equal(t, group1.Name, resBody[1].Name)
			assert.Equal(t, group1.Description, resBody[1].Description)
			assert.Equal(t, group1.Budget, resBody[1].Budget)
		}
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		assert.NoError(t, err)

		groups := []*model.Group{}

		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			GetGroups(ctx).
			Return(groups, nil)

		var resBody []*GroupOverview
		statusCode, _ := th.doRequest(t, echo.GET, "/api/groups", nil, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Len(t, resBody, 0)
	})

	t.Run("FailedToGetGroups", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		assert.NoError(t, err)

		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			GetGroups(ctx).
			Return(nil, errors.New("failed to get groups"))

		statusCode, _ := th.doRequest(t, echo.GET, "/api/groups", nil, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}
