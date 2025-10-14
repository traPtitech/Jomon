package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/internal/testutil"
	"github.com/traPtitech/Jomon/internal/testutil/random"
)

func TestEntRepository_GetApplicationTargets(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "get_application_targets")
	require.NoError(t, err)
	repo := NewEntRepository(client)
	client2, err := setup(t, ctx, "get_application_targets2")
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
		target1 := &ApplicationTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 100000),
		}
		target2 := &ApplicationTarget{
			Target: user2.ID,
			Amount: random.Numeric(t, 100000),
		}

		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, []*ApplicationTarget{target1, target2},
			user1.ID)
		require.NoError(t, err)
		// CreatedAt の差を1秒以内に収めるためにここで time.Now を取る
		exp := []*ApplicationTargetDetail{
			{Target: target1.Target, Amount: target1.Amount, CreatedAt: time.Now()},
			{Target: target2.Target, Amount: target2.Amount, CreatedAt: time.Now()},
		}
		got, err := repo.GetApplicationTargets(ctx, application.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(ApplicationTargetDetail{}, "ID", "PaidAt"),
			cmpopts.SortSlices(func(l, r *ApplicationTargetDetail) bool {
				return l.Target.ID() < r.Target.ID()
			}))
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()

		user, err := repo2.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo2.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, nil,
			user.ID)
		require.NoError(t, err)
		got, err := repo2.GetApplicationTargets(ctx, application.ID)
		require.NoError(t, err)
		require.Empty(t, got)
	})
}

func TestEntRepository_createApplicationTargets(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "create_application_targets")
	require.NoError(t, err)
	repo := NewEntRepository(client)

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
		target1 := &ApplicationTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 100000),
		}
		target2 := &ApplicationTarget{
			Target: user2.ID,
			Amount: random.Numeric(t, 100000),
		}

		got, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, []*ApplicationTarget{target1, target2},
			user1.ID)
		require.NoError(t, err)
		exp := []*ApplicationTargetDetail{
			{Target: target1.Target, Amount: target1.Amount, CreatedAt: time.Now()},
			{Target: target2.Target, Amount: target2.Amount, CreatedAt: time.Now()},
		}
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(ApplicationTargetDetail{}, "ID", "PaidAt"),
			cmpopts.SortSlices(func(l, r *ApplicationTargetDetail) bool {
				return l.Target.ID() < r.Target.ID()
			}))
		testutil.RequireEqual(t, exp, got.Targets, opts...)
	})
}

func TestEntRepository_deleteApplicationTargets(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "delete_application_targets")
	require.NoError(t, err)
	repo := NewEntRepository(client)
	client2, err := setup(t, ctx, "delete_application_targets2")
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
		target1 := &ApplicationTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 100000),
		}
		target2 := &ApplicationTarget{
			Target: user2.ID,
			Amount: random.Numeric(t, 100000),
		}

		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, []*ApplicationTarget{target1, target2},
			user1.ID)
		require.NoError(t, err)
		_, err = repo.UpdateApplication(
			ctx,
			application.ID,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, []*ApplicationTarget{})
		require.NoError(t, err)
		got, err := repo.GetApplicationTargets(ctx, application.ID)
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		user1, err := repo2.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		user2, err := repo2.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		target1 := &ApplicationTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 100000),
		}
		target2 := &ApplicationTarget{
			Target: user2.ID,
			Amount: random.Numeric(t, 100000),
		}

		application, err := repo2.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, nil,
			user1.ID)
		require.NoError(t, err)
		_, err = repo2.UpdateApplication(
			ctx,
			application.ID,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, []*ApplicationTarget{target1, target2})
		require.NoError(t, err)
		got, err := repo2.GetApplicationTargets(ctx, application.ID)
		require.NoError(t, err)
		require.Len(t, got, 2)
	})
}
