package model

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetTags(t *testing.T) {
	client, err := setup(t)
	assert.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tag1, _ := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		tag2, _ := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))

		got, err := repo.GetTags(ctx)
		assert.NoError(t, err)
		if assert.Len(t, got, 2) && got[0].ID == tag1.ID {
			assert.Equal(t, got[0].ID, tag1.ID)
			assert.Equal(t, got[0].Name, tag1.Name)
			assert.Equal(t, got[0].Description, tag1.Description)
			assert.Equal(t, got[1].ID, tag2.ID)
			assert.Equal(t, got[1].Name, tag2.Name)
			assert.Equal(t, got[1].Description, tag2.Description)
		} else if assert.Len(t, got, 2) {
			assert.Equal(t, got[0].ID, tag2.ID)
			assert.Equal(t, got[0].Name, tag2.Name)
			assert.Equal(t, got[0].Description, tag2.Description)
			assert.Equal(t, got[1].ID, tag1.ID)
			assert.Equal(t, got[1].Name, tag1.Name)
			assert.Equal(t, got[1].Description, tag1.Description)
		}
	})
}

func TestEntRepository_CreateTag(t *testing.T) {
	client, err := setup(t)
	assert.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		name := random.AlphaNumeric(t, 20)
		description := random.AlphaNumeric(t, 30)
		tag, err := repo.CreateTag(ctx, name, description)

		assert.NoError(t, err)
		assert.Equal(t, name, tag.Name)
		assert.Equal(t, description, tag.Description)
	})
}

func TestEntRepository_UpdateTag(t *testing.T) {
	client, err := setup(t)
	assert.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success1", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		assert.NoError(t, err)

		name := random.AlphaNumeric(t, 20)

		updated, err := repo.UpdateTag(ctx, tag.ID, name, tag.Description)

		assert.NoError(t, err)
		assert.Equal(t, tag.ID, updated.ID)
		assert.Equal(t, name, updated.Name)
		assert.Equal(t, tag.Description, updated.Description)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		assert.NoError(t, err)

		description := random.AlphaNumeric(t, 30)

		updated, err := repo.UpdateTag(ctx, tag.ID, tag.Name, description)

		assert.NoError(t, err)
		assert.Equal(t, tag.ID, updated.ID)
		assert.Equal(t, tag.Name, updated.Name)
		assert.Equal(t, description, updated.Description)
	})
}

func TestEntRepository_DeleteTag(t *testing.T) {
	client, err := setup(t)
	assert.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		name := random.AlphaNumeric(t, 20)
		description := random.AlphaNumeric(t, 30)
		tag, err := repo.CreateTag(ctx, name, description)
		assert.NoError(t, err)

		err = repo.DeleteTag(ctx, tag.ID)
		assert.NoError(t, err)
	})
}
