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

func TestEntRepository_GetGroups(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, _, err := setup(t, ctx, "get_groups")
	require.NoError(t, err)
	repo := NewEntRepository(client)
	client2, _, err := setup(t, ctx, "get_groups2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		groups, err := repo.GetGroups(ctx)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		testutil.RequireEqual(t, []*Group{group}, groups, opts...)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		groups, err := repo2.GetGroups(ctx)
		require.NoError(t, err)
		require.Empty(t, groups)
	})
}

func TestEntRepository_GetGroup(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, _, err := setup(t, ctx, "get_group")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		created, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		got, err := repo.GetGroup(ctx, created.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		testutil.RequireEqual(t, created, got, opts...)
	})

	t.Run("UnknownGroup", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetGroup(ctx, uuid.New())
		require.Error(t, err)
	})
}

func TestEntRepository_CreateGroup(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, _, err := setup(t, ctx, "create_group")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		name := random.AlphaNumeric(t, 20)
		description := random.AlphaNumeric(t, 15)
		created, err := repo.CreateGroup(ctx, name, description, &budget)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(Group{}, "ID"))
		exp := &Group{
			Name:        name,
			Description: description,
			Budget:      &budget,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}
		testutil.RequireEqual(t, exp, created, opts...)
	})

	t.Run("SuccessWithNilBudget", func(t *testing.T) {
		t.Parallel()
		name := random.AlphaNumeric(t, 20)
		description := random.AlphaNumeric(t, 15)
		created, err := repo.CreateGroup(ctx, name, description, nil)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(Group{}, "ID"))
		exp := &Group{
			Name:        name,
			Description: description,
			Budget:      nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}
		testutil.RequireEqual(t, exp, created, opts...)
	})

	t.Run("FailedWithEmptyName", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		_, err := repo.CreateGroup(ctx, "", random.AlphaNumeric(t, 15), &budget)
		require.Error(t, err)
	})
}

func TestEntRepository_UpdateGroup(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, _, err := setup(t, ctx, "update_group")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		created, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		updatedBudget := random.Numeric(t, 10000)
		ug := Group{
			ID:          created.ID,
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 15),
			Budget:      &updatedBudget,
		}
		updated, err := repo.UpdateGroup(ctx, created.ID, ug.Name, ug.Description, ug.Budget)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &Group{
			ID:          created.ID,
			Name:        ug.Name,
			Description: ug.Description,
			Budget:      &updatedBudget,
			CreatedAt:   created.CreatedAt,
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}
		testutil.RequireEqual(t, exp, updated, opts...)
	})

	t.Run("UnknownGroup", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		_, err := repo.UpdateGroup(
			ctx,
			uuid.New(),
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.Error(t, err)
	})

	t.Run("SuccessWithNilBudget", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		created, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		ug := Group{
			ID:          created.ID,
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 15),
			Budget:      nil,
		}
		updated, err := repo.UpdateGroup(ctx, created.ID, ug.Name, ug.Description, ug.Budget)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &Group{
			ID:          created.ID,
			Name:        ug.Name,
			Description: ug.Description,
			Budget:      nil,
			CreatedAt:   created.CreatedAt,
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}
		testutil.RequireEqual(t, exp, updated, opts...)
	})

	t.Run("FailedWithEmptyName", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		_, err = repo.UpdateGroup(ctx, group.ID, "", random.AlphaNumeric(t, 15), &budget)
		require.Error(t, err)
	})
}

func TestEntRepository_DeleteGroup(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, _, err := setup(t, ctx, "delete_group")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		err = repo.DeleteGroup(ctx, group.ID)
		require.NoError(t, err)

		groups, err := repo.GetGroups(ctx)
		require.NoError(t, err)
		require.Empty(t, groups)
	})

	t.Run("UnknownGroup", func(t *testing.T) {
		t.Parallel()
		err := repo.DeleteGroup(ctx, uuid.New())
		require.Error(t, err)
	})
}

