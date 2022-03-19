package model

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetMembers(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_members")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, &[]User{*owner})
		require.NoError(t, err)

		user1, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		user2, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)

		_, err = repo.CreateMember(ctx, group.ID, user1.ID)
		require.NoError(t, err)
		_, err = repo.CreateMember(ctx, group.ID, user2.ID)
		require.NoError(t, err)

		got, err := repo.GetMembers(ctx, group.ID)
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
		ctx := context.Background()
		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, &[]User{*owner})
		require.NoError(t, err)

		got, err := repo.GetMembers(ctx, group.ID)
		assert.NoError(t, err)
		assert.Equal(t, got, []*Member{})
	})
}

func TestEntRepository_CreateMember(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "create_member")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, &[]User{*owner})
		require.NoError(t, err)

		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)

		member, err := repo.CreateMember(ctx, group.ID, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, member.ID, user.ID)
	})

	t.Run("UnknownGroup", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		_, err = repo.CreateMember(ctx, uuid.New(), user.ID)
		assert.Error(t, err)
	})

	t.Run("UnknownUser", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, &[]User{*owner})
		require.NoError(t, err)

		_, err = repo.CreateMember(ctx, group.ID, uuid.New())
		assert.Error(t, err)
	})
}

func TestEntRepository_DeleteMember(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "delete_member")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, &[]User{*owner})
		require.NoError(t, err)

		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		member, err := repo.CreateMember(ctx, group.ID, user.ID)
		require.NoError(t, err)

		err = repo.DeleteMember(ctx, group.ID, member.ID)
		assert.NoError(t, err)
	})

	t.Run("UnknownGroup", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)

		err = repo.DeleteMember(ctx, uuid.New(), user.ID)
		assert.Error(t, err)
	})
}

func TestEntRepository_GetOwners(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_owners")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, &[]User{*owner})
		require.NoError(t, err)

		user1, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		user2, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)

		_, err = repo.CreateOwner(ctx, group.ID, user1.ID)
		require.NoError(t, err)
		_, err = repo.CreateOwner(ctx, group.ID, user2.ID)
		require.NoError(t, err)

		got, err := repo.GetOwners(ctx, group.ID)
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
		ctx := context.Background()
		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, &[]User{*owner})
		require.NoError(t, err)

		got, err := repo.GetOwners(ctx, group.ID)
		assert.NoError(t, err)
		assert.Equal(t, got, []*Owner{})
	})
}

func TestEntRepository_CreateOwner(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "create_owner")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, nil)
		require.NoError(t, err)

		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)

		owner, err := repo.CreateOwner(ctx, group.ID, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, owner.ID, user.ID)
	})

	t.Run("UnknownUser", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, &[]User{*owner})
		require.NoError(t, err)

		_, err = repo.CreateOwner(ctx, group.ID, uuid.New())
		assert.Error(t, err)
	})
}

func TestEntRepository_DeleteOwner(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "delete_owner")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, nil)
		require.NoError(t, err)

		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		owner, err := repo.CreateOwner(ctx, group.ID, user.ID)
		require.NoError(t, err)

		err = repo.DeleteOwner(ctx, group.ID, owner.ID)
		assert.NoError(t, err)
	})
}
