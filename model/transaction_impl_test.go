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

func TestEntRepository_GetTransactions(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "get_transactions")
	require.NoError(t, err)
	repo := NewEntRepository(client)
	client2, err := setup(t, ctx, "get_transactions2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2)
	client3, err := setup(t, ctx, "get_transactions3")
	require.NoError(t, err)
	repo3 := NewEntRepository(client3)
	client4, err := setup(t, ctx, "get_transactions4")
	require.NoError(t, err)
	repo4 := NewEntRepository(client4)
	client5, err := setup(t, ctx, "get_transactions5")
	require.NoError(t, err)
	repo5 := NewEntRepository(client5)
	client6, err := setup(t, ctx, "get_transactions6")
	require.NoError(t, err)
	repo6 := NewEntRepository(client6)
	client7, err := setup(t, ctx, "get_transactions7")
	require.NoError(t, err)
	repo7 := NewEntRepository(client7)
	client8, err := setup(t, ctx, "get_transactions8")
	require.NoError(t, err)
	repo8 := NewEntRepository(client8)
	client9, err := setup(t, ctx, "get_transactions9")
	require.NoError(t, err)
	repo9 := NewEntRepository(client9)
	client10, err := setup(t, ctx, "get_transactions10")
	require.NoError(t, err)
	repo10 := NewEntRepository(client10)
	client11, err := setup(t, ctx, "get_transactions11")
	require.NoError(t, err)
	repo11 := NewEntRepository(client11)

	t.Run("SuccessWithSortCreatedAt", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		title := random.AlphaNumeric(t, 20)
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		tx1, err := repo.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, application.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		tx2, err := repo.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, application.ID)
		require.NoError(t, err)

		// Get Transactions
		sort := "created_at"
		query := TransactionQuery{
			Sort: &sort,
		}
		got, err := repo.GetTransactions(ctx, query)
		require.NoError(t, err)
		require.Len(t, got, 2)
		opts := testutil.ApproxEqualOptions()
		exp := []*TransactionResponse{tx2, tx1}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithSortCreatedAtDesc", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo2.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		title := random.AlphaNumeric(t, 20)
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		application, err := repo2.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		tx1, err := repo2.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, application.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		tx2, err := repo2.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, application.ID)
		require.NoError(t, err)

		// Get Transactions
		sort := "-created_at"
		query := TransactionQuery{
			Sort: &sort,
		}
		got, err := repo2.GetTransactions(ctx, query)
		require.NoError(t, err)
		require.Len(t, got, 2)
		opts := testutil.ApproxEqualOptions()
		exp := []*TransactionResponse{tx1, tx2}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithSortAmount", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo3.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		title := random.AlphaNumeric(t, 20)
		target := random.AlphaNumeric(t, 20)
		application, err := repo3.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		tx1, err := repo3.CreateTransaction(ctx, title, 100, target, nil, uuid.Nil, application.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		tx2, err := repo3.CreateTransaction(ctx, title, 10000, target, nil, uuid.Nil, application.ID)
		require.NoError(t, err)

		// Get Transactions
		sort := "amount"
		query := TransactionQuery{
			Sort: &sort,
		}
		got, err := repo3.GetTransactions(ctx, query)
		require.NoError(t, err)
		require.Len(t, got, 2)
		opts := testutil.ApproxEqualOptions()
		exp := []*TransactionResponse{tx1, tx2}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithSortAmountDesc", func(t *testing.T) {
		err := dropAll(t, ctx, client)
		require.NoError(t, err)
		ctx := testutil.NewContext(t)

		// Create user
		// nolint:contextcheck
		user, err := repo4.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		title := random.AlphaNumeric(t, 20)
		target := random.AlphaNumeric(t, 20)
		// nolint:contextcheck
		application, err := repo4.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		// nolint:contextcheck
		tx1, err := repo4.CreateTransaction(ctx, title, 100, target, nil, uuid.Nil, application.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		// nolint:contextcheck
		tx2, err := repo4.CreateTransaction(ctx, title, 10000, target, nil, uuid.Nil, application.ID)
		require.NoError(t, err)

		// Get Transactions
		sort := "-amount"
		query := TransactionQuery{
			Sort: &sort,
		}
		// nolint:contextcheck
		got, err := repo4.GetTransactions(ctx, query)
		require.NoError(t, err)
		require.Len(t, got, 2)
		opts := testutil.ApproxEqualOptions()
		exp := []*TransactionResponse{tx2, tx1}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithNoneSort", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo5.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		title := random.AlphaNumeric(t, 20)
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		application, err := repo5.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		tx1, err := repo5.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, application.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		tx2, err := repo5.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, application.ID)
		require.NoError(t, err)

		// Get Transactions
		sort := ""
		query := TransactionQuery{
			Sort: &sort,
		}
		got, err := repo5.GetTransactions(ctx, query)
		require.NoError(t, err)
		require.Len(t, got, 2)
		opts := testutil.ApproxEqualOptions()
		exp := []*TransactionResponse{tx2, tx1}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithTarget", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo6.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		title := random.AlphaNumeric(t, 20)
		amount := random.Numeric(t, 100000)
		target1 := random.AlphaNumeric(t, 20)
		target2 := random.AlphaNumeric(t, 20)
		application, err := repo6.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		tx, err := repo6.CreateTransaction(ctx, title, amount, target1, nil, uuid.Nil, application.ID)
		require.NoError(t, err)
		_, err = repo6.CreateTransaction(ctx, title, amount, target2, nil, uuid.Nil, application.ID)
		require.NoError(t, err)

		// Get Transactions
		query := TransactionQuery{
			Target: &target1,
		}
		got, err := repo6.GetTransactions(ctx, query)
		require.NoError(t, err)
		require.Len(t, got, 1)
		opts := testutil.ApproxEqualOptions()
		exp := []*TransactionResponse{tx}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithSinceUntil", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo7.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		title := random.AlphaNumeric(t, 20)
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		application, err := repo7.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		_, err = repo7.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, application.ID)
		require.NoError(t, err)

		// Get Transactions
		since := time.Now()
		until := time.Now().Add(time.Hour * 24)
		query := TransactionQuery{
			Since: since,
			Until: until,
		}

		time.Sleep(1 * time.Second)
		tx, err := repo7.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, application.ID)
		require.NoError(t, err)

		got, err := repo7.GetTransactions(ctx, query)
		require.NoError(t, err)
		require.Len(t, got, 1)
		opts := testutil.ApproxEqualOptions()
		exp := []*TransactionResponse{tx}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithTag", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo8.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		title := random.AlphaNumeric(t, 20)
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		tag, err := repo8.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		application, err := repo8.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		_, err = repo8.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, application.ID)
		require.NoError(t, err)

		// Get Transactions
		query := TransactionQuery{
			Tag: &tag.Name,
		}

		tx, err := repo8.CreateTransaction(
			ctx,
			title, amount, target,
			[]uuid.UUID{tag.ID}, uuid.Nil, application.ID)
		require.NoError(t, err)

		got, err := repo8.GetTransactions(ctx, query)
		require.NoError(t, err)
		require.Len(t, got, 1)
		opts := testutil.ApproxEqualOptions()
		exp := []*TransactionResponse{tx}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithGroup", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo9.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		title := random.AlphaNumeric(t, 20)
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		budget := random.Numeric(t, 100000)
		group, err := repo9.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)
		application, err := repo9.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		_, err = repo9.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, application.ID)
		require.NoError(t, err)

		// Get Transactions
		query := TransactionQuery{
			Group: &group.Name,
		}

		tx, err := repo9.CreateTransaction(ctx, title, amount, target, nil, group.ID, application.ID)
		require.NoError(t, err)

		got, err := repo9.GetTransactions(ctx, query)
		require.NoError(t, err)
		require.Len(t, got, 1)
		opts := testutil.ApproxEqualOptions()
		exp := []*TransactionResponse{tx}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithApplication", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo10.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		title := random.AlphaNumeric(t, 20)
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		application, err := repo10.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		_, err = repo10.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, uuid.Nil)
		require.NoError(t, err)

		// Get Transactions
		query := TransactionQuery{
			Application: application.ID,
		}

		tx, err := repo10.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, application.ID)
		require.NoError(t, err)

		got, err := repo10.GetTransactions(ctx, query)
		require.NoError(t, err)
		require.Len(t, got, 1)
		opts := testutil.ApproxEqualOptions()
		exp := []*TransactionResponse{tx}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		// Get Transactions
		query := TransactionQuery{}
		got, err := repo11.GetTransactions(ctx, query)
		require.NoError(t, err)
		require.Empty(t, got)
	})
}

