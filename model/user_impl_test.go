package model

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetUsers(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_user")
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "get_user2")
	assert.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		got, err := repo.GetUsers(ctx)
		assert.NoError(t, err)
		assert.Len(t, got, 0)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		user1, err := repo2.CreateUser(ctx, "user1", "user1", true)
		assert.NoError(t, err)
		user2, err := repo2.CreateUser(ctx, "user2", "user2", true)
		assert.NoError(t, err)

		got, err := repo2.GetUsers(ctx)
		assert.NoError(t, err)
		if assert.Len(t, got, 2) && got[0].ID == user1.ID {
			assert.Equal(t, got[0].ID, user1.ID)
			assert.Equal(t, got[0].Name, user1.Name)
			assert.Equal(t, got[0].DisplayName, user1.DisplayName)
			assert.Equal(t, got[0].Admin, user1.Admin)
			assert.Equal(t, got[1].ID, user2.ID)
			assert.Equal(t, got[1].Name, user2.Name)
			assert.Equal(t, got[1].DisplayName, user2.DisplayName)
			assert.Equal(t, got[1].Admin, user2.Admin)
		} else if assert.Len(t, got, 2) {
			assert.Equal(t, got[0].ID, user2.ID)
			assert.Equal(t, got[0].Name, user2.Name)
			assert.Equal(t, got[0].DisplayName, user2.DisplayName)
			assert.Equal(t, got[0].Admin, user2.Admin)
			assert.Equal(t, got[1].ID, user1.ID)
			assert.Equal(t, got[1].Name, user1.Name)
			assert.Equal(t, got[1].DisplayName, user1.DisplayName)
			assert.Equal(t, got[1].Admin, user1.Admin)
		}
	})
}

func TestEntRepository_CreateUser(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "create_user")
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		name := random.AlphaNumeric(t, 20)
		dn := random.AlphaNumeric(t, 20)
		admin := random.Numeric(t, 2) == 1

		user, err := repo.CreateUser(ctx, name, dn, admin)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, name)
		assert.Equal(t, user.DisplayName, dn)
		assert.Equal(t, user.Admin, admin)
	})

	t.Run("MissingName", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		name := ""
		dn := random.AlphaNumeric(t, 20)
		admin := random.Numeric(t, 2) == 1

		_, err := repo.CreateUser(ctx, name, dn, admin)
		assert.Error(t, err)
	})
}

func TestEntRepository_GetUserByName(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_user_by_name")
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "get_user_by_name2")
	assert.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		name := random.AlphaNumeric(t, 20)
		dn := random.AlphaNumeric(t, 20)
		admin := random.Numeric(t, 2) == 1

		user, err := repo.CreateUser(ctx, name, dn, admin)
		assert.NoError(t, err)

		got, err := repo.GetUserByName(ctx, name)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, got.ID)
		assert.Equal(t, user.Name, got.Name)
		assert.Equal(t, user.DisplayName, got.DisplayName)
		assert.Equal(t, user.Admin, got.Admin)
	})

	t.Run("UnknownName", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		_, err = repo2.GetUserByName(ctx, random.AlphaNumeric(t, 20))
		assert.Error(t, err)
	})
}

func TestEntRepository_GetUserByID(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_user_by_id")
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "get_user_by_id2")
	assert.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		name := random.AlphaNumeric(t, 20)
		dn := random.AlphaNumeric(t, 20)
		admin := random.Numeric(t, 2) == 1

		user, err := repo.CreateUser(ctx, name, dn, admin)
		assert.NoError(t, err)

		got, err := repo.GetUserByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, got.ID)
		assert.Equal(t, user.Name, got.Name)
		assert.Equal(t, user.DisplayName, got.DisplayName)
		assert.Equal(t, user.Admin, got.Admin)
	})

	t.Run("UnknownUserID", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		_, err := repo2.GetUserByID(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestEntRepository_UpdateUser(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "update_user")
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		name := random.AlphaNumeric(t, 20)
		dn := random.AlphaNumeric(t, 20)
		admin := random.Numeric(t, 2) == 1

		user, err := repo.CreateUser(ctx, name, dn, admin)
		assert.NoError(t, err)

		uname := random.AlphaNumeric(t, 20)
		udn := random.AlphaNumeric(t, 20)
		uadmin := random.Numeric(t, 2) == 1
		got, err := repo.UpdateUser(ctx, user.ID, uname, udn, uadmin)
		assert.NoError(t, err)

		assert.NoError(t, err)
		assert.Equal(t, got.Name, uname)
		assert.Equal(t, got.DisplayName, udn)
		assert.Equal(t, got.Admin, uadmin)
	})

	t.Run("MissingName", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		name := random.AlphaNumeric(t, 20)
		dn := random.AlphaNumeric(t, 20)
		admin := random.Numeric(t, 2) == 1

		user, err := repo.CreateUser(ctx, name, dn, admin)
		assert.NoError(t, err)

		uname := ""
		udn := random.AlphaNumeric(t, 20)
		uadmin := random.Numeric(t, 2) == 1
		_, err = repo.UpdateUser(ctx, user.ID, uname, udn, uadmin)
		assert.Error(t, err)
	})
}
