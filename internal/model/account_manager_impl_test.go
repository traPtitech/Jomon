package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/internal/testutil"
	"github.com/traPtitech/Jomon/internal/testutil/random"
)

func TestEntRepository_GetAccountManagers(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "get_accountManagers")
	require.NoError(t, err)
	repo := NewEntRepository(client)
	client2, err := setup(t, ctx, "get_accountManagers2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2)

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

		got, err := repo.GetAccountManagers(ctx)
		require.NoError(t, err)
		exp := []*AccountManager{
			{ID: user1.ID},
			{ID: user2.ID},
		}
		require.ElementsMatch(t, exp, got)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()

		got, err := repo2.GetAccountManagers(ctx)
		require.NoError(t, err)
		require.Empty(t, got)
	})
}

func TestEntRepository_AddAccountManagers(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "add_accountManagers")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			false)
		require.NoError(t, err)

		err = repo.AddAccountManagers(ctx, []uuid.UUID{user.ID})
		require.NoError(t, err)

		u, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		require.True(t, u.AccountManager)
	})
}

func TestEntRepository_DeleteAccountManagers(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "delete_accountManagers")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)

		err = repo.DeleteAccountManagers(ctx, []uuid.UUID{user.ID})
		require.NoError(t, err)

		u, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		require.False(t, u.AccountManager)
	})
}
