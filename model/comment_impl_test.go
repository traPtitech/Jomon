package model

import (
	"context"
	"testing"

	"github.com/google/uuid"
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
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{}, nil, user.ID)
		require.NoError(t, err)

		comment1, err := repo.CreateComment(ctx, random.AlphaNumeric(t, 30), request.ID, user.ID)
		require.NoError(t, err)
		comment2, err := repo.CreateComment(ctx, random.AlphaNumeric(t, 30), request.ID, user.ID)
		require.NoError(t, err)
		
		got, err := repo.GetComments(ctx, request.ID)
		assert.NoError(t, err)
		if assert.Len(t, got, 2) && got[0].ID == comment1.ID {
			assert.Equal(t, got[0].ID, comment1.ID)
			assert.Equal(t, got[0].User, comment1.User)
			assert.Equal(t, got[0].Comment, comment1.Comment)
			assert.Equal(t, got[1].ID, comment2.ID)
			assert.Equal(t, got[1].User, comment2.User)
			assert.Equal(t, got[1].Comment, comment2.Comment)
		}else if assert.Len(t, got, 2) {
			assert.Equal(t, got[0].ID, comment2.ID)
			assert.Equal(t, got[0].User, comment2.User)
			assert.Equal(t, got[0].Comment, comment2.Comment)
			assert.Equal(t, got[1].ID, comment1.ID)
			assert.Equal(t, got[1].User, comment1.User)
			assert.Equal(t, got[1].Comment, comment1.Comment)
		}
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{}, nil, user.ID)
		require.NoError(t, err)

		got, err := repo.GetComments(ctx, request.ID)
		assert.NoError(t, err)
		assert.Len(t, got, 0)
	})

	t.Run("UnknownRequest", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetComments(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestEntRepository_CreateComments(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx)
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{}, nil, user.ID)
		require.NoError(t, err)

		comment := random.AlphaNumeric(t, 30)
		user2, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		created, err := repo.CreateComment(ctx, comment, request.ID, user2.ID)
		assert.NoError(t, err)
		assert.Equal(t, created.User, user2.ID)
		assert.Equal(t, created.Comment, comment)
	})

	t.Run("UnknownRequest", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		_, err = repo.CreateComment(ctx, random.AlphaNumeric(t, 30), uuid.New(), user.ID)
		assert.Error(t, err)
	})

	t.Run("UnknownUser", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{}, nil, user.ID)
		require.NoError(t, err)

		_, err = repo.CreateComment(ctx, random.AlphaNumeric(t, 30), request.ID, uuid.New())
		assert.Error(t, err)
	})
}

func TestEntREpository_UpdateComment(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx)
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{}, nil, user.ID)
		require.NoError(t, err)
		created, err := repo.CreateComment(ctx, random.AlphaNumeric(t, 30), request.ID, user.ID)
		require.NoError(t, err)

    comment := random.AlphaNumeric(t, 30)
		updated, err := repo.UpdateComment(ctx, comment, request.ID, created.ID)
		assert.NoError(t, err)
		assert.Equal(t, updated.ID, created.ID)
		assert.Equal(t, updated.User, created.User)
		assert.Equal(t, updated.Comment, comment)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request1, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{}, nil, user.ID)
		require.NoError(t, err)
		comment, err := repo.CreateComment(ctx, random.AlphaNumeric(t, 30), request1.ID, user.ID)
		require.NoError(t, err)

		request2, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{}, nil, user.ID)
		require.NoError(t, err)
		updated, err := repo.UpdateComment(ctx, comment.Comment, request2.ID, comment.ID)
		assert.NoError(t, err)
		assert.Equal(t, updated.ID, comment.ID)
		assert.Equal(t, updated.User, comment.User)
		assert.Equal(t, updated.Comment, comment.Comment)

		got, err := repo.GetComments(ctx, request2.ID)
		require.NoError(t, err)
		assert.Equal(t, got[0].ID, updated.ID)
		assert.Equal(t, got[0].User, updated.User)
		assert.Equal(t, got[0].Comment, updated.Comment)
	})

	t.Run("UnknownComment", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{}, nil, user.ID)
		require.NoError(t, err)

		_, err = repo.UpdateComment(ctx, random.AlphaNumeric(t, 30), request.ID, uuid.New())
		assert.Error(t, err)
	})

	t.Run("UnknownRequest", func(t *testing.T) {
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{}, nil, user.ID)
		require.NoError(t, err)
		comment, err := repo.CreateComment(ctx, random.AlphaNumeric(t, 30), request.ID, user.ID)
		require.NoError(t, err)

		_, err = repo.UpdateComment(ctx, random.AlphaNumeric(t, 30), uuid.New(), comment.ID)
		assert.Error(t, err)
	})
}

func TestEntRepository_DeleteComment(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx)
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{}, nil, user.ID)
		require.NoError(t, err)
		comment, err := repo.CreateComment(ctx, random.AlphaNumeric(t, 30), request.ID, user.ID)
		require.NoError(t, err)

		err = repo.DeleteComment(ctx, request.ID, comment.ID)
		assert.NoError(t, err)
	})

	t.Run("UnknownRequest", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{}, nil, user.ID)
		require.NoError(t, err)
		comment, err := repo.CreateComment(ctx, random.AlphaNumeric(t, 30), request.ID, user.ID)
		require.NoError(t, err)

		err = repo.DeleteComment(ctx, uuid.New(), comment.ID)
		assert.Error(t, err)
	})

	t.Run("UnknownComment", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 100), []*Tag{}, nil, user.ID)
		require.NoError(t, err)

		err = repo.DeleteComment(ctx, request.ID, uuid.New())
		assert.Error(t, err)
	})
}