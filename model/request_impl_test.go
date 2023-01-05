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
	client4, storage4, err := setup(t, ctx, "get_requests4")
	require.NoError(t, err)
	repo4 := NewEntRepository(client4, storage4)
	client5, storage5, err := setup(t, ctx, "get_requests5")
	require.NoError(t, err)
	repo5 := NewEntRepository(client5, storage5)
	client6, storage6, err := setup(t, ctx, "get_requests6")
	require.NoError(t, err)
	repo6 := NewEntRepository(client6, storage6)
	client7, storage7, err := setup(t, ctx, "get_requests7")
	require.NoError(t, err)
	repo7 := NewEntRepository(client7, storage7)
	client8, storage8, err := setup(t, ctx, "get_requests8")
	require.NoError(t, err)
	repo8 := NewEntRepository(client8, storage8)

	t.Run("SuccessWithSortCreatedAt", func(t *testing.T) {
		t.Parallel()
		user1, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		user2, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		request1, err := repo.CreateRequest(ctx, random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user1.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		request2, err := repo.CreateRequest(ctx, random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user2.ID)
		require.NoError(t, err)

		sort := "created_at"

		got, err := repo.GetRequests(ctx, RequestQuery{
			Sort: &sort,
		})
		assert.NoError(t, err)
		if assert.Len(t, got, 2) && assert.Equal(t, got[1].ID, request1.ID) {
			assert.Equal(t, got[1].ID, request1.ID)
			assert.Equal(t, got[1].Status, request1.Status)
			assert.Equal(t, got[1].Title, request1.Title)
			assert.Equal(t, got[1].Content, request1.Content)
			assert.Equal(t, got[1].Tags[0].ID, request1.Tags[0].ID)
			assert.Equal(t, got[1].Tags[0].Name, request1.Tags[0].Name)
			assert.Equal(t, got[0].ID, request2.ID)
			assert.Equal(t, got[0].Status, request2.Status)
			assert.Equal(t, got[0].Title, request2.Title)
			assert.Equal(t, got[0].Content, request2.Content)
			assert.Equal(t, got[0].Tags[0].ID, request1.Tags[0].ID)
			assert.Equal(t, got[0].Tags[0].Name, request1.Tags[0].Name)
		}
	})

	t.Run("SuccessWithReverseSortCreatedAt", func(t *testing.T) {
		t.Parallel()
		user1, err := repo2.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		user2, err := repo2.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo2.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo2.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		request1, err := repo2.CreateRequest(ctx, random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user1.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		request2, err := repo2.CreateRequest(ctx, random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user2.ID)
		require.NoError(t, err)

		sort := "-created_at"

		got, err := repo2.GetRequests(ctx, RequestQuery{
			Sort: &sort,
		})
		assert.NoError(t, err)
		if assert.Len(t, got, 2) && assert.Equal(t, got[0].ID, request1.ID) {
			assert.Equal(t, got[0].ID, request1.ID)
			assert.Equal(t, got[0].Status, request1.Status)
			assert.Equal(t, got[0].Title, request1.Title)
			assert.Equal(t, got[0].Content, request1.Content)
			assert.Equal(t, got[0].Tags[0].ID, request1.Tags[0].ID)
			assert.Equal(t, got[0].Tags[0].Name, request1.Tags[0].Name)
			assert.Equal(t, got[1].ID, request2.ID)
			assert.Equal(t, got[1].Status, request2.Status)
			assert.Equal(t, got[1].Title, request2.Title)
			assert.Equal(t, got[1].Content, request2.Content)
			assert.Equal(t, got[1].Tags[0].ID, request1.Tags[0].ID)
			assert.Equal(t, got[1].Tags[0].Name, request1.Tags[0].Name)
		}
	})

	t.Run("SuccessWithSortTitle", func(t *testing.T) {
		t.Parallel()
		user1, err := repo3.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		user2, err := repo3.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo3.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo3.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		request1, err := repo3.CreateRequest(ctx, "b", random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user1.ID)
		require.NoError(t, err)
		request2, err := repo3.CreateRequest(ctx, "a", random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user2.ID)
		require.NoError(t, err)

		sort := "title"

		got, err := repo3.GetRequests(ctx, RequestQuery{
			Sort: &sort,
		})
		assert.NoError(t, err)
		if assert.Len(t, got, 2) && assert.Equal(t, got[0].ID, request2.ID) {
			assert.Equal(t, got[0].ID, request2.ID)
			assert.Equal(t, got[0].Status, request2.Status)
			assert.Equal(t, got[0].Title, request2.Title)
			assert.Equal(t, got[0].Content, request2.Content)
			assert.Equal(t, got[0].Tags[0].ID, request2.Tags[0].ID)
			assert.Equal(t, got[0].Tags[0].Name, request2.Tags[0].Name)
			assert.Equal(t, got[1].ID, request1.ID)
			assert.Equal(t, got[1].Status, request1.Status)
			assert.Equal(t, got[1].Title, request1.Title)
			assert.Equal(t, got[1].Content, request1.Content)
			assert.Equal(t, got[1].Tags[0].ID, request1.Tags[0].ID)
			assert.Equal(t, got[1].Tags[0].Name, request1.Tags[0].Name)
		}
	})

	t.Run("SuccessWithReverseSortTitle", func(t *testing.T) {
		t.Parallel()
		user1, err := repo4.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		user2, err := repo4.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo4.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo4.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		request1, err := repo4.CreateRequest(ctx, "b", random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user1.ID)
		require.NoError(t, err)
		request2, err := repo4.CreateRequest(ctx, "a", random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user2.ID)
		require.NoError(t, err)

		sort := "-title"

		got, err := repo4.GetRequests(ctx, RequestQuery{
			Sort: &sort,
		})
		assert.NoError(t, err)
		if assert.Len(t, got, 2) && assert.Equal(t, got[0].ID, request1.ID) {
			assert.Equal(t, got[0].ID, request1.ID)
			assert.Equal(t, got[0].Status, request1.Status)
			assert.Equal(t, got[0].Title, request1.Title)
			assert.Equal(t, got[0].Content, request1.Content)
			assert.Equal(t, got[0].Tags[0].ID, request1.Tags[0].ID)
			assert.Equal(t, got[0].Tags[0].Name, request1.Tags[0].Name)
			assert.Equal(t, got[1].ID, request2.ID)
			assert.Equal(t, got[1].Status, request2.Status)
			assert.Equal(t, got[1].Title, request2.Title)
			assert.Equal(t, got[1].Content, request2.Content)
			assert.Equal(t, got[1].Tags[0].ID, request2.Tags[0].ID)
			assert.Equal(t, got[1].Tags[0].Name, request2.Tags[0].Name)
		}
	})

	t.Run("SuccessWithQueryTarget", func(t *testing.T) {
		t.Parallel()
		user1, err := repo5.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		user2, err := repo5.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo5.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target1 := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}
		target2 := &RequestTarget{
			Target: user2.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo5.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		request1, err := repo5.CreateRequest(ctx, "b", random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target1}, group, user1.ID)
		require.NoError(t, err)
		_, err = repo5.CreateRequest(ctx, "a", random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target2}, group, user2.ID)
		require.NoError(t, err)

		target := target1.Target
		got, err := repo5.GetRequests(ctx, RequestQuery{
			Target: &target,
		})
		assert.NoError(t, err)
		if assert.Len(t, got, 1) {
			assert.Equal(t, got[0].ID, request1.ID)
			assert.Equal(t, got[0].Status, request1.Status)
			assert.Equal(t, got[0].Title, request1.Title)
			assert.Equal(t, got[0].Content, request1.Content)
			assert.Equal(t, got[0].Tags[0].ID, request1.Tags[0].ID)
			assert.Equal(t, got[0].Tags[0].Name, request1.Tags[0].Name)
		}
	})

	t.Run("SuccessWithQuerySince", func(t *testing.T) {
		t.Parallel()
		user1, err := repo6.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		user2, err := repo6.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo6.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo6.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		request1, err := repo6.CreateRequest(ctx, "b", random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user1.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		request2, err := repo6.CreateRequest(ctx, "a", random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user2.ID)
		require.NoError(t, err)

		since := request1.CreatedAt.Add(10 * time.Millisecond)
		got, err := repo6.GetRequests(ctx, RequestQuery{
			Since: &since,
		})
		assert.NoError(t, err)
		if assert.Len(t, got, 1) {
			assert.Equal(t, got[0].ID, request2.ID)
			assert.Equal(t, got[0].Status, request2.Status)
			assert.Equal(t, got[0].Title, request2.Title)
			assert.Equal(t, got[0].Content, request2.Content)
			assert.Equal(t, got[0].Tags[0].ID, request2.Tags[0].ID)
			assert.Equal(t, got[0].Tags[0].Name, request2.Tags[0].Name)
		}
	})

	t.Run("SuccessWithQueryUntil", func(t *testing.T) {
		t.Parallel()
		user1, err := repo7.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		user2, err := repo7.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo7.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo7.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		request1, err := repo7.CreateRequest(ctx, "b", random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user1.ID)
		require.NoError(t, err)
		time.Sleep(2 * time.Second)
		request2, err := repo7.CreateRequest(ctx, "a", random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user2.ID)
		require.NoError(t, err)

		until := request2.CreatedAt.Add(-1 * time.Second)
		got, err := repo7.GetRequests(ctx, RequestQuery{
			Until: &until,
		})
		assert.NoError(t, err)
		if assert.Len(t, got, 1) {
			assert.Equal(t, got[0].ID, request1.ID)
			assert.Equal(t, got[0].Status, request1.Status)
			assert.Equal(t, got[0].Title, request1.Title)
			assert.Equal(t, got[0].Content, request1.Content)
			assert.Equal(t, got[0].Tags[0].ID, request1.Tags[0].ID)
			assert.Equal(t, got[0].Tags[0].Name, request1.Tags[0].Name)
		}
	})

	t.Run("SuccessWithQueryStatus", func(t *testing.T) {
		t.Parallel()
		user1, err := repo8.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		user2, err := repo8.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo8.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo8.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		request1, err := repo8.CreateRequest(ctx, "b", random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user1.ID)
		require.NoError(t, err)
		time.Sleep(2 * time.Second)
		_, err = repo8.CreateRequest(ctx, "a", random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user2.ID)
		require.NoError(t, err)

		time.Sleep(1 * time.Second)

		status := "accepted"
		_, err = repo8.CreateStatus(ctx, request1.ID, user1.ID, Accepted)
		require.NoError(t, err)

		got, err := repo8.GetRequests(ctx, RequestQuery{
			Status: &status,
		})
		assert.NoError(t, err)
		if assert.Len(t, got, 1) {
			assert.Equal(t, got[0].ID, request1.ID)
			assert.Equal(t, got[0].Status, Accepted)
			assert.Equal(t, got[0].Title, request1.Title)
			assert.Equal(t, got[0].Content, request1.Content)
			assert.Equal(t, got[0].Tags[0].ID, request1.Tags[0].ID)
			assert.Equal(t, got[0].Tags[0].Name, request1.Tags[0].Name)
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
		ctx := context.Background()
		title := random.AlphaNumeric(t, 40)
		content := random.AlphaNumeric(t, 100)
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		request, err := repo.CreateRequest(ctx, title, content, []*Tag{tag}, []*RequestTarget{target}, group, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, request.CreatedBy, user.ID)
		assert.Equal(t, request.Status, Status(1))
		assert.Equal(t, request.Title, title)
		assert.Equal(t, request.Content, content)
		assert.Equal(t, request.Tags, []*Tag{tag})
		assert.Equal(t, request.Group, group)
	})

	t.Run("UnknownUser", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		title := random.AlphaNumeric(t, 40)
		content := random.AlphaNumeric(t, 100)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)

		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		_, err = repo.CreateRequest(ctx, title, content, []*Tag{tag}, []*RequestTarget{}, group, uuid.New())
		assert.Error(t, err)
	})

	t.Run("UnknownTag", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		title := random.AlphaNumeric(t, 40)
		content := random.AlphaNumeric(t, 100)
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)

		date := time.Now()
		tag := &Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}

		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)

		_, err = repo.CreateRequest(ctx, title, content, []*Tag{tag}, []*RequestTarget{}, group, user.ID)
		assert.Error(t, err)
	})

	t.Run("UnknownGroup", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		title := random.AlphaNumeric(t, 40)
		content := random.AlphaNumeric(t, 100)
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
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

		_, err = repo.CreateRequest(ctx, title, content, []*Tag{tag}, []*RequestTarget{}, group, user.ID)
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
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user.ID)
		require.NoError(t, err)

		got, err := repo.GetRequest(ctx, request.ID)
		assert.NoError(t, err)
		assert.Equal(t, got.CreatedBy, user.ID)
		assert.Equal(t, got.Status, Status(1))
		assert.Equal(t, got.Title, request.Title)
		assert.Equal(t, got.Content, request.Content)
		assert.Equal(t, got.Tags[0].ID, request.Tags[0].ID)
		assert.Equal(t, got.Tags[0].Name, request.Tags[0].Name)
		assert.Equal(t, got.Targets[0].Target, request.Targets[0].Target)
		assert.Equal(t, got.Targets[0].Amount, request.Targets[0].Amount)
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

