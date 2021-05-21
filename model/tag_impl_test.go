package model

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_GetTags(t *testing.T) {
	client, _ := setup(t)
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