func TestEntRepository_GetMembers(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, _, err := setup(t, ctx, "get_members")
	require.NoError(t, err)
	repo := NewEntRepository(client)
	client2, _, err := setup(t, ctx, "get_members2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		user1, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			true)
		require.NoError(t, err)
		user2, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			true)
		require.NoError(t, err)

		_, err = repo.AddMembers(ctx, group.ID, []uuid.UUID{user1.ID, user2.ID})
		require.NoError(t, err)

		got, err := repo.GetMembers(ctx, group.ID)
		require.NoError(t, err)
		sortOpt := cmpopts.SortSlices(func(a, b *Member) bool {
			return a.ID.ID() < b.ID.ID()
		})
		exp := []*Member{
			{ID: user1.ID},
			{ID: user2.ID},
		}
		testutil.RequireEqual(t, exp, got, sortOpt)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		got, err := repo2.GetMembers(ctx, group.ID)
		require.NoError(t, err)
		require.Empty(t, got)
	})
}

func TestEntRepository_CreateMember(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, _, err := setup(t, ctx, "create_member")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			true)
		require.NoError(t, err)

		got, err := repo.AddMembers(ctx, group.ID, []uuid.UUID{user.ID})
		require.NoError(t, err)
		exp := []*Member{{ID: user.ID}}
		testutil.RequireEqual(t, exp, got)
	})

	t.Run("UnknownUser", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		_, err = repo.AddMembers(ctx, group.ID, []uuid.UUID{uuid.New()})
		require.Error(t, err)
	})
}

func TestEntRepository_DeleteMember(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, _, err := setup(t, ctx, "delete_member")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			true)
		require.NoError(t, err)
		member, err := repo.AddMembers(ctx, group.ID, []uuid.UUID{user.ID})
		require.NoError(t, err)

		err = repo.DeleteMembers(ctx, group.ID, []uuid.UUID{member[0].ID})
		require.NoError(t, err)
	})
}

func TestEntRepository_GetOwners(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, _, err := setup(t, ctx, "get_owners")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		user1, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			true)
		require.NoError(t, err)
		user2, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			true)
		require.NoError(t, err)

		_, err = repo.AddOwners(ctx, group.ID, []uuid.UUID{user1.ID, user2.ID})
		require.NoError(t, err)

		got, err := repo.GetOwners(ctx, group.ID)
		require.NoError(t, err)
		sortOpt := cmpopts.SortSlices(func(a, b *Owner) bool {
			return a.ID.ID() < b.ID.ID()
		})
		exp := []*Owner{
			{ID: user1.ID},
			{ID: user2.ID},
		}
		testutil.RequireEqual(t, exp, got, sortOpt)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		got, err := repo.GetOwners(ctx, group.ID)
		require.NoError(t, err)
		require.Empty(t, got)
	})
}

// FIXME: これAddOwnersでは?
func TestEntRepository_CreateOwner(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, _, err := setup(t, ctx, "create_owner")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			true)
		require.NoError(t, err)

		owner, err := repo.AddOwners(ctx, group.ID, []uuid.UUID{user.ID})
		require.NoError(t, err)
		exp := []*Owner{{ID: user.ID}}
		testutil.RequireEqual(t, exp, owner)
	})

	t.Run("UnknownUser", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		_, err = repo.AddOwners(ctx, group.ID, []uuid.UUID{uuid.New()})
		require.Error(t, err)
	})
}

func TestEntRepository_DeleteOwner(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, _, err := setup(t, ctx, "delete_owner")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			&budget)
		require.NoError(t, err)

		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 15),
			true)
		require.NoError(t, err)
		_, err = repo.AddOwners(ctx, group.ID, []uuid.UUID{user.ID})
		require.NoError(t, err)

		err = repo.DeleteOwners(ctx, group.ID, []uuid.UUID{user.ID})
		require.NoError(t, err)
	})
}
