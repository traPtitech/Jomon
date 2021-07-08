package router

import (
	"context"
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

// Testing GetTags Method
func TestHandlers_GetTags(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		assert.NoError(t, err)
		date := time.Now()

		tag1 := &model.Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		tag2 := &model.Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		tags := []*model.Tag{tag1, tag2}

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			GetTags(ctx).
			Return(tags, nil)

		statusCode, _ := th.doRequest(t, echo.GET, "/api/tags", nil, nil)
		assert.Equal(t, http.StatusOK, statusCode)
	})
}
