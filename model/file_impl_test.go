package model

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_CreateFile(t *testing.T) {
	client, storage, err := setup(t)
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
