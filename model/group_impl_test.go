package model

import (
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/Jomon/testutil/random"
)

//TestEntRepository_GetMembersのto do 1,CreateGroupをする 2,CreateUserをする 3,CreateMemberをする 4,GetMembersをする

func TestEntRepository_GetMembers(t *testing.T) {
	client, storage, err := setup(t)
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		owner, _ := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		budget := rand.Int()
		group, _ := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, &[]User{*owner})
    got, err := repo.GetMembers(ctx, group.ID)
		assert.NoError(t, err)
		assert.Equal(t, got, []*User{})

		user1, _ := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		user2, _ := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)

		member1, _ := repo.CreateMember(ctx, group.ID, user1.ID)
		member2, _ := repo.CreateMember(ctx, group.ID, user2.ID)
		assert.Equal(t, member1.ID, user1.ID)
		assert.Equal(t, member2.ID, user2.ID)

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
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("sucsess", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		owner, _ := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		budget := rand.Int()
		group, _ := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, &[]User{*owner})

		user, _ := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		member, _ := repo.CreateMember(ctx, group.ID, user.ID)
		
		assert.Equal(t, member.ID, user.ID)
	})
}

func TestEntRepository_DeleteMember(t *testing.T) {
	client, storage, err := setup(t)
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("sucsess", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		owner, _ := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
		budget := rand.Int()
		group, _ := repo.CreateGroup(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), &budget, &[]User{*owner})

		user, _ := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 15), true)
    member, _ := repo.CreateMember(ctx, group.ID, user.ID)
		
		err = repo.DeleteMember(ctx, group.ID, member.ID)
		assert.NoError(t, err)
	})
}
