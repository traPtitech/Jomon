package model

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetTransactions(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_transactions")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "get_transactions2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)
	client3, storage3, err := setup(t, ctx, "get_transactions3")
	require.NoError(t, err)
	repo3 := NewEntRepository(client3, storage3)
	client4, storage4, err := setup(t, ctx, "get_transactions4")
	require.NoError(t, err)
	repo4 := NewEntRepository(client4, storage4)
	client5, storage5, err := setup(t, ctx, "get_transactions5")
	require.NoError(t, err)
	repo5 := NewEntRepository(client5, storage5)
	client6, storage6, err := setup(t, ctx, "get_transactions6")
	require.NoError(t, err)
	repo6 := NewEntRepository(client6, storage6)
	client7, storage7, err := setup(t, ctx, "get_transactions7")
	require.NoError(t, err)
	repo7 := NewEntRepository(client7, storage7)
	client8, storage8, err := setup(t, ctx, "get_transactions8")
	require.NoError(t, err)
	repo8 := NewEntRepository(client8, storage8)
	client9, storage9, err := setup(t, ctx, "get_transactions9")
	require.NoError(t, err)
	repo9 := NewEntRepository(client9, storage9)
	client10, storage10, err := setup(t, ctx, "get_transactions10")
	require.NoError(t, err)
	repo10 := NewEntRepository(client10, storage10)
	client11, storage11, err := setup(t, ctx, "get_transactions11")
	require.NoError(t, err)
	repo11 := NewEntRepository(client11, storage11)

	t.Run("SuccessWithSortCreatedAt", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		request, err := repo.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), nil, nil, nil, user.ID)
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

	t.Run("SuccessWithSortCreatedAtDesc", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo2.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		request, err := repo2.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), nil, nil, nil, user.ID)
		require.NoError(t, err)

		tx1, err := repo2.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		tx2, err := repo2.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)

		// Get Transactions
		sort := "-created_at"
		query := TransactionQuery{
			Sort: &sort,
		}
		got, err := repo2.GetTransactions(ctx, query)
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

	t.Run("SuccessWithSortAmount", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo3.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		target := random.AlphaNumeric(t, 20)
		request, err := repo3.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), nil, nil, nil, user.ID)
		require.NoError(t, err)

		tx1, err := repo3.CreateTransaction(ctx, 100, target, nil, nil, &request.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		tx2, err := repo3.CreateTransaction(ctx, 10000, target, nil, nil, &request.ID)
		require.NoError(t, err)

		// Get Transactions
		sort := "amount"
		query := TransactionQuery{
			Sort: &sort,
		}
		got, err := repo3.GetTransactions(ctx, query)
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

	t.Run("SuccessWithSortAmountDesc", func(t *testing.T) {
		err := dropAll(t, ctx, client)
		require.NoError(t, err)
		ctx := context.Background()

		// Create user
		user, err := repo4.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		target := random.AlphaNumeric(t, 20)
		request, err := repo4.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), nil, nil, nil, user.ID)
		require.NoError(t, err)

		tx1, err := repo4.CreateTransaction(ctx, 100, target, nil, nil, &request.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		tx2, err := repo4.CreateTransaction(ctx, 10000, target, nil, nil, &request.ID)
		require.NoError(t, err)

		// Get Transactions
		sort := "-amount"
		query := TransactionQuery{
			Sort: &sort,
		}
		got, err := repo4.GetTransactions(ctx, query)
		assert.NoError(t, err)
		if assert.Len(t, got, 2) {
			assert.Equal(t, tx2.ID, got[0].ID)
			assert.Equal(t, tx2.Amount, got[0].Amount)
			assert.Equal(t, tx2.Target, got[0].Target)
			assert.Equal(t, tx2.CreatedAt, got[0].CreatedAt)
			assert.Equal(t, tx2.UpdatedAt, got[0].UpdatedAt)
			assert.Equal(t, tx1.ID, got[1].ID)
			assert.Equal(t, tx1.Amount, got[1].Amount)
			assert.Equal(t, tx1.Target, got[1].Target)
			assert.Equal(t, tx1.CreatedAt, got[1].CreatedAt)
			assert.Equal(t, tx1.UpdatedAt, got[1].UpdatedAt)
		}
	})

	t.Run("SuccessWithNoneSort", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo5.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		request, err := repo5.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), nil, nil, nil, user.ID)
		require.NoError(t, err)

		tx1, err := repo5.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		tx2, err := repo5.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)

		// Get Transactions
		sort := ""
		query := TransactionQuery{
			Sort: &sort,
		}
		got, err := repo5.GetTransactions(ctx, query)
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

	t.Run("SuccessWithTarget", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo6.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		amount := random.Numeric(t, 100000)
		target1 := random.AlphaNumeric(t, 20)
		target2 := random.AlphaNumeric(t, 20)
		request, err := repo6.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), nil, nil, nil, user.ID)
		require.NoError(t, err)

		tx, err := repo6.CreateTransaction(ctx, amount, target1, nil, nil, &request.ID)
		require.NoError(t, err)
		_, err = repo6.CreateTransaction(ctx, amount, target2, nil, nil, &request.ID)
		require.NoError(t, err)

		// Get Transactions
		query := TransactionQuery{
			Target: &target1,
		}
		got, err := repo6.GetTransactions(ctx, query)
		assert.NoError(t, err)
		if assert.Len(t, got, 1) {
			assert.Equal(t, tx.ID, got[0].ID)
			assert.Equal(t, tx.Amount, got[0].Amount)
			assert.Equal(t, tx.Target, got[0].Target)
			assert.Equal(t, tx.CreatedAt, got[0].CreatedAt)
			assert.Equal(t, tx.UpdatedAt, got[0].UpdatedAt)
		}
	})

	t.Run("SuccessWithSinceUntil", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo7.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		request, err := repo7.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), nil, nil, nil, user.ID)
		require.NoError(t, err)

		_, err = repo7.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)

		// Get Transactions
		since := time.Now()
		until := time.Now().Add(time.Hour * 24)
		query := TransactionQuery{
			Since: &since,
			Until: &until,
		}

		time.Sleep(1 * time.Second)
		tx, err := repo7.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)

		got, err := repo7.GetTransactions(ctx, query)
		assert.NoError(t, err)
		if assert.Len(t, got, 1) {
			assert.Equal(t, tx.ID, got[0].ID)
			assert.Equal(t, tx.Amount, got[0].Amount)
			assert.Equal(t, tx.Target, got[0].Target)
			assert.Equal(t, tx.CreatedAt, got[0].CreatedAt)
			assert.Equal(t, tx.UpdatedAt, got[0].UpdatedAt)
		}
	})

	t.Run("SuccessWithTag", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo8.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		tag, err := repo8.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		request, err := repo8.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), nil, nil, nil, user.ID)
		require.NoError(t, err)

		_, err = repo8.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)

		// Get Transactions
		query := TransactionQuery{
			Tag: &tag.Name,
		}

		tx, err := repo8.CreateTransaction(ctx, amount, target, []*uuid.UUID{&tag.ID}, nil, &request.ID)
		require.NoError(t, err)

		got, err := repo8.GetTransactions(ctx, query)
		assert.NoError(t, err)
		if assert.Len(t, got, 1) {
			assert.Equal(t, tx.ID, got[0].ID)
			assert.Equal(t, tx.Amount, got[0].Amount)
			assert.Equal(t, tx.Target, got[0].Target)
			assert.Equal(t, tx.CreatedAt, got[0].CreatedAt)
			assert.Equal(t, tx.UpdatedAt, got[0].UpdatedAt)
		}
	})

	t.Run("SuccessWithGroup", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo9.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		budget := random.Numeric(t, 100000)
		group, err := repo9.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &budget)
		require.NoError(t, err)
		request, err := repo9.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), nil, nil, nil, user.ID)
		require.NoError(t, err)

		_, err = repo9.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)

		// Get Transactions
		query := TransactionQuery{
			Group: &group.Name,
		}

		tx, err := repo9.CreateTransaction(ctx, amount, target, nil, &group.ID, &request.ID)
		require.NoError(t, err)

		got, err := repo9.GetTransactions(ctx, query)
		assert.NoError(t, err)
		if assert.Len(t, got, 1) {
			assert.Equal(t, tx.ID, got[0].ID)
			assert.Equal(t, tx.Amount, got[0].Amount)
			assert.Equal(t, tx.Target, got[0].Target)
			assert.Equal(t, tx.CreatedAt, got[0].CreatedAt)
			assert.Equal(t, tx.UpdatedAt, got[0].UpdatedAt)
		}
	})

	t.Run("SuccessWithRequest", func(t *testing.T) {
		t.Parallel()

		// Create user
		user, err := repo10.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		request, err := repo10.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), nil, nil, nil, user.ID)
		require.NoError(t, err)

		_, err = repo10.CreateTransaction(ctx, amount, target, nil, nil, nil)
		require.NoError(t, err)

		// Get Transactions
		query := TransactionQuery{
			Request: &request.ID,
		}

		tx, err := repo10.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)

		got, err := repo10.GetTransactions(ctx, query)
		assert.NoError(t, err)
		if assert.Len(t, got, 1) {
			assert.Equal(t, tx.ID, got[0].ID)
			assert.Equal(t, tx.Amount, got[0].Amount)
			assert.Equal(t, tx.Target, got[0].Target)
			assert.Equal(t, tx.CreatedAt, got[0].CreatedAt)
			assert.Equal(t, tx.UpdatedAt, got[0].UpdatedAt)
		}
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		// Get Transactions
		query := TransactionQuery{}
		got, err := repo11.GetTransactions(ctx, query)
		assert.NoError(t, err)
		assert.Len(t, got, 0)
	})
}

