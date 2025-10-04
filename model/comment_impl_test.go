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

func TestEntRepository_GetComments(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "get_comments")
	require.NoError(t, err)
	repo := NewEntRepository(client)
	client2, err := setup(t, ctx, "get_comments2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{}, []*ApplicationTarget{},
			user.ID)
		require.NoError(t, err)

		comment1, err := repo.CreateComment(
			ctx,
			random.AlphaNumeric(t, 30),
			application.ID,
			user.ID,
		)
		require.NoError(t, err)
		comment2, err := repo.CreateComment(
			ctx,
			random.AlphaNumeric(t, 30),
			application.ID,
			user.ID,
		)
		require.NoError(t, err)

		got, err := repo.GetComments(ctx, application.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.SortSlices(func(l, r *Comment) bool {
				return l.ID.ID() < r.ID.ID()
			}))
		exp := []*Comment{comment1, comment2}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		user, err := repo2.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo2.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{}, []*ApplicationTarget{},
			user.ID)
		require.NoError(t, err)

		got, err := repo2.GetComments(ctx, application.ID)
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("UnknownApplication", func(t *testing.T) {
		t.Parallel()
		_, err := repo.GetComments(ctx, uuid.New())
		require.Error(t, err)
	})
}

func TestEntRepository_CreateComment(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "create_comment")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{}, []*ApplicationTarget{},
			user.ID)
		require.NoError(t, err)

		comment := random.AlphaNumeric(t, 30)
		user2, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		created, err := repo.CreateComment(ctx, comment, application.ID, user2.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(Comment{}, "ID"))
		exp := &Comment{
			User:      user2.ID,
			Comment:   comment,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		testutil.RequireEqual(t, exp, created, opts...)
	})

	t.Run("UnknownApplication", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		_, err = repo.CreateComment(ctx, random.AlphaNumeric(t, 30), uuid.New(), user.ID)
		require.Error(t, err)
	})

	t.Run("UnknownUser", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{}, []*ApplicationTarget{},
			user.ID)
		require.NoError(t, err)

		_, err = repo.CreateComment(ctx, random.AlphaNumeric(t, 30), application.ID, uuid.New())
		require.Error(t, err)
	})
}

func TestEntRepository_UpdateComment(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "update_comment")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{}, []*ApplicationTarget{},
			user.ID)
		require.NoError(t, err)
		created, err := repo.CreateComment(ctx, random.AlphaNumeric(t, 30), application.ID, user.ID)
		require.NoError(t, err)

		comment := random.AlphaNumeric(t, 30)
		updated, err := repo.UpdateComment(ctx, comment, application.ID, created.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &Comment{
			ID:        created.ID,
			User:      created.User,
			Comment:   comment,
			CreatedAt: created.CreatedAt,
			UpdatedAt: time.Now(),
		}
		testutil.RequireEqual(t, exp, updated, opts...)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application1, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{}, []*ApplicationTarget{},
			user.ID)
		require.NoError(t, err)
		comment, err := repo.CreateComment(
			ctx,
			random.AlphaNumeric(t, 30),
			application1.ID,
			user.ID,
		)
		require.NoError(t, err)

		application2, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{}, []*ApplicationTarget{},
			user.ID)
		require.NoError(t, err)
		updated, err := repo.UpdateComment(ctx, comment.Comment, application2.ID, comment.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := &Comment{
			ID:        comment.ID,
			User:      comment.User,
			Comment:   comment.Comment,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: time.Now(),
		}
		testutil.RequireEqual(t, exp, updated, opts...)

		got, err := repo.GetComments(ctx, application2.ID)
		require.NoError(t, err)
		testutil.RequireEqual(t, []*Comment{updated}, got, opts...)
	})

	t.Run("UnknownComment", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{}, []*ApplicationTarget{},
			user.ID)
		require.NoError(t, err)

		_, err = repo.UpdateComment(
			ctx,
			random.AlphaNumeric(t, 30),
			application.ID, uuid.New())
		require.Error(t, err)
	})

	t.Run("UnknownApplication", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{}, []*ApplicationTarget{},
			user.ID)
		require.NoError(t, err)
		comment, err := repo.CreateComment(ctx, random.AlphaNumeric(t, 30), application.ID, user.ID)
		require.NoError(t, err)

		_, err = repo.UpdateComment(ctx, random.AlphaNumeric(t, 30), uuid.New(), comment.ID)
		require.Error(t, err)
	})
}

func TestEntRepository_DeleteComment(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "delete_comment")
	require.NoError(t, err)
	repo := NewEntRepository(client)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{}, []*ApplicationTarget{},
			user.ID)
		require.NoError(t, err)
		comment, err := repo.CreateComment(ctx, random.AlphaNumeric(t, 30), application.ID, user.ID)
		require.NoError(t, err)

		err = repo.DeleteComment(ctx, application.ID, comment.ID)
		require.NoError(t, err)

		comments, err := repo.GetComments(ctx, application.ID)
		require.NoError(t, err)
		require.Empty(t, comments)
	})

	t.Run("UnknownApplication", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{}, []*ApplicationTarget{},
			user.ID)
		require.NoError(t, err)
		comment, err := repo.CreateComment(ctx, random.AlphaNumeric(t, 30), application.ID, user.ID)
		require.NoError(t, err)

		err = repo.DeleteComment(ctx, uuid.New(), comment.ID)
		require.Error(t, err)
	})

	t.Run("UnknownComment", func(t *testing.T) {
		t.Parallel()
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{}, []*ApplicationTarget{},
			user.ID)
		require.NoError(t, err)

		err = repo.DeleteComment(ctx, application.ID, uuid.New())
		require.Error(t, err)
	})
}
