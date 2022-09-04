package model

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetRequests(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_requests")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "get_requests2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)
	client3, storage3, err := setup(t, ctx, "get_requests3")
	require.NoError(t, err)
	repo3 := NewEntRepository(client3, storage3)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user1, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		user2, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)

		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		request1, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user1.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		request2, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user2.ID)
		require.NoError(t, err)

		sort := "created_at"

		got, err := repo.GetRequests(ctx, RequestQuery{
			Sort: &sort,
		})
		assert.NoError(t, err)
		if assert.Len(t, got, 2) && assert.Equal(t, got[1].ID, request1.ID) {
			assert.Equal(t, got[1].ID, request1.ID)
			assert.Equal(t, got[1].Status, request1.Status)
			assert.Equal(t, got[1].Amount, request1.Amount)
			assert.Equal(t, got[1].Title, request1.Title)
			assert.Equal(t, got[1].Content, request1.Content)
			assert.Equal(t, got[1].Tags[0].ID, request1.Tags[0].ID)
			assert.Equal(t, got[1].Tags[0].Name, request1.Tags[0].Name)
			assert.Equal(t, got[1].Tags[0].Description, request1.Tags[0].Description)
			assert.Equal(t, got[0].ID, request2.ID)
			assert.Equal(t, got[0].Status, request2.Status)
			assert.Equal(t, got[0].Amount, request2.Amount)
			assert.Equal(t, got[0].Title, request2.Title)
			assert.Equal(t, got[0].Content, request2.Content)
			assert.Equal(t, got[0].Tags[0].ID, request1.Tags[0].ID)
			assert.Equal(t, got[0].Tags[0].Name, request1.Tags[0].Name)
			assert.Equal(t, got[0].Tags[0].Description, request1.Tags[0].Description)
		}
	})

	t.Run("SuccessWithReverseOrder", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, err)
		user1, err := repo2.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		user2, err := repo2.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo2.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)

		budget := random.Numeric(t, 10000)
		group, err := repo2.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		request1, err := repo2.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user1.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		request2, err := repo2.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user2.ID)
		require.NoError(t, err)

		sort := "-created_at"

		got, err := repo2.GetRequests(ctx, RequestQuery{
			Sort: &sort,
		})
		assert.NoError(t, err)
		if assert.Len(t, got, 2) && assert.Equal(t, got[0].ID, request1.ID) {
			assert.Equal(t, got[0].ID, request1.ID)
			assert.Equal(t, got[0].Status, request1.Status)
			assert.Equal(t, got[0].Amount, request1.Amount)
			assert.Equal(t, got[0].Title, request1.Title)
			assert.Equal(t, got[0].Content, request1.Content)
			assert.Equal(t, got[0].Tags[0].ID, request1.Tags[0].ID)
			assert.Equal(t, got[0].Tags[0].Name, request1.Tags[0].Name)
			assert.Equal(t, got[0].Tags[0].Description, request1.Tags[0].Description)
			assert.Equal(t, got[1].ID, request2.ID)
			assert.Equal(t, got[1].Status, request2.Status)
			assert.Equal(t, got[1].Amount, request2.Amount)
			assert.Equal(t, got[1].Title, request2.Title)
			assert.Equal(t, got[1].Content, request2.Content)
			assert.Equal(t, got[1].Tags[0].ID, request1.Tags[0].ID)
			assert.Equal(t, got[1].Tags[0].Name, request1.Tags[0].Name)
			assert.Equal(t, got[1].Tags[0].Description, request1.Tags[0].Description)
		}
	})

	t.Run("SuccessWithStatus", func(t *testing.T) {
		t.Parallel()
		user1, err := repo3.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		user2, err := repo3.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo3.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)

		budget := random.Numeric(t, 10000)
		group, err := repo3.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		request1, err := repo3.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user1.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		request2, err := repo3.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user2.ID)
		require.NoError(t, err)

		time.Sleep(1 * time.Second)

		status := "accepted"
		_, err = repo3.CreateStatus(ctx, request2.ID, user2.ID, Accepted)
		require.NoError(t, err)

		got, err := repo3.GetRequests(ctx, RequestQuery{
			Status: &status,
		})
		assert.NoError(t, err)
		if assert.Len(t, got, 1) {
			assert.Equal(t, got[0].ID, request2.ID)
			assert.Equal(t, got[0].Status, Accepted)
			assert.Equal(t, got[0].Amount, request2.Amount)
			assert.Equal(t, got[0].Title, request2.Title)
			assert.Equal(t, got[0].Content, request2.Content)
			assert.Equal(t, got[0].Tags[0].ID, request1.Tags[0].ID)
			assert.Equal(t, got[0].Tags[0].Name, request1.Tags[0].Name)
			assert.Equal(t, got[0].Tags[0].Description, request1.Tags[0].Description)
		}
	})
}

