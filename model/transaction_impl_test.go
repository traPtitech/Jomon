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
		err := dropAll(t, ctx, client)
		require.NoError(t, err)
		ctx := context.Background()

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

	t.Run("SuccessWithNoneSort", func(t *testing.T) {
		err := dropAll(t, ctx, client)
		require.NoError(t, err)
		ctx := context.Background()

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
		sort := ""
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

	t.Run("SuccessWithTarget", func(t *testing.T) {
		err := dropAll(t, ctx, client)
		require.NoError(t, err)
		ctx := context.Background()

		// Create user
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		amount := random.Numeric(t, 100000)
		target1 := random.AlphaNumeric(t, 20)
		target2 := random.AlphaNumeric(t, 20)
		request, err := repo.CreateRequest(ctx, amount, random.AlphaNumeric(t, 20), nil, nil, user.ID)
		require.NoError(t, err)

		tx, err := repo.CreateTransaction(ctx, amount, target1, nil, nil, &request.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		_, err = repo.CreateTransaction(ctx, amount, target2, nil, nil, &request.ID)
		require.NoError(t, err)

		// Get Transactions
		query := TransactionQuery{
			Target: &target1,
		}
		got, err := repo.GetTransactions(ctx, query)
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
		err := dropAll(t, ctx, client)
		require.NoError(t, err)
		ctx := context.Background()

		// Create user
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		request, err := repo.CreateRequest(ctx, amount, random.AlphaNumeric(t, 20), nil, nil, user.ID)
		require.NoError(t, err)

		_, err = repo.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)

		// Get Transactions
		since := time.Now()
		until := time.Now().Add(time.Hour * 24)
		query := TransactionQuery{
			Since: &since,
			Until: &until,
		}

		time.Sleep(1 * time.Second)
		tx, err := repo.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)

		got, err := repo.GetTransactions(ctx, query)
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
		err := dropAll(t, ctx, client)
		require.NoError(t, err)
		ctx := context.Background()

		// Create user
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), random.Numeric(t, 1) == 0)
		require.NoError(t, err)

		// Create Transactions
		amount := random.Numeric(t, 100000)
		target := random.AlphaNumeric(t, 20)
		request, err := repo.CreateRequest(ctx, amount, random.AlphaNumeric(t, 20), nil, nil, user.ID)
		require.NoError(t, err)

		_, err = repo.CreateTransaction(ctx, amount, target, nil, nil, nil)
		require.NoError(t, err)

		// Get Transactions
		query := TransactionQuery{
			Request: &request.ID,
		}

		time.Sleep(1 * time.Second)
		tx, err := repo.CreateTransaction(ctx, amount, target, nil, nil, &request.ID)
		require.NoError(t, err)

		got, err := repo.GetTransactions(ctx, query)
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
		err := dropAll(t, ctx, client)
		require.NoError(t, err)
		ctx := context.Background()

		// Get Transactions
		query := TransactionQuery{}
		got, err := repo.GetTransactions(ctx, query)
		assert.NoError(t, err)
		assert.Len(t, got, 0)
	})
}

func TestEntRepository_CreateTransaction(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx)
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
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)

		// Create group
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &amount, []*User{user})
		require.NoError(t, err)

		// Create Transactions
		target := random.AlphaNumeric(t, 20)
		request, err := repo.CreateRequest(ctx, amount, random.AlphaNumeric(t, 20), nil, nil, user.ID)
		require.NoError(t, err)

		tx, err := repo.CreateTransaction(ctx, amount, target, []*uuid.UUID{&tag.ID}, &group.ID, &request.ID)
		assert.NoError(t, err)
		if assert.NotNil(t, tx) {
			assert.Equal(t, amount, tx.Amount)
			assert.Equal(t, target, tx.Target)
			if assert.Len(t, tx.Tags, 1) {
				assert.Equal(t, tag.ID, tx.Tags[0].ID)
				assert.Equal(t, tag.Name, tx.Tags[0].Name)
				assert.Equal(t, tag.Description, tx.Tags[0].Description)
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
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), &amount, []*User{user})
		require.NoError(t, err)

		// Create Transactions
		target := random.AlphaNumeric(t, 20)
		request, err := repo.CreateRequest(ctx, amount, random.AlphaNumeric(t, 20), nil, nil, user.ID)
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
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		require.NoError(t, err)

		// Create Transactions
		target := random.AlphaNumeric(t, 20)
		request, err := repo.CreateRequest(ctx, amount, random.AlphaNumeric(t, 20), nil, nil, user.ID)
		require.NoError(t, err)

		tx, err := repo.CreateTransaction(ctx, amount, target, []*uuid.UUID{&tag.ID}, nil, &request.ID)
		assert.NoError(t, err)
		if assert.NotNil(t, tx) {
			assert.Equal(t, amount, tx.Amount)
			assert.Equal(t, target, tx.Target)
			if assert.Len(t, tx.Tags, 1) {
				assert.Equal(t, tag.ID, tx.Tags[0].ID)
				assert.Equal(t, tag.Name, tx.Tags[0].Name)
				assert.Equal(t, tag.Description, tx.Tags[0].Description)
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
