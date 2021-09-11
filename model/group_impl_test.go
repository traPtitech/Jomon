package model

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetMembers(t *testing.T) {
	client, storage, err := setup(t)
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, &[]User{*owner})
		require.NoError(t, err)

		got, err := repo.GetMembers(ctx, group.ID)
		assert.NoError(t, err)
		assert.Equal(t, got, []*User{})

		user1, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		user2, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)

		member1, err := repo.CreateMember(ctx, group.ID, user1.ID)
		require.NoError(t, err)
		member2, err := repo.CreateMember(ctx, group.ID, user2.ID)
		require.NoError(t, err)
		require.Equal(t, member1.ID, user1.ID)
		require.Equal(t, member2.ID, user2.ID)

		got, err = repo.GetMembers(ctx, group.ID)
		assert.NoError(t, err)
		if assert.Len(t, got, 2) && got[0].ID == user1.ID {
			assert.Equal(t, got[0].ID, user1.ID)
			assert.Equal(t, got[0].Name, user1.Name)
			assert.Equal(t, got[0].DisplayName, user1.DisplayName)
			assert.Equal(t, got[1].ID, user2.ID)
			assert.Equal(t, got[1].Name, user2.Name)
			assert.Equal(t, got[1].DisplayName, user2.DisplayName)
		} else if assert.Len(t, got, 2) {
			assert.Equal(t, got[0].ID, user2.ID)
			assert.Equal(t, got[0].Name, user2.Name)
			assert.Equal(t, got[0].DisplayName, user2.DisplayName)
			assert.Equal(t, got[1].ID, user1.ID)
			assert.Equal(t, got[1].Name, user1.Name)
			assert.Equal(t, got[1].DisplayName, user1.DisplayName)
		}
	})
}

func TestEntRepository_CreateMember(t *testing.T) {
	client, storage, err := setup(t)
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Sucsess", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, &[]User{*owner})
		require.NoError(t, err)

		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)

		member, err := repo.CreateMember(ctx, group.ID, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, member.ID, user.ID)

		FalseMember, err := repo.CreateMember(ctx, uuid.New(), user.ID)
		if FalseMember == nil {
			assert.NoError(t, err)
		}

	})
}

func TestEntRepository_DeleteMember(t *testing.T) {
	client, storage, err := setup(t)
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Sucsess", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		owner, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		budget := random.Numeric(t, 100000)
		group, err := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, &[]User{*owner})
		require.NoError(t, err)

		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		require.NoError(t, err)
		member, err := repo.CreateMember(ctx, group.ID, user.ID)
		require.NoError(t, err)
		require.Equal(t, member.ID, user.ID)

		err = repo.DeleteMember(ctx, group.ID, member.ID)
		assert.NoError(t, err)
	})
}