func TestEntRepository_CreateRequest(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "create_request")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		amount := random.Numeric(t, 100000)
		title := random.AlphaNumeric(t, 40)
		content := random.AlphaNumeric(t, 100)
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)

		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		request, err := repo.CreateRequest(ctx, amount, title, content, []*Tag{tag}, group, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, request.CreatedBy, user.ID)
		assert.Equal(t, request.Status, Status(1))
		assert.Equal(t, request.Amount, amount)
		assert.Equal(t, request.Title, title)
		assert.Equal(t, request.Content, content)
		assert.Equal(t, request.Tags, []*Tag{tag})
		assert.Equal(t, request.Group, group)
	})

	t.Run("UnknownUser", func(t *testing.T) {
		t.Parallel()
		amount := random.Numeric(t, 100000)
		title := random.AlphaNumeric(t, 40)
		content := random.AlphaNumeric(t, 100)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)

		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		_, err = repo.CreateRequest(ctx, amount, title, content, []*Tag{tag}, group, uuid.New())
		assert.Error(t, err)
	})

	t.Run("UnknownTag", func(t *testing.T) {
		t.Parallel()
		amount := random.Numeric(t, 100000)
		title := random.AlphaNumeric(t, 40)
		content := random.AlphaNumeric(t, 100)
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)

		date := time.Now()
		tag := &Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 30),
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		_, err = repo.CreateRequest(ctx, amount, title, content, []*Tag{tag}, group, user.ID)
		assert.Error(t, err)
	})

	t.Run("UnknownGroup", func(t *testing.T) {
		t.Parallel()
		amount := random.Numeric(t, 100000)
		title := random.AlphaNumeric(t, 40)
		content := random.AlphaNumeric(t, 100)
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)

		date := time.Now()
		budget := random.Numeric(t, 100000)
		group := &Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 20),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		_, err = repo.CreateRequest(ctx, amount, title, content, []*Tag{tag}, group, user.ID)
		assert.Error(t, err)
	})
}