func TestEntRepository_UpdateRequest(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "update_request")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user.ID)
		require.NoError(t, err)

		updatedRequest, err := repo.UpdateRequest(ctx, request.ID, request.Title, request.Content, []*Tag{tag}, []*RequestTarget{target}, group)
		assert.NoError(t, err)
		assert.Equal(t, updatedRequest.ID, request.ID)
		assert.Equal(t, updatedRequest.Status, request.Status)
		assert.Equal(t, updatedRequest.Title, request.Title)
		assert.Equal(t, updatedRequest.Content, request.Content)
		assert.Equal(t, updatedRequest.Tags[0].ID, tag.ID)
		assert.Equal(t, updatedRequest.Tags[0].Name, tag.Name)
		assert.Equal(t, updatedRequest.Group.ID, request.Group.ID)
		assert.Equal(t, updatedRequest.Group.Name, request.Group.Name)
		assert.Equal(t, updatedRequest.Group.Description, request.Group.Description)
		assert.Equal(t, updatedRequest.Group.Budget, request.Group.Budget)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user.ID)
		require.NoError(t, err)

		title := random.AlphaNumeric(t, 40)
		updatedRequest, err := repo.UpdateRequest(ctx, request.ID, title, request.Content, []*Tag{tag}, []*RequestTarget{target}, group)
		assert.NoError(t, err)
		assert.Equal(t, updatedRequest.ID, request.ID)
		assert.Equal(t, updatedRequest.Status, request.Status)
		assert.Equal(t, updatedRequest.Title, title)
		assert.Equal(t, updatedRequest.Content, request.Content)
		assert.Equal(t, updatedRequest.Tags[0].ID, tag.ID)
		assert.Equal(t, updatedRequest.Tags[0].Name, tag.Name)
		assert.Equal(t, updatedRequest.Group.ID, request.Group.ID)
		assert.Equal(t, updatedRequest.Group.Name, request.Group.Name)
		assert.Equal(t, updatedRequest.Group.Description, request.Group.Description)
		assert.Equal(t, updatedRequest.Group.Budget, request.Group.Budget)
	})

	t.Run("Success3", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user.ID)
		require.NoError(t, err)
		content := random.AlphaNumeric(t, 100)
		updatedRequest, err := repo.UpdateRequest(ctx, request.ID, request.Title, content, []*Tag{tag}, []*RequestTarget{target}, group)
		assert.NoError(t, err)
		assert.Equal(t, updatedRequest.ID, request.ID)
		assert.Equal(t, updatedRequest.Status, request.Status)
		assert.Equal(t, updatedRequest.Title, request.Title)
		assert.Equal(t, updatedRequest.Content, content)
		assert.Equal(t, updatedRequest.Tags[0].ID, tag.ID)
		assert.Equal(t, updatedRequest.Tags[0].Name, tag.Name)
		assert.Equal(t, updatedRequest.Group.ID, request.Group.ID)
		assert.Equal(t, updatedRequest.Group.Name, request.Group.Name)
		assert.Equal(t, updatedRequest.Group.Description, request.Group.Description)
		assert.Equal(t, updatedRequest.Group.Budget, request.Group.Budget)
	})

	t.Run("UnknownTag", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user.ID)
		require.NoError(t, err)

		date := time.Now()
		unknownTag := &Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		_, err = repo.UpdateRequest(ctx, request.ID, request.Title, request.Content, []*Tag{unknownTag}, []*RequestTarget{target}, group)
		assert.Error(t, err)
	})

	t.Run("UnknownGroup", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, []*RequestTarget{target}, group, user.ID)
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
		_, err = repo.UpdateRequest(ctx, request.ID, request.Title, request.Content, []*Tag{tag}, []*RequestTarget{target}, unknownGroup)
		assert.Error(t, err)
	})
}
