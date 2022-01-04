package model

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetComments(t *testing.T) {
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
		request, err := repo.CreateRequest(ctx, amount, title, content, []*Tag{}, nil, user.ID)
		require.NoError(t, err)
		comment1 := random.AlphaNumeric(t, 30)
		comment2 := random.AlphaNumeric(t, 30)

		created1, err := repo.CreateComment(ctx, comment1, request.ID, user.ID)
		require.NoError(t, err)
		created2, err := repo.CreateComment(ctx, comment2, request.ID, user.ID)
		require.NoError(t, err)

		got, err := repo.GetComments(ctx, request.ID)
		assert.NoError(t, err)
		if assert.Len(t, got, 2) && got[0].ID == created1.ID {
			assert.Equal(t, got[0].ID, created1.ID)
			assert.Equal(t, got[0].User, created1.User)
			assert.Equal(t, got[0].Comment, created1.Comment)
			assert.Equal(t, got[1].ID, created2.ID)
			assert.Equal(t, got[1].User, created2.User)
			assert.Equal(t, got[1].Comment, created2.Comment)
		} else if assert.Len(t, got, 2) {
			assert.Equal(t, got[0].ID, created2.ID)
			assert.Equal(t, got[0].User, created2.User)
			assert.Equal(t, got[0].Comment, created2.Comment)
			assert.Equal(t, got[1].ID, created1.ID)
			assert.Equal(t, got[1].User, created1.User)
			assert.Equal(t, got[1].Comment, created1.Comment)
		}
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()

		amount := random.Numeric(t, 100000)
		title := random.AlphaNumeric(t, 40)
		content := random.AlphaNumeric(t, 100)
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, amount, title, content, []*Tag{}, nil, user.ID)
		require.NoError(t, err)

		got, err := repo.GetComments(ctx, request.ID)
		assert.NoError(t, err)
		assert.Len(t, got, 0)
	})
}

func TestEntRepository_CreateComments(t *testing.T) {
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
		request, err := repo.CreateRequest(ctx, amount, title, content, []*Tag{}, nil, user.ID)
		require.NoError(t, err)

		comment1 := random.AlphaNumeric(t, 30)
		comment2 := random.AlphaNumeric(t, 30)
		created1, err := repo.CreateComment(ctx, comment1, request.ID, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, created1.User, user.ID)
		assert.Equal(t, created1.Comment, comment1)

		created2, err := repo.CreateComment(ctx, comment2, request.ID, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, created2.User, user.ID)
		assert.Equal(t, created2.Comment, comment2)
	})
}

func TestEntREpository_UpdateComment(t *testing.T) {
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
		request, err := repo.CreateRequest(ctx, amount, title, content, []*Tag{}, nil, user.ID)
		require.NoError(t, err)
		comment := random.AlphaNumeric(t, 30)
		created, err := repo.CreateComment(ctx, comment, request.ID, user.ID)
		require.NoError(t, err)

		updateComment := random.AlphaNumeric(t, 30)
		updated, err := repo.UpdateComment(ctx, updateComment, request.ID, created.ID)
		assert.NoError(t, err)
		assert.Equal(t, updated.ID, created.ID)
		assert.Equal(t, updated.User, created.User)
		assert.Equal(t, updated.Comment, updateComment)
	})
}
