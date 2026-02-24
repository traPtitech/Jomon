package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/internal/nulltime"
	"github.com/traPtitech/Jomon/internal/service"
	"github.com/traPtitech/Jomon/internal/testutil"
	"github.com/traPtitech/Jomon/internal/testutil/random"
)

func toExpectedApplicationResponse(
	t *testing.T, rd *service.ApplicationDetail,
) *service.ApplicationResponse {
	t.Helper()
	return &service.ApplicationResponse{
		ID:        rd.ID,
		Status:    rd.Status,
		CreatedAt: rd.CreatedAt,
		UpdatedAt: rd.UpdatedAt,
		CreatedBy: rd.CreatedBy,
		Title:     rd.Title,
		Content:   rd.Content,
		Tags:      rd.Tags,
		//Targets:   rd.Targets,
		//Statuses:  rd.Statuses,
	}
}

func TestEntRepository_GetApplications(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "get_applications")
	require.NoError(t, err)
	repo := NewEntRepository(client)
	client2, err := setup(t, ctx, "get_applications2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2)
	client3, err := setup(t, ctx, "get_applications3")
	require.NoError(t, err)
	repo3 := NewEntRepository(client3)
	client4, err := setup(t, ctx, "get_applications4")
	require.NoError(t, err)
	repo4 := NewEntRepository(client4)
	client5, err := setup(t, ctx, "get_applications5")
	require.NoError(t, err)
	repo5 := NewEntRepository(client5)
	client6, err := setup(t, ctx, "get_applications6")
	require.NoError(t, err)
	repo6 := NewEntRepository(client6)
	client7, err := setup(t, ctx, "get_applications7")
	require.NoError(t, err)
	repo7 := NewEntRepository(client7)
	client8, err := setup(t, ctx, "get_applications8")
	require.NoError(t, err)
	repo8 := NewEntRepository(client8)
	client9, err := setup(t, ctx, "get_applications9")
	require.NoError(t, err)
	repo9 := NewEntRepository(client9)

	t.Run("SuccessWithSortCreatedAt", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		user1, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		user2, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &service.ApplicationTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		application1, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target},
			user1.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		application2, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target},
			user2.ID)
		require.NoError(t, err)

		sort := "created_at"

		got, err := repo.GetApplications(ctx, service.ApplicationQuery{
			Sort: &sort,
		})
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.SortSlices(func(a, b *service.ApplicationResponse) bool {
				return a.ID.ID() < b.ID.ID()
			}))
		exp := []*service.ApplicationResponse{
			toExpectedApplicationResponse(t, application1),
			toExpectedApplicationResponse(t, application2),
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithReverseSortCreatedAt", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		user1, err := repo2.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		user2, err := repo2.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo2.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &service.ApplicationTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		application1, err := repo2.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target},
			user1.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		application2, err := repo2.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target},
			user2.ID)
		require.NoError(t, err)

		sort := "-created_at"

		got, err := repo2.GetApplications(ctx, service.ApplicationQuery{
			Sort: &sort,
		})
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.SortSlices(func(a, b *service.ApplicationResponse) bool {
				return a.ID.ID() < b.ID.ID()
			}))
		exp := []*service.ApplicationResponse{
			toExpectedApplicationResponse(t, application1),
			toExpectedApplicationResponse(t, application2),
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithSortTitle", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		user1, err := repo3.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		user2, err := repo3.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo3.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &service.ApplicationTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		application1, err := repo3.CreateApplication(
			ctx,
			"b",
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target},
			user1.ID)
		require.NoError(t, err)
		application2, err := repo3.CreateApplication(
			ctx,
			"a",
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target},
			user2.ID)
		require.NoError(t, err)

		sort := "title"

		got, err := repo3.GetApplications(ctx, service.ApplicationQuery{
			Sort: &sort,
		})
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.SortSlices(func(a, b *service.ApplicationResponse) bool {
				return a.ID.ID() < b.ID.ID()
			}))
		exp := []*service.ApplicationResponse{
			toExpectedApplicationResponse(t, application2),
			toExpectedApplicationResponse(t, application1),
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithReverseSortTitle", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		user1, err := repo4.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		user2, err := repo4.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo4.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &service.ApplicationTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		application1, err := repo4.CreateApplication(
			ctx,
			"b",
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target},
			user1.ID)
		require.NoError(t, err)
		application2, err := repo4.CreateApplication(
			ctx,
			"a",
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target},
			user2.ID)
		require.NoError(t, err)

		sort := "-title"

		got, err := repo4.GetApplications(ctx, service.ApplicationQuery{
			Sort: &sort,
		})
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.SortSlices(func(a, b *service.ApplicationResponse) bool {
				return a.ID.ID() < b.ID.ID()
			}))
		exp := []*service.ApplicationResponse{
			toExpectedApplicationResponse(t, application1),
			toExpectedApplicationResponse(t, application2),
		}
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithQueryTarget", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		user1, err := repo5.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		user2, err := repo5.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo5.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target1 := &service.ApplicationTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}
		target2 := &service.ApplicationTarget{
			Target: user2.ID,
			Amount: random.Numeric(t, 10000),
		}

		application1, err := repo5.CreateApplication(
			ctx,
			"b",
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target1},
			user1.ID)
		require.NoError(t, err)
		_, err = repo5.CreateApplication(
			ctx,
			"a",
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target2},
			user2.ID)
		require.NoError(t, err)

		target := target1.Target
		got, err := repo5.GetApplications(ctx, service.ApplicationQuery{
			Target: target,
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		exp := toExpectedApplicationResponse(t, application1)
		opts := testutil.ApproxEqualOptions()
		testutil.RequireEqual(t, exp, got[0], opts...)
	})

	t.Run("SuccessWithQuerySince", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		user1, err := repo6.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		user2, err := repo6.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo6.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &service.ApplicationTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		application1, err := repo6.CreateApplication(
			ctx,
			"b",
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target},
			user1.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		application2, err := repo6.CreateApplication(
			ctx,
			"a",
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target},
			user2.ID)
		require.NoError(t, err)

		since := application1.CreatedAt.Add(10 * time.Millisecond)
		got, err := repo6.GetApplications(ctx, service.ApplicationQuery{
			Since: nulltime.FromTime(&since),
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		exp := toExpectedApplicationResponse(t, application2)
		opts := testutil.ApproxEqualOptions()
		testutil.RequireEqual(t, exp, got[0], opts...)
	})

	t.Run("SuccessWithQueryUntil", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		user1, err := repo7.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		user2, err := repo7.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo7.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &service.ApplicationTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		application1, err := repo7.CreateApplication(
			ctx,
			"b",
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target},
			user1.ID)
		require.NoError(t, err)
		time.Sleep(2 * time.Second)
		application2, err := repo7.CreateApplication(
			ctx,
			"a",
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target},
			user2.ID)
		require.NoError(t, err)

		until := application2.CreatedAt.Add(-1 * time.Second)
		got, err := repo7.GetApplications(ctx, service.ApplicationQuery{
			Until: nulltime.FromTime(&until),
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		exp := toExpectedApplicationResponse(t, application1)
		opts := testutil.ApproxEqualOptions()
		testutil.RequireEqual(t, exp, got[0], opts...)
	})

	t.Run("SuccessWithQueryStatus", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		user1, err := repo8.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		user2, err := repo8.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo8.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &service.ApplicationTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		application1, err := repo8.CreateApplication(
			ctx,
			"b",
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target},
			user1.ID)
		require.NoError(t, err)
		time.Sleep(2 * time.Second)
		_, err = repo8.CreateApplication(
			ctx,
			"a",
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag},
			[]*service.ApplicationTarget{target},
			user2.ID)
		require.NoError(t, err)

		time.Sleep(1 * time.Second)

		status := "accepted"
		_, err = repo8.CreateStatus(ctx, application1.ID, user1.ID, service.Accepted)
		require.NoError(t, err)

		got, err := repo8.GetApplications(ctx, service.ApplicationQuery{
			Status: &status,
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		exp := toExpectedApplicationResponse(t, application1)
		exp.Status = service.Accepted
		opts := testutil.ApproxEqualOptions()
		testutil.RequireEqual(t, exp, got[0], opts...)
	})

	t.Run("SuccessWithQueryCreatedBy", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		user1, err := repo9.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		user2, err := repo9.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		target := &service.ApplicationTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		application1, err := repo9.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*service.Tag{},
			[]*service.ApplicationTarget{target},
			user1.ID)
		require.NoError(t, err)
		_, err = repo9.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*service.Tag{},
			[]*service.ApplicationTarget{target},
			user2.ID)
		require.NoError(t, err)

		got, err := repo9.GetApplications(ctx, service.ApplicationQuery{
			CreatedBy: user1.ID,
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		exp := toExpectedApplicationResponse(t, application1)
		opts := testutil.ApproxEqualOptions()
		testutil.RequireEqual(t, exp, got[0], opts...)
	})
}

