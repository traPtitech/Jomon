package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/testutil"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestEntRepository_CreateFile(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "create_file")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		var tags []*Tag
		var targets []*ApplicationTarget
		var group *Group
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 50),
			tags, targets,
			group, user.ID)
		require.NoError(t, err)

		mimetype := "image/png"

		name := random.AlphaNumeric(t, 20)

		file, err := repo.CreateFile(ctx, name, mimetype, application.ID, user.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(File{}, "ID"))
		exp := &File{
			Name:      name,
			MimeType:  mimetype,
			CreatedBy: user.ID,
			CreatedAt: time.Now(),
		}
		testutil.RequireEqual(t, exp, file, opts...)
	})

	t.Run("UnknownApplication", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application := Application{
			ID: uuid.New(),
		}

		mimetype := "image/png"

		name := random.AlphaNumeric(t, 20)

		_, err = repo.CreateFile(ctx, name, mimetype, application.ID, user.ID)
		require.Error(t, err)
	})

	t.Run("MissingName", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		var tags []*Tag
		var targets []*ApplicationTarget
		var group *Group
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 50),
			tags, targets,
			group, user.ID)
		require.NoError(t, err)

		mimetype := "image/png"

		_, err = repo.CreateFile(ctx, "", mimetype, application.ID, user.ID)
		require.Error(t, err)
	})
}

func TestEntRepository_GetFile(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "get_file")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		var tags []*Tag
		var targets []*ApplicationTarget
		var group *Group
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 50),
			tags, targets,
			group, user.ID)
		require.NoError(t, err)

		mimetype := "image/png"

		name := random.AlphaNumeric(t, 20)

		file, err := repo.CreateFile(ctx, name, mimetype, application.ID, user.ID)
		require.NoError(t, err)
		got, err := repo.GetFile(ctx, file.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &File{
			ID:        file.ID,
			Name:      name,
			MimeType:  mimetype,
			CreatedBy: user.ID,
			CreatedAt: file.CreatedAt,
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("UnknownFile", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		_, err = repo.GetFile(ctx, uuid.New())
		require.Error(t, err)
	})
}

func TestEntRepository_DeleteFile(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "delete_file")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		var tags []*Tag
		var targets []*ApplicationTarget
		var group *Group
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 50),
			tags, targets,
			group, user.ID)
		require.NoError(t, err)

		mimetype := "image/png"

		name := random.AlphaNumeric(t, 20)

		file, err := repo.CreateFile(ctx, name, mimetype, application.ID, user.ID)
		require.NoError(t, err)

		err = repo.DeleteFile(ctx, file.ID)
		require.NoError(t, err)

		r, err := repo.GetApplication(ctx, application.ID)
		require.NoError(t, err)
		require.Empty(t, r.Files)
	})

	t.Run("UnknownFile", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)

		err = repo.DeleteFile(ctx, uuid.New())
		require.Error(t, err)
	})
}
