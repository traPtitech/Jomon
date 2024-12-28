package model

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/Jomon/testutil"
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
		t.Parallel()
		tag1, _ := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		tag2, _ := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))

		got, err := repo.GetTags(ctx)
		assert.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.SortSlices(func(l, r *Tag) bool {
				return l.ID.ID() < r.ID.ID()
			}))
		exp := []*Tag{tag1, tag2}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		got, err := repo2.GetTags(ctx)
		assert.NoError(t, err)
		assert.Empty(t, got)
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
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(Tag{}, "ID"))
		exp := &Tag{
			Name:      name,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		}
		testutil.RequireEqual(t, exp, tag, opts...)
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
		created, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		assert.NoError(t, err)

		name := random.AlphaNumeric(t, 20)

		updated, err := repo.UpdateTag(ctx, created.ID, name)

		assert.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &Tag{
			ID:        created.ID,
			Name:      name,
			CreatedAt: created.CreatedAt,
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		}
		testutil.RequireEqual(t, exp, updated, opts...)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		assert.NoError(t, err)

		updated, err := repo.UpdateTag(ctx, tag.ID, tag.Name)

		assert.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		testutil.RequireEqual(t, tag, updated, opts...)
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
