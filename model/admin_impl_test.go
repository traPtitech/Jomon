package model

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetAdmins(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_admins")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "get_admins2")
	assert.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user1, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		user2, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)

		got, err := repo.GetAdmins(ctx)
		assert.NoError(t, err)
		if assert.Len(t, got, 2) && got[0].ID == user1.ID {
			assert.Equal(t, got[0].ID, user1.ID)
			assert.Equal(t, got[1].ID, user2.ID)
		} else if assert.Len(t, got, 2) {
			assert.Equal(t, got[0].ID, user2.ID)
			assert.Equal(t, got[1].ID, user1.ID)
		}
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()

		got, err := repo2.GetAdmins(ctx)
		assert.NoError(t, err)
		assert.Len(t, got, 0)
	})
}

func TestEntRepository_AddAdmins(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "add_admins")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			false)
		require.NoError(t, err)

		err = repo.AddAdmins(ctx, []uuid.UUID{user.ID})
		assert.NoError(t, err)

		u, err := repo.GetUserByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.True(t, u.Admin)
	})
}

func TestEntRepository_DeleteAdmins(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "delete_admins")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)

		err = repo.DeleteAdmins(ctx, []uuid.UUID{user.ID})
		assert.NoError(t, err)

		u, err := repo.GetUserByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.False(t, u.Admin)
	})
}
