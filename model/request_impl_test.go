package model

/*
func TestEntRepository_GetRequests(t *testing.T) {
	client, storage, err := setup(t)
	assert.NoError(t, err)
	repo := NewEntRepository(client, storage)
	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		user1, _ := repo.GetUser(ctx) // TODO: impl
		tag1, _ := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		tag2, _ := repo.CreateTag(ctx, random.AlphaNumeric(t, 20), random.AlphaNumeric(t, 30))
		request1, _ := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 1000), []*Tag{tag1}, nil, user1.ID)
		request2, _ := repo.CreateRequest(ctx, random.Numeric(t, 1000000), random.AlphaNumeric(t, 40), random.AlphaNumeric(t, 1000), []*Tag{tag1}, nil, user1.ID)
		sort := "created_at"
		got, err := repo.GetRequests(ctx, RequestQuery{
			Sort: &sort,
		})
		assert.NoError(t, err)
		// TODO: impl
	})
}
*/
