package model

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_CreateStatus(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx)
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		amount := random.Numeric(t, 100000)
		title := random.AlphaNumeric(t, 40)
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, amount, title, nil, nil, nil, user.ID)
		require.NoError(t, err)

		status := Status(random.Numeric(t, 5) + 1)
		created, err := repo.CreateStatus(ctx, request.ID, user.ID, status)
		assert.NoError(t, err)
		assert.Equal(t, created.Status, status)
	})

	t.Run("InvalidStatus", func(t *testing.T) {
		t.Parallel()
		amount := random.Numeric(t, 100000)
		title := random.AlphaNumeric(t, 40)
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, amount, title, nil, nil, nil, user.ID)
		require.NoError(t, err)

		invalidStatus := Status(6)
		_, err = repo.CreateStatus(ctx, request.ID, user.ID, invalidStatus)
		assert.Error(t, err)
	})

	t.Run("UnknownRequestID", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)

		status := Status(random.Numeric(t, 5) + 1)
		_, err = repo.CreateStatus(ctx, uuid.New(), user.ID, status)
		assert.Error(t, err)
	})

	t.Run("UnknownUserID", func(t *testing.T) {
		t.Parallel()
		amount := random.Numeric(t, 100000)
		title := random.AlphaNumeric(t, 40)
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, amount, title, nil, nil, nil, user.ID)
		require.NoError(t, err)

		status := Status(random.Numeric(t, 5) + 1)
		_, err = repo.CreateStatus(ctx, request.ID, uuid.New(), status)
		assert.Error(t, err)
	})
}
