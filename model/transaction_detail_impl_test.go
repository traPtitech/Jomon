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

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tx, err := client.Tx(ctx)
		require.NoError(t, err)
		defer func() {
			if v := recover(); v != nil {
				_ = tx.Rollback()
				panic(v)
			}
		}()

		title := random.AlphaNumeric(t, 20)
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 10)
		// Create TransactionDetail
		td, err := repo.createTransactionDetail(ctx, tx, title, amount, target)
		assert.NoError(t, err)
		err = tx.Commit()
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

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tx, err := client.Tx(ctx)
		require.NoError(t, err)
		defer func() {
			if v := recover(); v != nil {
				_ = tx.Rollback()
				panic(v)
			}
		}()

		title := "Hoge"
		amount := 100
		target := "hoge"

		// Create Transaction
		trns, err := repo.CreateTransaction(ctx, title, amount, target, nil, nil, nil)
		require.NoError(t, err)

		// Update TransactionDetail
		updateTitle := "Fuga"
		updatedAmount := 1000
		updatedTarget := "fuga"
		td, err := repo.updateTransactionDetail(
			ctx, tx, trns.ID,
			updateTitle, updatedAmount, updatedTarget)
		assert.NoError(t, err)
		err = tx.Commit()
		assert.NoError(t, err)
		assert.NotNil(t, td)
		assert.Equal(t, td.Amount, updatedAmount)
		assert.Equal(t, td.Target, updatedTarget)
	})
}
