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

// Test GET /api/requests. this test uses mock, so query tests are in model.
func TestHandlers_GetRequests(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		assert.NoError(t, err)
		date1 := time.Now()
		date2 := date1.Add(time.Hour)

		request1 := &model.RequestResponse{
			ID:        uuid.New(),
			Status:    model.Submitted.String(),
			CreatedBy: uuid.New(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date1,
			UpdatedAt: date1,
		}
		request2 := &model.RequestResponse{
			ID:        uuid.New(),
			Status:    model.Submitted.String(),
			CreatedBy: uuid.New(),
			Amount:    random.Numeric(t, 1000000),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			CreatedAt: date2,
			UpdatedAt: date2,
		}
		requests := []*model.RequestResponse{request1, request2}

		var title string
		var target string
		var year int
		var since time.Time
		var until time.Time
		var tag string
		var group string

		query := model.RequestQuery{
			Sort:   &title,
			Target: &target,
			Year:   &year,
			Since:  &since,
			Until:  &until,
			Tag:    &tag,
			Group:  &group,
		}

		ctx := context.Background()
		th.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(ctx, query).
			Return(requests, nil)

		var resBody []*RequestResponse
		statusCode, _ := th.doRequest(t, echo.GET, "/api/requests", nil, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Len(t, resBody, 2)
		assert.Equal(t, resBody[0].ID, request1.ID)
		assert.Equal(t, resBody[0].Status, request1.Status)
		assert.Equal(t, resBody[0].CreatedBy, request1.CreatedBy)
		assert.Equal(t, resBody[0].Amount, request1.Amount)
		assert.Equal(t, resBody[0].Title, request1.Title)
		assert.Equal(t, resBody[0].Content, request1.Content)
		assert.Equal(t, resBody[1].ID, request2.ID)
		assert.Equal(t, resBody[1].Status, request2.Status)
		assert.Equal(t, resBody[1].CreatedBy, request2.CreatedBy)
		assert.Equal(t, resBody[1].Amount, request2.Amount)
		assert.Equal(t, resBody[1].Title, request2.Title)
		assert.Equal(t, resBody[1].Content, request2.Content)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		assert.NoError(t, err)

		requests := []*model.RequestResponse{}

		var title string
		var target string
		var year int
		var since time.Time
		var until time.Time
		var tag string
		var group string

		query := model.RequestQuery{
			Sort:   &title,
			Target: &target,
			Year:   &year,
			Since:  &since,
			Until:  &until,
			Tag:    &tag,
			Group:  &group,
		}

		ctx := context.Background()
		th.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(ctx, query).
			Return(requests, nil)

		var resBody []*RequestResponse
		statusCode, _ := th.doRequest(t, echo.GET, "/api/requests", nil, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Len(t, resBody, 0)
	})

	t.Run("FailedToGetRequests", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		assert.NoError(t, err)

		var title string
		var target string
		var year int
		var since time.Time
		var until time.Time
		var tag string
		var group string

		query := model.RequestQuery{
			Sort:   &title,
			Target: &target,
			Year:   &year,
			Since:  &since,
			Until:  &until,
			Tag:    &tag,
			Group:  &group,
		}

		ctx := context.Background()
		th.Repository.MockRequestRepository.
			EXPECT().
			GetRequests(ctx, query).
			Return(nil, errors.New("Failed to get requests."))

		statusCode, _ := th.doRequest(t, echo.GET, "/api/requests", nil, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}

func TestHandlers_PostRequest(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		assert.NoError(t, err)

		req := &Request{
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			Amount:    random.Numeric(t, 1000000),
		}
		var tags []*model.Tag
		var group *model.Group

		date := time.Now()

		ctx := context.Background()
		th.Repository.MockRequestRepository.
			EXPECT().
			CreateRequest(ctx, req.Amount, req.Title, req.Content, tags, group, req.CreatedBy).
			Return(&model.RequestDetail{
				ID:        uuid.New(),
				Status:    model.Submitted.String(),
				Amount:    req.Amount,
				Title:     req.Title,
				Content:   req.Content,
				Comments:  []*model.Comment{},
				Files:     []*uuid.UUID{},
				Statuses:  []*model.RequestStatus{},
				Tags:      tags,
				Group:     group,
				CreatedAt: date,
				UpdatedAt: date,
				CreatedBy: req.CreatedBy,
			}, nil)

		var resBody RequestResponse
		statusCode, _ := th.doRequest(t, echo.POST, "/api/requests", &req, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, resBody.Title, req.Title)
		assert.Equal(t, resBody.Content, req.Content)
		assert.Equal(t, resBody.Amount, req.Amount)
		assert.Equal(t, resBody.CreatedBy, req.CreatedBy)
	})

	t.Run("SuccessWithTags", func(t *testing.T) {
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
		tags := []*model.Tag{tag1}
		var group *model.Group

		req := &Request{
			CreatedBy: uuid.New(),
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			Amount:    random.Numeric(t, 1000000),
			Tags:      []*uuid.UUID{&tag1.ID},
		}

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			GetTag(ctx, tag1.ID).
			Return(tag1, nil)
		th.Repository.MockRequestRepository.
			EXPECT().
			CreateRequest(ctx, req.Amount, req.Title, req.Content, tags, group, req.CreatedBy).
			Return(&model.RequestDetail{
				ID:        uuid.New(),
				Status:    model.Submitted.String(),
				Amount:    req.Amount,
				Title:     req.Title,
				Content:   req.Content,
				Comments:  []*model.Comment{},
				Files:     []*uuid.UUID{},
				Statuses:  []*model.RequestStatus{},
				Tags:      tags,
				Group:     group,
				CreatedAt: date,
				UpdatedAt: date,
				CreatedBy: req.CreatedBy,
			}, nil)

		var resBody RequestResponse
		statusCode, _ := th.doRequest(t, echo.POST, "/api/requests", &req, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, resBody.Title, req.Title)
		assert.Equal(t, resBody.Content, req.Content)
		assert.Equal(t, resBody.Amount, req.Amount)
		assert.Equal(t, resBody.CreatedBy, req.CreatedBy)
		if assert.Len(t, resBody.Tags, 1) {
			assert.Equal(t, resBody.Tags[0].ID, tag1.ID)
			assert.Equal(t, resBody.Tags[0].Name, tag1.Name)
			assert.Equal(t, resBody.Tags[0].Description, tag1.Description)
		}
	})
}