func TestEntRepository_GetTransaction(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "get_transaction")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		ctx := testutil.NewContext(t)

		// Create user
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		title := random.AlphaNumeric(t, 20)
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		_, err = repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		tx, err := repo.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, uuid.Nil)
		require.NoError(t, err)

		// Get Transaction
		got, err := repo.GetTransaction(ctx, tx.ID)
		require.NoError(t, err)
		require.NotNil(t, got)
		opts := testutil.ApproxEqualOptions()
		testutil.RequireEqual(t, tx, got, opts...)
	})
}

func TestEntRepository_CreateTransaction(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "create_transaction")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		amount := random.Numeric(t, 100000)

		// Create user
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create tag
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)

		// Create group
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&amount)
		require.NoError(t, err)

		// Create Transactions
		title := random.AlphaNumeric(t, 20)
		target := random.AlphaNumeric(t, 20)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		tx, err := repo.CreateTransaction(
			ctx,
			title, amount, target,
			[]uuid.UUID{tag.ID}, group.ID, application.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		// FIXME: #831
		opts = append(opts,
			cmpopts.IgnoreFields(TransactionResponse{}, "ID", "UpdatedAt"))
		exp := &TransactionResponse{
			Title:       title,
			Amount:      amount,
			Target:      target,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Application: application.ID,
			Tags:        []*Tag{tag},
			Group:       group,
		}
		testutil.RequireEqual(t, exp, tx, opts...)
	})

	t.Run("SuccessWithoutTags", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		title := random.AlphaNumeric(t, 20)
		amount := random.Numeric(t, 100000)

		// Create user
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create group
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&amount)
		require.NoError(t, err)

		// Create Transactions
		target := random.AlphaNumeric(t, 20)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		tx, err := repo.CreateTransaction(ctx, title, amount, target, nil, group.ID, application.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		// FIXME: #831
		opts = append(opts,
			cmpopts.IgnoreFields(TransactionResponse{}, "ID", "UpdatedAt"))
		exp := &TransactionResponse{
			Title:       title,
			Amount:      amount,
			Target:      target,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Application: application.ID,
			Tags:        []*Tag{},
			Group:       group,
		}
		testutil.RequireEqual(t, exp, tx, opts...)
	})

	t.Run("SuccessWithoutGroup", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		title := random.AlphaNumeric(t, 20)
		amount := random.Numeric(t, 100000)

		// Create user
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create tag
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)

		// Create Transactions
		target := random.AlphaNumeric(t, 20)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		tx, err := repo.CreateTransaction(
			ctx,
			title, amount, target,
			[]uuid.UUID{tag.ID}, uuid.Nil, application.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		// FIXME: #831
		opts = append(opts,
			cmpopts.IgnoreFields(TransactionResponse{}, "ID", "UpdatedAt"))
		exp := &TransactionResponse{
			Title:       title,
			Amount:      amount,
			Target:      target,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Application: application.ID,
			Tags:        []*Tag{tag},
			Group:       nil,
		}
		testutil.RequireEqual(t, exp, tx, opts...)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		// Create Transactions
		title := random.AlphaNumeric(t, 20)
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)

		tx, err := repo.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, uuid.Nil)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		// FIXME: #831
		opts = append(opts,
			cmpopts.IgnoreFields(TransactionResponse{}, "ID", "UpdatedAt"))
		exp := &TransactionResponse{
			Title:       title,
			Amount:      amount,
			Target:      target,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Application: uuid.Nil,
			Tags:        []*Tag{},
			Group:       nil,
		}
		testutil.RequireEqual(t, exp, tx, opts...)
	})

	t.Run("SuccessWithNegativeAmount", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		// Create Transactions
		title := random.AlphaNumeric(t, 20)
		amount := -1 * random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)

		tx, err := repo.CreateTransaction(ctx, title, amount, target, nil, uuid.Nil, uuid.Nil)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		// FIXME: #831
		opts = append(opts,
			cmpopts.IgnoreFields(TransactionResponse{}, "ID", "UpdatedAt"))
		exp := &TransactionResponse{
			Title:     title,
			Amount:    amount,
			Target:    target,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Tags:      []*Tag{},
			Group:     nil,
		}
		testutil.RequireEqual(t, exp, tx, opts...)
	})
}

