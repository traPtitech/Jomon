package model

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil/random"
)

// Targetだけ何もいれないか、初期値をいれるとinvalid memory address or nil pointer dereference になる。仕組み不明
// implのtarget覗いてみたらなんかランダムでいれてもinvalid memory address or nil pointer dereference

// func TestEntRepository_GetRequests(t *testing.T) {
// 	ctx := context.Background()
// 	client, storage, err := setup(t, ctx)
// 	require.NoError(t, err)
// 	repo := NewEntRepository(client, storage)

// 	t.Run("Success", func(t *testing.T) {
// 		t.Parallel()
// 		user1, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
// 		require.NoError(t, err)
// 		// user2, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
// 		// require.NoError(t, err)
// 		tag1, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
// 		require.NoError(t, err)
// 		// tag2, _ := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
// 		request1, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag1}, nil, user1.ID)
// 		require.NoError(t, err)
// 		// request2, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag2}, nil, user2.ID)
// 		// require.NoError(t, err)

// 		// sort := "created_at"
// 		target := random.AlphaNumeric(t, 20)
// 		year := 2021
// 		// since := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local)
// 		// until := time.Now()
// 		var title string

// 		var since time.Time
// 		var until time.Time
// 		var tag string
// 		var group string

// 		got, err := repo.GetRequests(ctx, RequestQuery{
// 			Sort:   &title,
// 			Target: &target,
// 			Year:   &year,
// 			Since:  &since,
// 			Until:  &until,
// 			Tag:    &tag,
// 			Group:  &group,
// 		})
// 		assert.NoError(t, err)
// 		if assert.Len(t, got, 1) && got[0].ID == request1.ID {
// 			assert.Equal(t, got[0].ID, request1.ID)
// 			assert.Equal(t, got[0].Status, request1.Status)
// 			assert.Equal(t, got[0].Amount, request1.Amount)
// 			assert.Equal(t, got[0].Title, request1.Title)
// 			assert.Equal(t, got[0].Content, request1.Content)
// 			assert.Equal(t, got[0].Tags, request1.Tags)
// 			// 	assert.Equal(t, got[1].ID, request2.ID)
// 			// 	assert.Equal(t, got[1].Status, request2.Status)
// 			// 	assert.Equal(t, got[1].Amount, request2.Amount)
// 			// 	assert.Equal(t, got[1].Title, request2.Title)
// 			// 	assert.Equal(t, got[1].Content, request2.Content)
// 			// 	assert.Equal(t, got[1].Tags, request2.Tags)
// 			// }else if assert.Len(t, got, 2) {
// 			//   assert.Equal(t, got[0].ID, request2.ID)
// 			// 	assert.Equal(t, got[0].Status, request2.Status)
// 			// 	assert.Equal(t, got[0].Amount, request2.Amount)
// 			// 	assert.Equal(t, got[0].Title, request2.Title)
// 			// 	assert.Equal(t, got[0].Content, request2.Content)
// 			// 	assert.Equal(t, got[0].Tags, request2.Tags)
// 			// 	assert.Equal(t, got[1].ID, request1.ID)
// 			// 	assert.Equal(t, got[1].Status, request1.Status)
// 			// 	assert.Equal(t, got[1].Amount, request1.Amount)
// 			// 	assert.Equal(t, got[1].Title, request1.Title)
// 			// 	assert.Equal(t, got[1].Content, request1.Content)
// 			// 	assert.Equal(t, got[1].Tags, request1.Tags)

// 		}
// 	})
// }

func TestEntRepository_CreateRequest(t *testing.T) {
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

		request1, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, request1.CreatedBy, user.ID)
		assert.Equal(t, request1.Tags, []*Tag{tag})
		request2, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag}, group, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, request2.CreatedBy, user.ID)
		assert.Equal(t, request2.Tags, []*Tag{tag})

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
		tag1, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)
		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget, &[]User{*owner})
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{tag1}, group, user1.ID)
		require.NoError(t, err)

		amount := 100
		_, err = repo.UpdateRequest(ctx, request.ID, amount, request.Title, request.Content, []*Tag{tag1}, group)
		assert.NoError(t, err)
		// assert.Equal(t, updatedRequest.ID, request.ID)
		// assert.Equal(t, updatedRequest.Amount, amount)
		// assert.Equal(t, updatedRequest.Title, request.Title)
		// assert.Equal(t, updatedRequest.Content, request.Content)
		// assert.Equal(t, updatedRequest.Tags, request.Tags)

	})

}
