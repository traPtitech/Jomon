package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/internal/testutil"
	"github.com/traPtitech/Jomon/internal/testutil/random"
)

func TestEntRepository_GetUserBySubject(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "get_user_by_subject")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		createdUser, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			false,
		)
		require.NoError(t, err)
		subject := random.AlphaNumeric(t, 32)

		err = repo.RegisterUserSubject(ctx, subject, createdUser.ID)
		require.NoError(t, err)

		got, err := repo.GetUserBySubject(ctx, subject)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		testutil.RequireEqual(t, createdUser, got, opts...)
	})

	t.Run("UnknownSubject", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		_, err := repo.GetUserBySubject(ctx, random.AlphaNumeric(t, 32))
		require.Error(t, err)
	})
}

func TestEntRepository_RegisterUserSubject(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "register_user_subject")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		createdUser, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			false,
		)
		require.NoError(t, err)
		subject := random.AlphaNumeric(t, 32)

		err = repo.RegisterUserSubject(ctx, subject, createdUser.ID)
		require.NoError(t, err)

		got, err := repo.GetUserBySubject(ctx, subject)
		require.NoError(t, err)
		require.Equal(t, createdUser.ID, got.ID)
	})

	t.Run("UnknownUser", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		err := repo.RegisterUserSubject(ctx, random.AlphaNumeric(t, 32), uuid.New())
		require.Error(t, err)
	})

	t.Run("DuplicateSubject", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		user1, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			false,
		)
		require.NoError(t, err)
		user2, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			false,
		)
		require.NoError(t, err)
		subject := random.AlphaNumeric(t, 32)

		err = repo.RegisterUserSubject(ctx, subject, user1.ID)
		require.NoError(t, err)
		err = repo.RegisterUserSubject(ctx, subject, user2.ID)
		require.Error(t, err)
	})
}