func TestEntRepository_GetTransaction(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_transaction")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()

		// Create user
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		_, err = repo.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), nil, nil, nil, user.ID)
		require.NoError(t, err)

		tx, err := repo.CreateTransaction(ctx, amount, target, nil, nil, nil)
		require.NoError(t, err)

		// Get Transaction
		got, err := repo.GetTransaction(ctx, tx.ID)
		assert.NoError(t, err)
		if assert.NotNil(t, got) {
			assert.Equal(t, tx.ID, got.ID)
			assert.Equal(t, tx.Amount, got.Amount)
			assert.Equal(t, tx.Target, got.Target)
			assert.Equal(t, tx.CreatedAt, got.CreatedAt)
			assert.Equal(t, tx.UpdatedAt, got.UpdatedAt)
		}
	})
}

func TestEntRepository_CreateTransaction(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "create_transaction")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		amount := random.Numeric(t, 100000)

		// Create user
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create tag
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)

		// Create group
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &amount)
		require.NoError(t, err)

		// Create Transactions
		target := random.AlphaNumeric(t, 20)
		request, err := repo.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), nil, nil, nil, user.ID)
		require.NoError(t, err)

		tx, err := repo.CreateTransaction(ctx, amount, target, []*uuid.UUID{&tag.ID}, &group.ID, &request.ID)
		assert.NoError(t, err)
		if assert.NotNil(t, tx) {
			assert.Equal(t, amount, tx.Amount)
			assert.Equal(t, target, tx.Target)
			if assert.Len(t, tx.Tags, 1) {
				assert.Equal(t, tag.ID, tx.Tags[0].ID)
				assert.Equal(t, tag.Name, tx.Tags[0].Name)
			}
			if assert.NotNil(t, tx.Group) {
				assert.Equal(t, group.ID, tx.Group.ID)
				assert.Equal(t, group.Name, tx.Group.Name)
				assert.Equal(t, group.Description, tx.Group.Description)
				assert.Equal(t, group.Budget, tx.Group.Budget)
			}
		}
	})

	t.Run("SuccessWithoutTags", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		amount := random.Numeric(t, 100000)

		// Create user
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create group
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &amount)
		require.NoError(t, err)

		// Create Transactions
		target := random.AlphaNumeric(t, 20)
		request, err := repo.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), nil, nil, nil, user.ID)
		require.NoError(t, err)

		tx, err := repo.CreateTransaction(ctx, amount, target, nil, &group.ID, &request.ID)
		assert.NoError(t, err)
		if assert.NotNil(t, tx) {
			assert.Equal(t, amount, tx.Amount)
			assert.Equal(t, target, tx.Target)
			assert.Len(t, tx.Tags, 0)
			if assert.NotNil(t, tx.Group) {
				assert.Equal(t, group.ID, tx.Group.ID)
				assert.Equal(t, group.Name, tx.Group.Name)
				assert.Equal(t, group.Description, tx.Group.Description)
				assert.Equal(t, group.Budget, tx.Group.Budget)
			}
		}
	})

	t.Run("SuccessWithoutGroup", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		amount := random.Numeric(t, 100000)

		// Create user
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create tag
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)

		// Create Transactions
		target := random.AlphaNumeric(t, 20)
		request, err := repo.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), nil, nil, nil, user.ID)
		require.NoError(t, err)

		tx, err := repo.CreateTransaction(ctx, amount, target, []*uuid.UUID{&tag.ID}, nil, &request.ID)
		assert.NoError(t, err)
		if assert.NotNil(t, tx) {
			assert.Equal(t, amount, tx.Amount)
			assert.Equal(t, target, tx.Target)
			if assert.Len(t, tx.Tags, 1) {
				assert.Equal(t, tag.ID, tx.Tags[0].ID)
				assert.Equal(t, tag.Name, tx.Tags[0].Name)
			}
			assert.Nil(t, tx.Group)
		}
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		// Create Transactions
		target := random.AlphaNumeric(t, 20)
		amount := random.Numeric(t, 100000)

		tx, err := repo.CreateTransaction(ctx, amount, target, nil, nil, nil)
		assert.NoError(t, err)
		if assert.NotNil(t, tx) {
			assert.Equal(t, amount, tx.Amount)
			assert.Equal(t, target, tx.Target)
			assert.Len(t, tx.Tags, 0)
			assert.Nil(t, tx.Group)
		}
	})

	t.Run("SuccessWithNegativeAmount", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		// Create Transactions
		target := random.AlphaNumeric(t, 20)
		amount := -1 * random.Numeric(t, 100000)

		tx, err := repo.CreateTransaction(ctx, amount, target, nil, nil, nil)
		assert.NoError(t, err)
		if assert.NotNil(t, tx) {
			assert.Equal(t, amount, tx.Amount)
			assert.Equal(t, target, tx.Target)
			assert.Len(t, tx.Tags, 0)
			assert.Nil(t, tx.Group)
		}
	})
}

