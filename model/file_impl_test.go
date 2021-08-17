package model

import (
	"context"
	"fmt"
	"mime"
	"strings"
	"testing"

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
}
