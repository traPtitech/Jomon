package model

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetRequestTargets(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_request_targets")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "get_request_targets2")
	require.NoError(t, err)
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
		got, err := repo.GetRequestTargets(ctx, request.ID)
		assert.NoError(t, err)
		if assert.Len(t, got, 2) && got[0].Target == target1.Target {
			assert.Equal(t, got[0].Target, target1.Target)
			assert.Equal(t, got[0].Amount, target1.Amount)
			assert.Equal(t, got[1].Target, target2.Target)
			assert.Equal(t, got[1].Amount, target2.Amount)
		} else if assert.Len(t, got, 2) {
			assert.Equal(t, got[0].Target, target2.Target)
			assert.Equal(t, got[0].Amount, target2.Amount)
			assert.Equal(t, got[1].Target, target1.Target)
			assert.Equal(t, got[1].Amount, target1.Amount)
		}
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
		assert.NoError(t, err)
		assert.Len(t, got, 0)
	})
}

func TestEntRepository_createRequestTargets(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "create_request_targets")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

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
		assert.NoError(t, err)
		if got.Targets[0].Target == target1.Target {
			assert.Equal(t, got.Targets[0].Target, target1.Target)
			assert.Equal(t, got.Targets[0].Amount, target1.Amount)
			assert.Equal(t, got.Targets[1].Target, target2.Target)
			assert.Equal(t, got.Targets[1].Amount, target2.Amount)
		} else {
			assert.Equal(t, got.Targets[0].Target, target2.Target)
			assert.Equal(t, got.Targets[0].Amount, target2.Amount)
			assert.Equal(t, got.Targets[1].Target, target1.Target)
			assert.Equal(t, got.Targets[1].Amount, target1.Amount)
		}
	})
}

func TestEntRepository_deleteRequestTargets(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "delete_request_targets")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "delete_request_targets2")
	require.NoError(t, err)
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
		assert.NoError(t, err)
		got, err := repo.GetRequestTargets(ctx, request.ID)
		assert.NoError(t, err)
		assert.Len(t, got, 0)
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
		assert.NoError(t, err)
		got, err := repo2.GetRequestTargets(ctx, request.ID)
		assert.NoError(t, err)
		assert.Len(t, got, 2)
	})
}
