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

func TestEntRepository_createTransactionDetail(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "create_transaction_detail")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
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
		td, err := repo.createTransactionDetail(ctx, tx, title, amount, target)
		require.NoError(t, err)
		err = tx.Commit()
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(TransactionDetail{}, "ID"))
		exp := &TransactionDetail{
			Title:     title,
			Amount:    amount,
			Target:    target,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		testutil.RequireEqual(t, exp, td, opts...)
	})
}

func TestEntRepository_updateTransactionDetail(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "update_transaction_detail")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
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
		trns, err := repo.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, uuid.Nil)
		require.NoError(t, err)

		// Update TransactionDetail
		updateTitle := "Fuga"
		updatedAmount := 1000
		updatedTarget := "fuga"
		td, err := repo.updateTransactionDetail(
			ctx, tx, trns.ID,
			updateTitle, updatedAmount, updatedTarget)
		require.NoError(t, err)
		err = tx.Commit()
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &TransactionDetail{
			ID:        td.ID,
			Title:     updateTitle,
			Amount:    updatedAmount,
			Target:    updatedTarget,
			CreatedAt: td.CreatedAt,
			UpdatedAt: time.Now(),
		}
		testutil.RequireEqual(t, exp, td, opts...)
	})
}
