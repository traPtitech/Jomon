package model

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_CreateFile(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx)
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		var tags []*Tag
		var group *Group
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), false)
		assert.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 50), tags, group, user.ID)
		assert.NoError(t, err)

		sampleText := "sampleData"

		mimetype := "image/png"

		src := strings.NewReader(sampleText)

		name := random.AlphaNumeric(t, 20)

		file, err := repo.CreateFile(ctx, src, name, mimetype, request.ID)
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

		sampleText := "sampleData"

		mimetype := "image/png"

		src := strings.NewReader(sampleText)

		name := random.AlphaNumeric(t, 20)

		_, err = repo.CreateFile(ctx, src, name, mimetype, request.ID)
		assert.Error(t, err)
	})

	t.Run("MissingName", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		var tags []*Tag
		var group *Group
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), false)
		assert.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 50), tags, group, user.ID)
		assert.NoError(t, err)

		sampleText := "sampleData"

		mimetype := "image/png"

		src := strings.NewReader(sampleText)

		_, err = repo.CreateFile(ctx, src, "", mimetype, request.ID)
		assert.Error(t, err)
	})
}

func TestEntRepository_GetFile(t *testing.T) {
	ctx := context.Background()
	client, storage, err := setup(t, ctx)
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		var tags []*Tag
		var group *Group
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), false)
		assert.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 50), tags, group, user.ID)
		assert.NoError(t, err)

		sampleText := "sampleData"

		mimetype := "image/png"

		src := strings.NewReader(sampleText)

		name := random.AlphaNumeric(t, 20)

		file, err := repo.CreateFile(ctx, src, name, mimetype, request.ID)
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
	client, storage, err := setup(t, ctx)
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		var tags []*Tag
		var group *Group
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), false)
		assert.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 50), tags, group, user.ID)
		assert.NoError(t, err)

		sampleText := "sampleData"

		mimetype := "image/png"

		src := strings.NewReader(sampleText)

		name := random.AlphaNumeric(t, 20)

		file, err := repo.CreateFile(ctx, src, name, mimetype, request.ID)
		assert.NoError(t, err)

		ext, err := mime.ExtensionsByType(file.MimeType)
		assert.NoError(t, err)

		storage.EXPECT().
			Delete(fmt.Sprintf("%s%s", file.ID.String(), ext[0])).
			Return(nil)

		got, err := repo.DeleteFile(ctx, file.ID)
		assert.NoError(t, err)
		assert.Equal(t, file.ID, got.ID)
		assert.Equal(t, file.Name, got.Name)
		assert.Equal(t, file.MimeType, got.MimeType)
	})

	t.Run("UnknownFile", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		_, err = repo.DeleteFile(ctx, uuid.New())
		assert.Error(t, err)
	})

	t.Run("UnknownMimeType", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		var tags []*Tag
		var group *Group
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), false)
		assert.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 50), tags, group, user.ID)
		assert.NoError(t, err)

		sampleText := "sampleData"

		mimetype := "po"

		src := strings.NewReader(sampleText)

		name := random.AlphaNumeric(t, 20)

		file, err := repo.CreateFile(ctx, src, name, mimetype, request.ID)
		assert.NoError(t, err)

		_, err = repo.DeleteFile(ctx, file.ID)
		assert.Error(t, err)
	})

	t.Run("MissingMimeType", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		var tags []*Tag
		var group *Group
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), false)
		assert.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 50), tags, group, user.ID)
		assert.NoError(t, err)

		sampleText := "sampleData"

		mimetype := ""

		src := strings.NewReader(sampleText)

		name := random.AlphaNumeric(t, 20)

		file, err := repo.CreateFile(ctx, src, name, mimetype, request.ID)
		assert.NoError(t, err)

		_, err = repo.DeleteFile(ctx, file.ID)
		assert.Error(t, err)
	})

	t.Run("FailedToDelete", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		var tags []*Tag
		var group *Group
		user, err := repo.CreateUser(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 20), false)
		assert.NoError(t, err)
		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 50), tags, group, user.ID)
		assert.NoError(t, err)

		sampleText := "sampleData"

		mimetype := "image/png"

		src := strings.NewReader(sampleText)

		name := random.AlphaNumeric(t, 20)

		file, err := repo.CreateFile(ctx, src, name, mimetype, request.ID)
		assert.NoError(t, err)

		ext, err := mime.ExtensionsByType(file.MimeType)
		assert.NoError(t, err)

		storage.EXPECT().
			Delete(fmt.Sprintf("%s%s", file.ID.String(), ext[0])).
			Return(errors.New("failed to delete file"))

		_, err = repo.DeleteFile(ctx, file.ID)
		assert.Error(t, err)
	})
}
