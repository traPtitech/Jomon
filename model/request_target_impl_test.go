package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetRequestTargets(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "get_request_targets")
	require.NoError(t, err)
	repo := NewEntRepository(client)
	client2, err := setup(t, ctx, "get_request_targets2")
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
		target1 := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 100000),
		}
		target2 := &RequestTarget{
			Target: user2.ID,
			Amount: random.Numeric(t, 100000),
		}

		request, err := repo.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, []*RequestTarget{target1, target2},
			nil, user1.ID)
		require.NoError(t, err)
		// CreatedAt の差を1秒以内に収めるためにここで time.Now を取る
		exp := []*RequestTargetDetail{
			{Target: target1.Target, Amount: target1.Amount, CreatedAt: time.Now()},
			{Target: target2.Target, Amount: target2.Amount, CreatedAt: time.Now()},
		}
		got, err := repo.GetRequestTargets(ctx, request.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(RequestTargetDetail{}, "ID", "PaidAt"),
			cmpopts.SortSlices(func(l, r *RequestTargetDetail) bool {
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
		request, err := repo2.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)
		got, err := repo2.GetRequestTargets(ctx, request.ID)
		require.NoError(t, err)
		require.Empty(t, got)
	})
}

func TestEntRepository_createRequestTargets(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "create_request_targets")
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
		target1 := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 100000),
		}
		target2 := &RequestTarget{
			Target: user2.ID,
			Amount: random.Numeric(t, 100000),
		}

		got, err := repo.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, []*RequestTarget{target1, target2},
			nil, user1.ID)
		require.NoError(t, err)
		exp := []*RequestTargetDetail{
			{Target: target1.Target, Amount: target1.Amount, CreatedAt: time.Now()},
			{Target: target2.Target, Amount: target2.Amount, CreatedAt: time.Now()},
		}
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(RequestTargetDetail{}, "ID", "PaidAt"),
			cmpopts.SortSlices(func(l, r *RequestTargetDetail) bool {
				return l.Target.ID() < r.Target.ID()
			}))
		testutil.RequireEqual(t, exp, got.Targets, opts...)
	})
}

func TestEntRepository_deleteRequestTargets(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "delete_request_targets")
	require.NoError(t, err)
	repo := NewEntRepository(client)
	client2, err := setup(t, ctx, "delete_request_targets2")
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
		target1 := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 100000),
		}
		target2 := &RequestTarget{
			Target: user2.ID,
			Amount: random.Numeric(t, 100000),
		}

		request, err := repo.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, []*RequestTarget{target1, target2},
			nil, user1.ID)
		require.NoError(t, err)
		_, err = repo.UpdateRequest(
			ctx,
			request.ID,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, []*RequestTarget{},
			nil)
		require.NoError(t, err)
		got, err := repo.GetRequestTargets(ctx, request.ID)
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
		target1 := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 100000),
		}
		target2 := &RequestTarget{
			Target: user2.ID,
			Amount: random.Numeric(t, 100000),
		}

		request, err := repo2.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, nil,
			nil, user1.ID)
		require.NoError(t, err)
		_, err = repo2.UpdateRequest(
			ctx,
			request.ID,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, []*RequestTarget{target1, target2},
			nil)
		require.NoError(t, err)
		got, err := repo2.GetRequestTargets(ctx, request.ID)
		require.NoError(t, err)
		require.Len(t, got, 2)
	})
}
