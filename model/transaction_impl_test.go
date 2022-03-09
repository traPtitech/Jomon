package model

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetTransactions(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx)
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("SuccessWithSortCreatedAt", func(t *testing.T) {
		ctx := context.Background()

		// Create user
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		request, err := repo.CreateRequest(ctx, amount, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), nil, nil, user.ID)
		require.NoError(t, err)

		tx1, err := repo.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		tx2, err := repo.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)

		// Get Transactions
		sort := "created_at"
		query := TransactionQuery{
			Sort: &sort,
		}
		got, err := repo.GetTransactions(ctx, query)
		assert.NoError(t, err)
		if assert.Len(t, got, 2) {
			assert.Equal(t, tx1.ID, got[1].ID)
			assert.Equal(t, tx1.Amount, got[1].Amount)
			assert.Equal(t, tx1.Target, got[1].Target)
			assert.Equal(t, tx1.CreatedAt, got[1].CreatedAt)
			assert.Equal(t, tx1.UpdatedAt, got[1].UpdatedAt)
			assert.Equal(t, tx2.ID, got[0].ID)
			assert.Equal(t, tx2.Amount, got[0].Amount)
			assert.Equal(t, tx2.Target, got[0].Target)
			assert.Equal(t, tx2.CreatedAt, got[0].CreatedAt)
			assert.Equal(t, tx2.UpdatedAt, got[0].UpdatedAt)
		}
	})

	t.Run("SuccessWithSortCreatedAtAsc", func(t *testing.T) {
		ctx := context.Background()
		err := dropAll(t, ctx, client)
		require.NoError(t, err)

		// Create user
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		request, err := repo.CreateRequest(ctx, amount, random.AlphaNumeric(t, 20), nil, nil, user.ID)
		require.NoError(t, err)

		tx1, err := repo.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		tx2, err := repo.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)

		// Get Transactions
		sort := "-created_at"
		query := TransactionQuery{
			Sort: &sort,
		}
		got, err := repo.GetTransactions(ctx, query)
		assert.NoError(t, err)
		if assert.Len(t, got, 2) {
			assert.Equal(t, tx1.ID, got[0].ID)
			assert.Equal(t, tx1.Amount, got[0].Amount)
			assert.Equal(t, tx1.Target, got[0].Target)
			assert.Equal(t, tx1.CreatedAt, got[0].CreatedAt)
			assert.Equal(t, tx1.UpdatedAt, got[0].UpdatedAt)
			assert.Equal(t, tx2.ID, got[1].ID)
			assert.Equal(t, tx2.Amount, got[1].Amount)
			assert.Equal(t, tx2.Target, got[1].Target)
			assert.Equal(t, tx2.CreatedAt, got[1].CreatedAt)
			assert.Equal(t, tx2.UpdatedAt, got[1].UpdatedAt)
		}
	})
}
