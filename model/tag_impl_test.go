package model

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetTags(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_tags")
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "get_tags2")
	assert.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)

	t.Run("Success", func(t *testing.T) {
		got, err := repo.GetTags(ctx)
		assert.NoError(t, err)
		assert.Len(t, got, 0)
	})

	t.Run("Success2", func(t *testing.T) {
		tag1, _ := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		tag2, _ := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))

		got, err := repo.GetTags(ctx)
		assert.NoError(t, err)
		if assert.Len(t, got, 2) && got[0].ID == tag1.ID {
			assert.Equal(t, got[0].ID, tag1.ID)
			assert.Equal(t, got[0].Name, tag1.Name)
			assert.Equal(t, got[1].ID, tag2.ID)
			assert.Equal(t, got[1].Name, tag2.Name)
		} else if assert.Len(t, got, 2) {
			assert.Equal(t, got[0].ID, tag2.ID)
			assert.Equal(t, got[0].Name, tag2.Name)
			assert.Equal(t, got[1].ID, tag1.ID)
			assert.Equal(t, got[1].Name, tag1.Name)
		}
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		got, err := repo2.GetTags(ctx)
		assert.NoError(t, err)
		assert.Len(t, got, 0)
	})
}

func TestEntRepository_CreateTag(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "create_tag")
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		name := random.AlphaNumeric(t, 20)
		tag, err := repo.CreateTag(ctx, name)

		assert.NoError(t, err)
		assert.Equal(t, name, tag.Name)
	})

	t.Run("MissingName", func(t *testing.T) {
		t.Parallel()
		name := ""
		_, err := repo.CreateTag(ctx, name)

		assert.Error(t, err)
	})
}

func TestEntRepository_UpdateTag(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "update_tag")
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success1", func(t *testing.T) {
		t.Parallel()
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		assert.NoError(t, err)

		name := random.AlphaNumeric(t, 20)

		updated, err := repo.UpdateTag(ctx, tag.ID, name)

		assert.NoError(t, err)
		assert.Equal(t, tag.ID, updated.ID)
		assert.Equal(t, name, updated.Name)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		assert.NoError(t, err)

		updated, err := repo.UpdateTag(ctx, tag.ID, tag.Name)

		assert.NoError(t, err)
		assert.Equal(t, tag.ID, updated.ID)
		assert.Equal(t, tag.Name, updated.Name)
	})

	t.Run("MissingName", func(t *testing.T) {
		t.Parallel()
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		assert.NoError(t, err)

		name := ""
		_, err = repo.UpdateTag(ctx, tag.ID, name)

		assert.Error(t, err)
	})
}

func TestEntRepository_DeleteTag(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "delete_tag")
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		name := random.AlphaNumeric(t, 20)
		tag, err := repo.CreateTag(ctx, name)
		assert.NoError(t, err)

		err = repo.DeleteTag(ctx, tag.ID)
		assert.NoError(t, err)
	})

	t.Run("UnknownTag", func(t *testing.T) {
		t.Parallel()
		err = repo.DeleteTag(ctx, uuid.New())
		assert.Error(t, err)
	})
}
