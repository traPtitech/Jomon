package model

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_CreateFile(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "create_file")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		var tags []*Tag
		var group *Group
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 50), tags, group, user.ID)
		assert.NoError(t, err)

		mimetype := "image/png"

		name := random.AlphaNumeric(t, 20)

		file, err := repo.CreateFile(ctx, name, mimetype, request.ID)
		assert.NoError(t, err)
		assert.Equal(t, name, file.Name)
		assert.Equal(t, mimetype, file.MimeType)
	})

	t.Run("UnknownRequest", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		assert.NoError(t, err)
		request := Request{
			ID: uuid.New(),
		}

		mimetype := "image/png"

		name := random.AlphaNumeric(t, 20)

		_, err = repo.CreateFile(ctx, name, mimetype, request.ID)
		assert.Error(t, err)
	})

	t.Run("MissingName", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		var tags []*Tag
		var group *Group
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 50), tags, group, user.ID)
		assert.NoError(t, err)

		mimetype := "image/png"

		_, err = repo.CreateFile(ctx, "", mimetype, request.ID)
		assert.Error(t, err)
	})
}

func TestEntRepository_GetFile(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "get_file")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		var tags []*Tag
		var group *Group
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 50), tags, group, user.ID)
		assert.NoError(t, err)

		mimetype := "image/png"

		name := random.AlphaNumeric(t, 20)

		file, err := repo.CreateFile(ctx, name, mimetype, request.ID)
		assert.NoError(t, err)
		got, err := repo.GetFile(ctx, file.ID)
		assert.NoError(t, err)
		assert.Equal(t, file.ID, got.ID)
		assert.Equal(t, file.Name, got.Name)
		assert.Equal(t, file.MimeType, got.MimeType)
	})

	t.Run("UnknownFile", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		_, err = repo.GetFile(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestEntRepository_DeleteFile(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx, "delete_file")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		var tags []*Tag
		var group *Group
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30), true)
		require.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 50), tags, group, user.ID)
		assert.NoError(t, err)

		mimetype := "image/png"

		name := random.AlphaNumeric(t, 20)

		file, err := repo.CreateFile(ctx, name, mimetype, request.ID)
		assert.NoError(t, err)

		err = repo.DeleteFile(ctx, file.ID)
		assert.NoError(t, err)

		r, err := repo.GetRequest(ctx, request.ID)
		require.NoError(t, err)
		assert.Len(t, r.Files, 0)
	})

	t.Run("UnknownFile", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		err = repo.DeleteFile(ctx, uuid.New())
		assert.Error(t, err)
	})
}
