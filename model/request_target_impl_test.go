package model

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetRequestTargets(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx)
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		target1 := &Target{
			Target: random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 100000),
		}
		target2 := &Target{
			Target: random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 100000),
		}

		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), nil, []*Target{target1, target2}, nil, user.ID)
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

	t.Run("UnknownID", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetRequestTargets(ctx, uuid.New())
		assert.Error(t, err)
	})
}

