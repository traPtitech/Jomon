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

func TestEntRepository_CreateStatus(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "create_status")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, nil,
			user.ID)
		require.NoError(t, err)

		status := Status(random.Numeric(t, 5) + 1)
		created, err := repo.CreateStatus(ctx, application.ID, user.ID, status)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(ApplicationStatus{}, "ID"))
		exp := &ApplicationStatus{
			CreatedBy: user.ID,
			Status:    status,
			CreatedAt: time.Now(),
		}
		testutil.RequireEqual(t, exp, created, opts...)
	})

	t.Run("InvalidStatus", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, nil,
			user.ID)
		require.NoError(t, err)

		invalidStatus := Status(6)
		_, err = repo.CreateStatus(ctx, application.ID, user.ID, invalidStatus)
		require.Error(t, err)
	})

	t.Run("UnknownApplicationID", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)

		status := Status(random.Numeric(t, 5) + 1)
		_, err = repo.CreateStatus(ctx, uuid.New(), user.ID, status)
		require.Error(t, err)
	})

	t.Run("UnknownUserID", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 40),
			nil, nil,
			user.ID)
		require.NoError(t, err)

		status := Status(random.Numeric(t, 5) + 1)
		_, err = repo.CreateStatus(ctx, application.ID, uuid.New(), status)
		require.Error(t, err)
	})
}
