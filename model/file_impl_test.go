package model

import (
	"context"
	"fmt"
	"mime"
	"strings"
	"testing"

	"github.com/google/uuid"
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

		request, err := repo.CreateRequest(ctx, random.Numeric(t, 100000), random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 50), nil, nil, nil)
		assert.NoError(t, err)

		sampleText := "sampleData"

		fileID := uuid.New()
		mimetype := "image/png"

		ext, err := mime.ExtensionsByType(mimetype)
		assert.NoError(t, err)
		assert.True(t, len(ext) != 0)

		filename := fmt.Sprintf("%s%s", fileID.String(), ext[0])

		src := strings.NewReader(sampleText)

		storage.EXPECT().
			Save(filename, src).
			Return(nil)

		file, err := repo.CreateFile(ctx, src, random.AlphaNumeric(t, 20), mimetype, request.ID)
		assert.NoError(t, err)
		assert.Equal(t, filename, file.Name)
		assert.Equal(t, mimetype, file.MimeType)
	})
}
