package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestHandlers_GetGroups(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
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
		require.NoError(t, err)

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
		require.NoError(t, err)

		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			GetGroups(ctx).
			Return(nil, errors.New("Failed to get groups"))

		statusCode, _ := th.doRequest(t, echo.GET, "/api/groups", nil, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}

func TestHandlers_GetMembers(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		user1 := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}
		user2 := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}
		members := []*model.User{user1, user2}

		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			GetMembers(ctx, group.ID).
			Return(members, nil)

		var resBody MemberResponse
		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		statusCode, _ := th.doRequest(t, echo.GET, path, nil, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Len(t, resBody.ID, 2)
		if resBody.ID[0] == user1.ID {
			assert.Equal(t, resBody.ID[0], user1.ID)
			assert.Equal(t, resBody.ID[1], user2.ID)
		} else {
			assert.Equal(t, resBody.ID[0], user2.ID)
			assert.Equal(t, resBody.ID[1], user1.ID)
		}
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			GetMembers(ctx, group.ID).
			Return([]*model.User{}, nil)

		var resBody MemberResponse
		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		statusCode, _ := th.doRequest(t, echo.GET, path, nil, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Len(t, resBody.ID, 0)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)

		path := "/api/groups/hoge/members" // Invalid UUID
		statusCode, _ := th.doRequest(t, echo.GET, path, nil, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)

		path := fmt.Sprintf("/api/groups/%s/members", uuid.Nil.String())
		statusCode, _ := th.doRequest(t, echo.GET, path, nil, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("UnknownID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)

		unknownID := uuid.New()
		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			GetMembers(ctx, unknownID).
			Return(nil, errors.New("Group not found"))

		path := fmt.Sprintf("/api/groups/%s/members", unknownID.String())
		statusCode, _ := th.doRequest(t, echo.GET, path, nil, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}

func TestHndlers_PostMember(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			CreateMember(ctx, group.ID, user.ID).
			Return(&model.Member{
				ID: user.ID,
			}, nil)

		req := Member{
			ID: user.ID,
		}

		var resBody Member
		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		statusCode, _ := th.doRequest(t, echo.POST, path, &req, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, user.ID, resBody.ID)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		req := Member{
			ID: user.ID,
		}

		path := "/api/groups/hoge/members" // Invalid UUID
		statusCode, _ := th.doRequest(t, echo.POST, path, &req, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		req := Member{
			ID: user.ID,
		}

		path := fmt.Sprintf("/api/groups/%s/members", uuid.Nil.String())
		statusCode, _ := th.doRequest(t, echo.POST, path, &req, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("UnkonwnGroupID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		unkonwnGroupID := uuid.New()
		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			CreateMember(ctx, unkonwnGroupID, user.ID).
			Return(nil, errors.New("unknown group id"))

		req := Member{
			ID: user.ID,
		}

		path := fmt.Sprintf("/api/groups/%s/members", unkonwnGroupID.String())
		statusCode, _ := th.doRequest(t, echo.POST, path, &req, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})

	t.Run("UnkonwnUserID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		unkonwnUserID := uuid.New()
		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			CreateMember(ctx, group.ID, unkonwnUserID).
			Return(nil, errors.New("unknown user id"))

		req := Member{
			ID: unkonwnUserID,
		}

		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		statusCode, _ := th.doRequest(t, echo.POST, path, &req, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}

func TestHandlers_DeleteMember(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			CreateMember(ctx, group.ID, user.ID).
			Return(&model.Member{
				ID: user.ID,
				}, nil)

		req := Member{
			ID: user.ID,
		}

		var resBody Member
		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		statusCode, _ := th.doRequest(t, echo.POST, path, &req, &resBody)
		require.Equal(t, http.StatusOK, statusCode)

		th.Repository.MockGroupRepository.
			EXPECT().
			DeleteMember(ctx, group.ID, user.ID).
			Return(nil)

		statusCode2, _ := th.doRequest(t, echo.DELETE, path, &req, nil)
		assert.Equal(t, http.StatusOK, statusCode2)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			CreateMember(ctx, group.ID, user.ID).
			Return(&model.Member{user.ID}, nil)

		req := Member{
			ID: user.ID,
		}

		var resBody Member
		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		statusCode, _ := th.doRequest(t, echo.POST, path, &req, &resBody)
		require.Equal(t, http.StatusOK, statusCode)

		path2 := "/api/groups/hoge/members" // Invalid UUID
		statusCode2, _ := th.doRequest(t, echo.DELETE, path2, &req, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode2)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			CreateMember(ctx, group.ID, user.ID).
			Return(&model.Member{user.ID}, nil)

		req := Member{
			ID: user.ID,
		}

		var resBody Member
		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		statusCode, _ := th.doRequest(t, echo.POST, path, &req, &resBody)
		require.Equal(t, http.StatusOK, statusCode)

		path2 := fmt.Sprintf("/api/groups/%s/members", uuid.Nil.String())
		statusCode2, _ := th.doRequest(t, echo.DELETE, path2, &req, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode2)
	})

	t.Run("UnknownGroupID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		user := &model.User{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			DisplayName: random.AlphaNumeric(t, 50),
			Admin:       true,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			CreateMember(ctx, group.ID, user.ID).
			Return(&model.Member{user.ID}, nil)

		req := Member{
			ID: user.ID,
		}

		var resBody Member
		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		statusCode, _ := th.doRequest(t, echo.POST, path, &req, &resBody)
		require.Equal(t, http.StatusOK, statusCode)

		unknownGroupID := uuid.New()
		th.Repository.MockGroupRepository.
			EXPECT().
			DeleteMember(ctx, unknownGroupID, user.ID).
			Return(errors.New("unkown group id"))

		path2 := fmt.Sprintf("/api/groups/%s/members", unknownGroupID.String())
		statusCode2, _ := th.doRequest(t, echo.DELETE, path2, &req, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode2)
	})

	t.Run("UnknownMemberID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		th, err := SetupTestHandlers(t, ctrl)
		require.NoError(t, err)
		date := time.Now()

		budget := random.Numeric(t, 1000000)
		group := &model.Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
			DeletedAt:   &date,
		}

		unknownMemberID := uuid.New()
		ctx := context.Background()
		th.Repository.MockGroupRepository.
			EXPECT().
			DeleteMember(ctx, group.ID, unknownMemberID).
			Return(errors.New("unkown member id"))

		req := Member{
			ID: unknownMemberID,
		}

		path := fmt.Sprintf("/api/groups/%s/members", group.ID.String())
		statusCode, _ := th.doRequest(t, echo.DELETE, path, &req, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})

}