func TestEntRepository_UpdateTransaction(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "update_transaction")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		amount := random.Numeric(t, 100000)

		// Create user
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create tag
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)

		// Create group
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &amount)
		require.NoError(t, err)

		// Create Transactions
		target := random.AlphaNumeric(t, 20)
		request, err := repo.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), nil, nil, nil, user.ID)
		require.NoError(t, err)

		tx, err := repo.CreateTransaction(ctx, amount, target, []*uuid.UUID{&tag.ID}, &group.ID, &request.ID)
		require.NoError(t, err)

		// Update Transactions
		amount = random.Numeric(t, 100000)

		// Create tag
		tag, err = repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)

		// Create group
		group, err = repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &amount)
		require.NoError(t, err)

		// Create Transactions
		target = random.AlphaNumeric(t, 20)
		request, err = repo.CreateRequest(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), nil, nil, nil, user.ID)
		require.NoError(t, err)

		tx, err = repo.UpdateTransaction(ctx, tx.ID, amount, target, []*uuid.UUID{&tag.ID}, &group.ID, &request.ID)
		assert.NoError(t, err)
		if assert.NotNil(t, tx) {
			assert.Equal(t, amount, tx.Amount)
			assert.Equal(t, target, tx.Target)
			if assert.Len(t, tx.Tags, 1) {
				assert.Equal(t, tag.ID, tx.Tags[0].ID)
				assert.Equal(t, tag.Name, tx.Tags[0].Name)
			}
			if assert.NotNil(t, tx.Group) {
				assert.Equal(t, group.ID, tx.Group.ID)
				assert.Equal(t, group.Name, tx.Group.Name)
				assert.Equal(t, group.Description, tx.Group.Description)
				assert.Equal(t, group.Budget, tx.Group.Budget)
			}
		}
	})
}