func TestEntRepository_GetRequest(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_request")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user.ID)
		require.NoError(t, err)

		got, err := repo.GetRequest(ctx, request.ID)
		assert.NoError(t, err)
		assert.Equal(t, got.CreatedBy, user.ID)
		assert.Equal(t, got.Status, Status(1))
		assert.Equal(t, got.Amount, request.Amount)
		assert.Equal(t, got.Title, request.Title)
		assert.Equal(t, got.Content, request.Content)
		assert.Equal(t, got.Tags[0].ID, request.Tags[0].ID)
		assert.Equal(t, got.Tags[0].Name, request.Tags[0].Name)
		assert.Equal(t, got.Tags[0].Description, request.Tags[0].Description)
		assert.Equal(t, got.Group.ID, request.Group.ID)
		assert.Equal(t, got.Group.Name, request.Group.Name)
		assert.Equal(t, got.Group.Description, request.Group.Description)
		assert.Equal(t, got.Group.Budget, request.Group.Budget)
	})

	t.Run("UnknownRequest", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetRequest(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestEntREpository_UpdateRequest(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "update_request")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user.ID)
		require.NoError(t, err)

		amount := random.Numeric(t, 100000)
		updatedRequest, err := repo.UpdateRequest(ctx, request.ID, amount, request.Title, request.Content, []*Tag{tag}, group)
		assert.NoError(t, err)
		assert.Equal(t, updatedRequest.ID, request.ID)
		assert.Equal(t, updatedRequest.Status, request.Status)
		assert.Equal(t, updatedRequest.Amount, amount)
		assert.Equal(t, updatedRequest.Title, request.Title)
		assert.Equal(t, updatedRequest.Content, request.Content)
		assert.Equal(t, updatedRequest.Tags[0].ID, tag.ID)
		assert.Equal(t, updatedRequest.Tags[0].Name, tag.Name)
		assert.Equal(t, updatedRequest.Tags[0].Description, tag.Description)
		assert.Equal(t, updatedRequest.Group.ID, request.Group.ID)
		assert.Equal(t, updatedRequest.Group.Name, request.Group.Name)
		assert.Equal(t, updatedRequest.Group.Description, request.Group.Description)
		assert.Equal(t, updatedRequest.Group.Budget, request.Group.Budget)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user.ID)
		require.NoError(t, err)

		title := random.AlphaNumeric(t, 40)
		updatedRequest, err := repo.UpdateRequest(ctx, request.ID, request.Amount, title, request.Content, []*Tag{tag}, group)
		assert.NoError(t, err)
		assert.Equal(t, updatedRequest.ID, request.ID)
		assert.Equal(t, updatedRequest.Status, request.Status)
		assert.Equal(t, updatedRequest.Amount, request.Amount)
		assert.Equal(t, updatedRequest.Title, title)
		assert.Equal(t, updatedRequest.Content, request.Content)
		assert.Equal(t, updatedRequest.Tags[0].ID, tag.ID)
		assert.Equal(t, updatedRequest.Tags[0].Name, tag.Name)
		assert.Equal(t, updatedRequest.Tags[0].Description, tag.Description)
		assert.Equal(t, updatedRequest.Group.ID, request.Group.ID)
		assert.Equal(t, updatedRequest.Group.Name, request.Group.Name)
		assert.Equal(t, updatedRequest.Group.Description, request.Group.Description)
		assert.Equal(t, updatedRequest.Group.Budget, request.Group.Budget)
	})

	t.Run("Success3", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user.ID)
		require.NoError(t, err)

		content := random.AlphaNumeric(t, 100)
		updatedRequest, err := repo.UpdateRequest(ctx, request.ID, request.Amount, request.Title, content, []*Tag{tag}, group)
		assert.NoError(t, err)
		assert.Equal(t, updatedRequest.ID, request.ID)
		assert.Equal(t, updatedRequest.Status, request.Status)
		assert.Equal(t, updatedRequest.Amount, request.Amount)
		assert.Equal(t, updatedRequest.Title, request.Title)
		assert.Equal(t, updatedRequest.Content, content)
		assert.Equal(t, updatedRequest.Tags[0].ID, tag.ID)
		assert.Equal(t, updatedRequest.Tags[0].Name, tag.Name)
		assert.Equal(t, updatedRequest.Tags[0].Description, tag.Description)
		assert.Equal(t, updatedRequest.Group.ID, request.Group.ID)
		assert.Equal(t, updatedRequest.Group.Name, request.Group.Name)
		assert.Equal(t, updatedRequest.Group.Description, request.Group.Description)
		assert.Equal(t, updatedRequest.Group.Budget, request.Group.Budget)
	})

	t.Run("UnknownTag", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user.ID)
		require.NoError(t, err)

		date := time.Now()
		unknownTag := &Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 30),
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		_, err = repo.UpdateRequest(ctx, request.ID, request.Amount, request.Title, request.Content, []*Tag{unknownTag}, group)
		assert.Error(t, err)
	})

	t.Run("UnknownGroup", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user.ID)
		require.NoError(t, err)

		date := time.Now()
		unknownBudget := random.Numeric(t, 100000)
		unknownGroup := &Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 30),
			Budget:      &unknownBudget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		_, err = repo.UpdateRequest(ctx, request.ID, request.Amount, request.Title, request.Content, []*Tag{tag}, unknownGroup)
		assert.Error(t, err)
	})
}