func TestEntRepository_CreateApplication(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "create_application")
	require.NoError(t, err)
	repo := NewEntRepository(client)
	client2, err := setup(t, ctx, "create_application2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2)
	client3, err := setup(t, ctx, "create_application3")
	require.NoError(t, err)
	repo3 := NewEntRepository(client3)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		title := random.AlphaNumeric(t, 40)
		content := random.AlphaNumeric(t, 100)
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &service.ApplicationTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}

		application, err := repo.CreateApplication(
			ctx,
			title, content,
			[]*service.Tag{tag}, []*service.ApplicationTarget{target},
			user.ID)
		require.NoError(t, err)
		exp := &service.ApplicationDetail{
			Status:  service.Submitted,
			Title:   title,
			Content: content,
			Tags:    []*service.Tag{tag},
			Targets: []*service.ApplicationTargetDetail{{
				Target: target.Target,
				Amount: target.Amount,
			}},
			Statuses: []*service.ApplicationStatus{{
				CreatedBy: user.ID,
				Status:    service.Submitted,
			}},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: user.ID,
		}
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(service.ApplicationDetail{}, "ID"),
			cmpopts.IgnoreFields(service.ApplicationTargetDetail{}, "ID", "PaidAt", "CreatedAt"),
			cmpopts.IgnoreFields(service.ApplicationStatus{}, "ID", "CreatedAt"))
		testutil.AssertEqual(t, exp, application, opts...)
	})

	t.Run("UnknownUser", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		title := random.AlphaNumeric(t, 40)
		content := random.AlphaNumeric(t, 100)
		tag, err := repo2.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)

		_, err = repo2.CreateApplication(
			ctx,
			title, content,
			[]*service.Tag{tag}, []*service.ApplicationTarget{},
			uuid.New())
		require.Error(t, err)
	})

	t.Run("UnknownTag", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		title := random.AlphaNumeric(t, 40)
		content := random.AlphaNumeric(t, 100)
		user, err := repo3.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)

		date := time.Now()
		tag := &service.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}

		_, err = repo3.CreateApplication(
			ctx,
			title, content,
			[]*service.Tag{tag}, []*service.ApplicationTarget{},
			user.ID)
		require.Error(t, err)
	})
}

