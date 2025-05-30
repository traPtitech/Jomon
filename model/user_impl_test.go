package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetUsers(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, storage, err := setup(t, ctx, "get_user")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "get_user2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		got, err := repo.GetUsers(ctx)
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		user1, err := repo2.CreateUser(ctx, "user1", "user1", true)
		require.NoError(t, err)
		user2, err := repo2.CreateUser(ctx, "user2", "user2", true)
		require.NoError(t, err)

		got, err := repo2.GetUsers(ctx)
		require.NoError(t, err)
		require.Len(t, got, 2)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.SortSlices(func(l, r *User) bool {
				return l.ID.ID() < r.ID.ID()
			}))
		exp := []*User{user1, user2}
		testutil.RequireEqual(t, exp, got, opts...)
	})
}

func TestEntRepository_CreateUser(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, storage, err := setup(t, ctx, "create_user")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		name := random.AlphaNumeric(t, 20)
		dn := random.AlphaNumeric(t, 20)
		admin := random.Numeric(t, 2) == 1

		user, err := repo.CreateUser(ctx, name, dn, admin)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts, cmpopts.IgnoreFields(User{}, "ID"))
		exp := &User{
			Name:        name,
			DisplayName: dn,
			Admin:       admin,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}
		testutil.RequireEqual(t, exp, user, opts...)
	})

	t.Run("MissingName", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		name := ""
		dn := random.AlphaNumeric(t, 20)
		admin := random.Numeric(t, 2) == 1

		_, err := repo.CreateUser(ctx, name, dn, admin)
		require.Error(t, err)
	})
}

func TestEntRepository_GetUserByName(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, storage, err := setup(t, ctx, "get_user_by_name")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "get_user_by_name2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		name := random.AlphaNumeric(t, 20)
		dn := random.AlphaNumeric(t, 20)
		admin := random.Numeric(t, 2) == 1

		user, err := repo.CreateUser(ctx, name, dn, admin)
		require.NoError(t, err)

		got, err := repo.GetUserByName(ctx, name)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		testutil.RequireEqual(t, user, got, opts...)
	})

	t.Run("UnknownName", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		_, err = repo2.GetUserByName(ctx, random.AlphaNumeric(t, 20))
		require.Error(t, err)
	})
}

func TestEntRepository_GetUserByID(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, storage, err := setup(t, ctx, "get_user_by_id")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "get_user_by_id2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		name := random.AlphaNumeric(t, 20)
		dn := random.AlphaNumeric(t, 20)
		admin := random.Numeric(t, 2) == 1

		user, err := repo.CreateUser(ctx, name, dn, admin)
		require.NoError(t, err)

		got, err := repo.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		testutil.RequireEqual(t, user, got, opts...)
	})

	t.Run("UnknownUserID", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		_, err := repo2.GetUserByID(ctx, uuid.New())
		require.Error(t, err)
	})
}

func TestEntRepository_UpdateUser(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, storage, err := setup(t, ctx, "update_user")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		name := random.AlphaNumeric(t, 20)
		dn := random.AlphaNumeric(t, 20)
		admin := random.Numeric(t, 2) == 1

		user, err := repo.CreateUser(ctx, name, dn, admin)
		require.NoError(t, err)

		uname := random.AlphaNumeric(t, 20)
		udn := random.AlphaNumeric(t, 20)
		uadmin := random.Numeric(t, 2) == 1
		got, err := repo.UpdateUser(ctx, user.ID, uname, udn, uadmin)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &User{
			ID:          user.ID,
			Name:        uname,
			DisplayName: udn,
			Admin:       uadmin,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   time.Now(),
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("MissingName", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		name := random.AlphaNumeric(t, 20)
		dn := random.AlphaNumeric(t, 20)
		admin := random.Numeric(t, 2) == 1

		user, err := repo.CreateUser(ctx, name, dn, admin)
		require.NoError(t, err)

		uname := ""
		udn := random.AlphaNumeric(t, 20)
		uadmin := random.Numeric(t, 2) == 1
		_, err = repo.UpdateUser(ctx, user.ID, uname, udn, uadmin)
		require.Error(t, err)
	})
}