func TestEntRepository_UpdateTransaction(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "update_transaction")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		title := random.AlphaNumeric(t, 20)
		amount := random.Numeric(t, 100000)

		// Create user
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 20),
			random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create tag
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)

		// Create group
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&amount)
		require.NoError(t, err)

		// Create Transactions
		target := random.AlphaNumeric(t, 20)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		tx, err := repo.CreateTransaction(
			ctx,
			title, amount, target,
			[]uuid.UUID{tag.ID}, group.ID, application.ID)
		require.NoError(t, err)

		// Update Transactions
		title = random.AlphaNumeric(t, 20)
		amount = random.Numeric(t, 100000)

		// Create tag
		tag, err = repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)

		// Create group
		group, err = repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&amount)
		require.NoError(t, err)

		// Create Transactions
		target = random.AlphaNumeric(t, 20)
		application, err = repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			nil, nil,
			nil, user.ID)
		require.NoError(t, err)

		updated, err := repo.UpdateTransaction(
			ctx,
			tx.ID, title, amount, target,
			[]uuid.UUID{tag.ID}, group.ID, application.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		// FIXME: #831
		opts = append(opts,
			cmpopts.IgnoreFields(TransactionResponse{}, "UpdatedAt"))
		exp := &TransactionResponse{
			ID:          tx.ID,
			Title:       title,
			Amount:      amount,
			Target:      target,
			CreatedAt:   tx.CreatedAt,
			UpdatedAt:   time.Now(),
			Application: application.ID,
			Tags:        []*Tag{tag},
			Group:       group,
		}
		testutil.RequireEqual(t, exp, updated, opts...)
	})
}
