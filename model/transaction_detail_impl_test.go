package model

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_createTransactionDetail(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "create_transaction_detail")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	r := repo.(*EntRepository)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 10)

		// Create TransactionDetail
		td, err := r.createTransactionDetail(ctx, amount, target)
		assert.NoError(t, err)
		assert.NotNil(t, td)
		assert.Equal(t, td.Amount, amount)
		assert.Equal(t, td.Target, target)
	})
}

func TestEntRepository_updateTransactionDetail(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "update_transaction_detail")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	r := repo.(*EntRepository)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		amount := 100
		target := "hoge"

		// Create Transaction
		tx, err := repo.CreateTransaction(ctx, amount, target, nil, nil, nil)
		require.NoError(t, err)

		// Update TransactionDetail
		updatedAmount := 1000
		updatedTarget := "fuga"
		td, err := r.updateTransactionDetail(ctx, tx.ID, updatedAmount, updatedTarget)
		assert.NoError(t, err)
		assert.NotNil(t, td)
		assert.Equal(t, td.Amount, updatedAmount)
		assert.Equal(t, td.Target, updatedTarget)
	})
}
