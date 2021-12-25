package model

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetRequests(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx)
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user1, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		user2, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)

		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget, &[]User{*owner})
		require.NoError(t, err)

		request1, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user1.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		request2, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user2.ID)
		require.NoError(t, err)

		sort := "created_at"
		year := 2021
		since := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local)
		until := time.Now()

		got, err := repo.GetRequests(ctx, RequestQuery{
			Sort: &sort,
			// Target: &target,
			Year:  &year,
			Since: &since,
			Until: &until,
			Tag:   &tag.Name,
			Group: &group.Name,
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
		user1, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		user2, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)

		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget, &[]User{*owner})
		require.NoError(t, err)

		request1, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user1.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		request2, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user2.ID)
		require.NoError(t, err)


		sort := "-created_at"
		year := 2021
		since := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local)
		until := time.Now()

		got, err := repo.GetRequests(ctx, RequestQuery{
			Sort: &sort,
			// Target: &target,
			Year:  &year,
			Since: &since,
			Until: &until,
			Tag:   &tag.Name,
			Group: &group.Name,
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
}

func TestEntRepository_CreateRequest(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx)
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

		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget, &[]User{*owner})
		require.NoError(t, err)

		request1, err := repo.CreateRequest(ctx, amount, title, content, []*Tag{}, group, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, request1.CreatedBy, user.ID)
		assert.Equal(t, request1.Amount, amount)
		assert.Equal(t, request1.Title, title)
		assert.Equal(t, request1.Content, content)
		assert.Equal(t, request1.Tags, []*Tag{})
		assert.Equal(t, request1.Group, group)
		request2, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, request2.CreatedBy, user.ID)
		assert.Equal(t, request2.Tags, []*Tag{tag})

	})
}

func TestEntRepository_GetRequest(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx)
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)
		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget, &[]User{*owner})
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user.ID)
		require.NoError(t, err)

		got, err := repo.GetRequest(ctx, request.ID)
		assert.NoError(t, err)
		assert.Equal(t, got.CreatedBy, user.ID)
		assert.Equal(t, got.Tags[0].ID, tag.ID)
		assert.Equal(t, got.Tags[0].Name, tag.Name)
		assert.Equal(t, got.Tags[0].Description, tag.Description)
		assert.Equal(t, got.Group.ID, group.ID)
		assert.Equal(t, got.Group.Name, group.Name)
		assert.Equal(t, got.Group.Description, group.Description)
		assert.Equal(t, *got.Group.Budget, *group.Budget)
	})
}

func TestEntREpository_UpdateRequest(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx)
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user1, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)
		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget, &[]User{*owner})
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{}, group, user1.ID)
		require.NoError(t, err)

		amount := 100
		updatedRequest, err := repo.UpdateRequest(ctx, request.ID, amount, request.Title, request.Content, []*Tag{tag}, group)
		assert.NoError(t, err)
		assert.Equal(t, updatedRequest.ID, request.ID)
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

}