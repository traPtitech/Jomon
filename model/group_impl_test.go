package model

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetGroups(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_groups")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "get_groups2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget)
		require.NoError(t, err)

		groups, err := repo.GetGroups(ctx)
		require.NoError(t, err)
		assert.Equal(t, 1, len(groups))
		assert.Equal(t, group.ID, groups[0].ID)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		groups, err := repo2.GetGroups(ctx)
		require.NoError(t, err)
		assert.Equal(t, 0, len(groups))
	})
}

func TestEntRepository_GetGroup(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_group")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget)
		require.NoError(t, err)

		g, err := repo.GetGroup(ctx, group.ID)
		assert.NoError(t, err)
		assert.Equal(t, group.ID, g.ID)
		assert.Equal(t, group.Name, g.Name)
		assert.Equal(t, group.Description, g.Description)
		assert.Equal(t, group.Budget, g.Budget)
	})

	t.Run("UnknownGroup", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetGroup(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestEntRepository_CreateGroup(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "create_group")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		name := random.AlphaNumeric(t, 20)
		description := random.AlphaNumeric(t, 15)
		group, err := repo.CreateGroup(ctx, name, description, &budget)
		require.NoError(t, err)
		assert.Equal(t, name, group.Name)
		assert.Equal(t, description, group.Description)
		assert.Equal(t, *group.Budget, budget)
	})

	t.Run("SuccessWithNilBudget", func(t *testing.T) {
		t.Parallel()
		name := random.AlphaNumeric(t, 20)
		description := random.AlphaNumeric(t, 15)
		group, err := repo.CreateGroup(ctx, name, description, nil)
		assert.NoError(t, err)
		assert.Equal(t, name, group.Name)
		assert.Equal(t, description, group.Description)
		assert.Nil(t, group.Budget)
	})

	t.Run("FailedWithEmptyName", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		_, err := repo.CreateGroup(ctx, "", random.AlphaNumeric(t, 15), &budget)
		assert.Error(t, err)
	})
}

func TestEntRepository_UpdateGroup(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "update_group")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget)
		require.NoError(t, err)

		updatedBudget := random.Numeric(t, 10000)
		ug := Group{
			ID:          group.ID,
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 15),
			Budget:      &updatedBudget,
		}
		updated, err := repo.UpdateGroup(ctx, group.ID, ug.Name, ug.Description, ug.Budget)
		assert.NoError(t, err)
		assert.Equal(t, ug.Name, updated.Name)
		assert.Equal(t, ug.Description, updated.Description)
		assert.Equal(t, ug.Budget, updated.Budget)
	})

	t.Run("UnknownGroup", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		_, err := repo.UpdateGroup(ctx, uuid.New(), random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget)
		assert.Error(t, err)
	})
}

func TestEntRepository_DeleteGroup(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "delete_group")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget)
		require.NoError(t, err)

		err = repo.DeleteGroup(ctx, group.ID)
		assert.NoError(t, err)
	})

	t.Run("UnknownGroup", func(t *testing.T) {
		t.Parallel()
		err := repo.DeleteGroup(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestEntRepository_GetMembers(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_members")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "get_members2")
	assert.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget)
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
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget)
		require.NoError(t, err)

		got, err := repo2.GetMembers(ctx, group.ID)
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
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget)
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
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget)
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
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget)
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
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget)
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
			assert.Equal(t, got[1].ID, user1.ID)
			assert.Equal(t, got[0].ID, user2.ID)
		}
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget)
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
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget)
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
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget)
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
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget)
		require.NoError(t, err)

		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		owner, err := repo.CreateOwner(ctx, group.ID, user.ID)
		require.NoError(t, err)

		err = repo.DeleteOwner(ctx, group.ID, owner.ID)
		assert.NoError(t, err)
	})
}
