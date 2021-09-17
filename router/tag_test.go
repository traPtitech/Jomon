package router

/*
// TODO: 直す
func TestHandlers_GetTags(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)
		date := time.Now()

		tag1 := &model.Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		tag2 := &model.Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		tags := []*model.Tag{tag1, tag2}

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			GetTags(ctx).
			Return(tags, nil)

		var resBody TagResponse
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.GET, "/api/tags", nil, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Len(t, resBody.Tags, 2)
		if resBody.Tags[0].ID == tag1.ID {
			assert.Equal(t, resBody.Tags[0].ID, tag1.ID)
			assert.Equal(t, resBody.Tags[0].Name, tag1.Name)
			assert.Equal(t, resBody.Tags[0].Description, tag1.Description)
			assert.Equal(t, resBody.Tags[1].ID, tag2.ID)
			assert.Equal(t, resBody.Tags[1].Name, tag2.Name)
			assert.Equal(t, resBody.Tags[1].Description, tag2.Description)
		} else {
			assert.Equal(t, resBody.Tags[0].ID, tag2.ID)
			assert.Equal(t, resBody.Tags[0].Name, tag2.Name)
			assert.Equal(t, resBody.Tags[0].Description, tag2.Description)
			assert.Equal(t, resBody.Tags[1].ID, tag1.ID)
			assert.Equal(t, resBody.Tags[1].Name, tag1.Name)
			assert.Equal(t, resBody.Tags[1].Description, tag1.Description)
		}
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)

		tags := []*model.Tag{}

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			GetTags(ctx).
			Return(tags, nil)

		var resBody TagResponse
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.GET, "/api/tags", nil, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Len(t, resBody.Tags, 0)
	})

	t.Run("FailedToGetTags", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			GetTags(ctx).
			Return(nil, errors.New("Failed to get tags."))

		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.GET, "/api/tags", nil, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}

// TODO: 直す
func TestHandlers_PostTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)
		date := time.Now()

		tag := &model.Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			CreateTag(ctx, tag.Name, tag.Description).
			Return(tag, nil)

		req := Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}

		var resBody TagOverview
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.POST, "/api/tags", &req, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, tag.ID, resBody.ID)
		assert.Equal(t, tag.Name, resBody.Name)
		assert.Equal(t, tag.Description, resBody.Description)
	})

	t.Run("MissingName", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)
		date := time.Now()

		tag := &model.Tag{
			ID:          uuid.New(),
			Name:        "",
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			CreateTag(ctx, tag.Name, tag.Description).
			Return(nil, errors.New("Tag name can't be empty."))

		req := Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}

		var resBody TagOverview
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.POST, "/api/tags", &req, &resBody)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})
}

// TODO: 直す
func TestHandlers_PutTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)
		date := time.Now()

		tag := &model.Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			UpdateTag(ctx, tag.ID, tag.Name, tag.Description).
			Return(tag, nil)

		req := Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}

		var resBody TagOverview
		path := fmt.Sprintf("/api/tags/%s", tag.ID.String())
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.PUT, path, &req, &resBody)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, tag.ID, resBody.ID)
		assert.Equal(t, tag.Name, resBody.Name)
		assert.Equal(t, tag.Description, resBody.Description)
	})

	t.Run("MissingName", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)
		date := time.Now()

		tag := &model.Tag{
			ID:          uuid.New(),
			Name:        "",
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			UpdateTag(ctx, tag.ID, tag.Name, tag.Description).
			Return(nil, errors.New("Tag name can't be empty."))

		req := Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}

		path := fmt.Sprintf("/api/tags/%s", tag.ID.String())
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.PUT, path, &req, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)
		date := time.Now()

		tag := &model.Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		req := Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}

		path := "/api/tags/hoge" // Invalid UUID
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.PUT, path, &req, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)
		date := time.Now()

		tag := &model.Tag{
			ID:          uuid.Nil,
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		req := Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}

		path := fmt.Sprintf("/api/tags/%s", tag.ID.String())
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.PUT, path, &req, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})
}

// TODO: 直す
func TestHandlers_DeleteTag(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)

		id := uuid.New()

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			DeleteTag(ctx, id).
			Return(nil)

		path := fmt.Sprintf("/api/tags/%s", id.String())
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.DELETE, path, nil, nil)
		assert.Equal(t, http.StatusOK, statusCode)
	})

	t.Run("UnknownID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)

		id := uuid.New()

		ctx := context.Background()
		th.Repository.MockTagRepository.
			EXPECT().
			DeleteTag(ctx, id).
			Return(errors.New("Tag not found"))

		path := fmt.Sprintf("/api/tags/%s", id.String())
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.DELETE, path, nil, nil)
		assert.Equal(t, http.StatusInternalServerError, statusCode)
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)

		path := "/api/tags/hoge"
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.DELETE, path, nil, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})

	t.Run("NilUUID", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		accessUser := makeUser(t, false)
		th, err := NewTestServer(t, ctrl, accessUser)
		assert.NoError(t, err)

		path := fmt.Sprintf("/api/tags/%s", uuid.Nil)
		statusCode, _ := th.doRequestWithLogin(t, accessUser, echo.DELETE, path, nil, nil)
		assert.Equal(t, http.StatusBadRequest, statusCode)
	})
}
*/