func TestEntRepository_GetApplication(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "get_application")
	require.NoError(t, err)
	repo := NewEntRepository(client)
	client2, err := setup(t, ctx, "get_application2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &service.ApplicationTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		require.NoError(t, err)
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag}, []*service.ApplicationTarget{target},
			user.ID)
		require.NoError(t, err)

		got, err := repo.GetApplication(ctx, application.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		testutil.AssertEqual(t, application, got, opts...)
	})

	t.Run("UnknownApplication", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		_, err := repo2.GetApplication(ctx, uuid.New())
		require.Error(t, err)
	})
}

func TestEntRepository_UpdateApplication(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, err := setup(t, ctx, "update_application")
	require.NoError(t, err)
	repo := NewEntRepository(client)
	client2, err := setup(t, ctx, "update_application2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2)
	client3, err := setup(t, ctx, "update_application3")
	require.NoError(t, err)
	repo3 := NewEntRepository(client3)
	client4, err := setup(t, ctx, "update_application4")
	require.NoError(t, err)
	repo4 := NewEntRepository(client4)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		user, err := repo.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &service.ApplicationTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		application, err := repo.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag}, []*service.ApplicationTarget{target},
			user.ID)
		require.NoError(t, err)
		// CreatedAt の差を1秒以内に収めるためにここで time.Now を取る
		expTarget := &service.ApplicationTargetDetail{
			Target:    target.Target,
			Amount:    target.Amount,
			CreatedAt: time.Now(),
		}

		updatedApplication, err := repo.UpdateApplication(
			ctx,
			application.ID, application.Title, application.Content,
			[]*service.Tag{tag}, []*service.ApplicationTarget{target})
		require.NoError(t, err)
		exp := &service.ApplicationDetail{
			ID:        application.ID,
			Status:    application.Status,
			Title:     application.Title,
			Content:   application.Content,
			Comments:  application.Comments,
			Files:     application.Files,
			Tags:      []*service.Tag{tag},
			Targets:   []*service.ApplicationTargetDetail{expTarget},
			Statuses:  application.Statuses,
			CreatedAt: application.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: application.CreatedBy,
		}
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(service.ApplicationTargetDetail{}, "ID", "PaidAt"))
		testutil.AssertEqual(t, exp, updatedApplication, opts...)
	})

	t.Run("Success2", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		user, err := repo2.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo2.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &service.ApplicationTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		application, err := repo2.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag}, []*service.ApplicationTarget{target},
			user.ID)
		require.NoError(t, err)
		// CreatedAt の差を1秒以内に収めるためにここで time.Now を取る
		expTarget := &service.ApplicationTargetDetail{
			Target:    target.Target,
			Amount:    target.Amount,
			CreatedAt: time.Now(),
		}

		title := random.AlphaNumeric(t, 40)
		updatedApplication, err := repo2.UpdateApplication(
			ctx,
			application.ID, title, application.Content,
			[]*service.Tag{tag}, []*service.ApplicationTarget{target})
		require.NoError(t, err)
		exp := &service.ApplicationDetail{
			ID:        application.ID,
			Status:    application.Status,
			Title:     title,
			Content:   application.Content,
			Comments:  application.Comments,
			Files:     application.Files,
			Tags:      []*service.Tag{tag},
			Targets:   []*service.ApplicationTargetDetail{expTarget},
			Statuses:  application.Statuses,
			CreatedAt: application.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: application.CreatedBy,
		}
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(service.ApplicationTargetDetail{}, "ID", "PaidAt"))
		testutil.AssertEqual(t, exp, updatedApplication, opts...)
	})

	t.Run("Success3", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		user, err := repo3.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo3.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &service.ApplicationTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		application, err := repo3.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag}, []*service.ApplicationTarget{target},
			user.ID)
		require.NoError(t, err)
		// CreatedAt の差を1秒以内に収めるためにここで time.Now を取る
		expTarget := &service.ApplicationTargetDetail{
			Target:    target.Target,
			Amount:    target.Amount,
			CreatedAt: time.Now(),
		}
		content := random.AlphaNumeric(t, 100)
		updatedApplication, err := repo3.UpdateApplication(
			ctx,
			application.ID, application.Title, content,
			[]*service.Tag{tag}, []*service.ApplicationTarget{target})
		require.NoError(t, err)
		exp := &service.ApplicationDetail{
			ID:        application.ID,
			Status:    application.Status,
			Title:     application.Title,
			Content:   content,
			Comments:  application.Comments,
			Files:     application.Files,
			Tags:      []*service.Tag{tag},
			Targets:   []*service.ApplicationTargetDetail{expTarget},
			Statuses:  application.Statuses,
			CreatedAt: application.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: application.CreatedBy,
		}
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(service.ApplicationTargetDetail{}, "ID", "PaidAt"))
		testutil.AssertEqual(t, exp, updatedApplication, opts...)
	})

	t.Run("UnknownTag", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		user, err := repo4.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo4.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &service.ApplicationTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		application, err := repo4.CreateApplication(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*service.Tag{tag}, []*service.ApplicationTarget{target},
			user.ID)
		require.NoError(t, err)

		date := time.Now()
		unknownTag := &service.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		_, err = repo4.UpdateApplication(
			ctx,
			application.ID, application.Title, application.Content,
			[]*service.Tag{unknownTag}, []*service.ApplicationTarget{target})
		require.Error(t, err)
	})
}
