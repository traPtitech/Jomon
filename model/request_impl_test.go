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

func (rd *RequestDetail) toExpectedRequestResponse(t *testing.T) *RequestResponse {
	t.Helper()
	return &RequestResponse{
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
		Group: rd.Group,
	}
}

func TestEntRepository_GetRequests(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, storage, err := setup(t, ctx, "get_requests")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "get_requests2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)
	client3, storage3, err := setup(t, ctx, "get_requests3")
	require.NoError(t, err)
	repo3 := NewEntRepository(client3, storage3)
	client4, storage4, err := setup(t, ctx, "get_requests4")
	require.NoError(t, err)
	repo4 := NewEntRepository(client4, storage4)
	client5, storage5, err := setup(t, ctx, "get_requests5")
	require.NoError(t, err)
	repo5 := NewEntRepository(client5, storage5)
	client6, storage6, err := setup(t, ctx, "get_requests6")
	require.NoError(t, err)
	repo6 := NewEntRepository(client6, storage6)
	client7, storage7, err := setup(t, ctx, "get_requests7")
	require.NoError(t, err)
	repo7 := NewEntRepository(client7, storage7)
	client8, storage8, err := setup(t, ctx, "get_requests8")
	require.NoError(t, err)
	repo8 := NewEntRepository(client8, storage8)
	client9, storage9, err := setup(t, ctx, "get_requests9")
	require.NoError(t, err)
	repo9 := NewEntRepository(client9, storage9)

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
		target := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)

		request1, err := repo.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target},
			group,
			user1.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		request2, err := repo.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target},
			group,
			user2.ID)
		require.NoError(t, err)

		sort := "created_at"

		got, err := repo.GetRequests(ctx, RequestQuery{
			Sort: &sort,
		})
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.SortSlices(func(a, b *RequestResponse) bool {
				return a.ID.ID() < b.ID.ID()
			}))
		exp := []*RequestResponse{
			request1.toExpectedRequestResponse(t),
			request2.toExpectedRequestResponse(t),
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
		target := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo2.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)

		request1, err := repo2.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target},
			group,
			user1.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		request2, err := repo2.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target},
			group,
			user2.ID)
		require.NoError(t, err)

		sort := "-created_at"

		got, err := repo2.GetRequests(ctx, RequestQuery{
			Sort: &sort,
		})
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.SortSlices(func(a, b *RequestResponse) bool {
				return a.ID.ID() < b.ID.ID()
			}))
		exp := []*RequestResponse{
			request1.toExpectedRequestResponse(t),
			request2.toExpectedRequestResponse(t),
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
		target := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo3.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)

		request1, err := repo3.CreateRequest(
			ctx,
			"b",
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target},
			group,
			user1.ID)
		require.NoError(t, err)
		request2, err := repo3.CreateRequest(
			ctx,
			"a",
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target},
			group,
			user2.ID)
		require.NoError(t, err)

		sort := "title"

		got, err := repo3.GetRequests(ctx, RequestQuery{
			Sort: &sort,
		})
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.SortSlices(func(a, b *RequestResponse) bool {
				return a.ID.ID() < b.ID.ID()
			}))
		exp := []*RequestResponse{
			request2.toExpectedRequestResponse(t),
			request1.toExpectedRequestResponse(t),
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
		target := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo4.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)

		request1, err := repo4.CreateRequest(
			ctx,
			"b",
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target},
			group,
			user1.ID)
		require.NoError(t, err)
		request2, err := repo4.CreateRequest(
			ctx,
			"a",
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target},
			group,
			user2.ID)
		require.NoError(t, err)

		sort := "-title"

		got, err := repo4.GetRequests(ctx, RequestQuery{
			Sort: &sort,
		})
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.SortSlices(func(a, b *RequestResponse) bool {
				return a.ID.ID() < b.ID.ID()
			}))
		exp := []*RequestResponse{
			request1.toExpectedRequestResponse(t),
			request2.toExpectedRequestResponse(t),
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
		target1 := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}
		target2 := &RequestTarget{
			Target: user2.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo5.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)

		request1, err := repo5.CreateRequest(
			ctx,
			"b",
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target1},
			group,
			user1.ID)
		require.NoError(t, err)
		_, err = repo5.CreateRequest(
			ctx,
			"a",
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target2},
			group,
			user2.ID)
		require.NoError(t, err)

		target := target1.Target
		got, err := repo5.GetRequests(ctx, RequestQuery{
			Target: target,
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		exp := request1.toExpectedRequestResponse(t)
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
		target := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo6.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)

		request1, err := repo6.CreateRequest(
			ctx,
			"b",
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target},
			group,
			user1.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		request2, err := repo6.CreateRequest(
			ctx,
			"a",
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target},
			group,
			user2.ID)
		require.NoError(t, err)

		since := request1.CreatedAt.Add(10 * time.Millisecond)
		got, err := repo6.GetRequests(ctx, RequestQuery{
			Since: &since,
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		exp := request2.toExpectedRequestResponse(t)
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
		target := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo7.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)

		request1, err := repo7.CreateRequest(
			ctx,
			"b",
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target},
			group,
			user1.ID)
		require.NoError(t, err)
		time.Sleep(2 * time.Second)
		request2, err := repo7.CreateRequest(
			ctx,
			"a",
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target},
			group,
			user2.ID)
		require.NoError(t, err)

		until := request2.CreatedAt.Add(-1 * time.Second)
		got, err := repo7.GetRequests(ctx, RequestQuery{
			Until: &until,
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		exp := request1.toExpectedRequestResponse(t)
		exp.Group.UpdatedAt = request2.Group.UpdatedAt
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
		target := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo8.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)

		request1, err := repo8.CreateRequest(
			ctx,
			"b",
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target},
			group,
			user1.ID)
		require.NoError(t, err)
		time.Sleep(2 * time.Second)
		request2, err := repo8.CreateRequest(
			ctx,
			"a",
			random.AlphaNumeric(t, 100),
			[]*Tag{tag},
			[]*RequestTarget{target},
			group,
			user2.ID)
		require.NoError(t, err)

		time.Sleep(1 * time.Second)

		status := "accepted"
		_, err = repo8.CreateStatus(ctx, request1.ID, user1.ID, Accepted)
		require.NoError(t, err)

		got, err := repo8.GetRequests(ctx, RequestQuery{
			Status: &status,
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		exp := request1.toExpectedRequestResponse(t)
		exp.Status = Accepted
		exp.Group.UpdatedAt = request2.Group.UpdatedAt
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
		target := &RequestTarget{
			Target: user1.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo9.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)

		request1, err := repo9.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{},
			[]*RequestTarget{target},
			group,
			user1.ID)
		require.NoError(t, err)
		_, err = repo9.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{},
			[]*RequestTarget{target},
			group,
			user2.ID)
		require.NoError(t, err)

		got, err := repo9.GetRequests(ctx, RequestQuery{
			CreatedBy: user1.ID,
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		exp := request1.toExpectedRequestResponse(t)
		opts := testutil.ApproxEqualOptions()
		testutil.RequireEqual(t, exp, got[0], opts...)
	})
}

func TestEntRepository_CreateRequest(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, storage, err := setup(t, ctx, "create_request")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "create_request2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)
	client3, storage3, err := setup(t, ctx, "create_request3")
	require.NoError(t, err)
	repo3 := NewEntRepository(client3, storage3)
	client4, storage4, err := setup(t, ctx, "create_request4")
	require.NoError(t, err)
	repo4 := NewEntRepository(client4, storage4)

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
		target := &RequestTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}

		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)

		request, err := repo.CreateRequest(
			ctx,
			title, content,
			[]*Tag{tag}, []*RequestTarget{target},
			group, user.ID)
		require.NoError(t, err)
		exp := &RequestDetail{
			Status:  Submitted,
			Title:   title,
			Content: content,
			Tags:    []*Tag{tag},
			Targets: []*RequestTargetDetail{{
				Target: target.Target,
				Amount: target.Amount,
			}},
			Statuses: []*RequestStatus{{
				CreatedBy: user.ID,
				Status:    Submitted,
			}},
			Group:     group,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: user.ID,
		}
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(RequestDetail{}, "ID"),
			cmpopts.IgnoreFields(RequestTargetDetail{}, "ID", "PaidAt", "CreatedAt"),
			cmpopts.IgnoreFields(RequestStatus{}, "ID", "CreatedAt"))
		testutil.AssertEqual(t, exp, request, opts...)
	})

	t.Run("UnknownUser", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		title := random.AlphaNumeric(t, 40)
		content := random.AlphaNumeric(t, 100)
		tag, err := repo2.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)

		budget := random.Numeric(t, 10000)
		group, err := repo2.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)

		_, err = repo2.CreateRequest(
			ctx,
			title, content,
			[]*Tag{tag}, []*RequestTarget{},
			group, uuid.New())
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
		tag := &Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}

		budget := random.Numeric(t, 10000)
		group, err := repo3.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)

		_, err = repo3.CreateRequest(
			ctx,
			title, content,
			[]*Tag{tag}, []*RequestTarget{},
			group, user.ID)
		require.Error(t, err)
	})

	t.Run("UnknownGroup", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		title := random.AlphaNumeric(t, 40)
		content := random.AlphaNumeric(t, 100)
		user, err := repo4.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo4.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)

		date := time.Now()
		budget := random.Numeric(t, 100000)
		group := &Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 20),
			Budget:      &budget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}

		_, err = repo4.CreateRequest(
			ctx,
			title, content,
			[]*Tag{tag}, []*RequestTarget{},
			group, user.ID)
		require.Error(t, err)
	})
}

func TestEntRepository_GetRequest(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, storage, err := setup(t, ctx, "get_request")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "get_request2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)

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
		target := &RequestTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)
		request, err := repo.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{tag}, []*RequestTarget{target},
			group, user.ID)
		require.NoError(t, err)

		got, err := repo.GetRequest(ctx, request.ID)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		testutil.AssertEqual(t, request, got, opts...)
	})

	t.Run("UnknownRequest", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		_, err := repo2.GetRequest(ctx, uuid.New())
		require.Error(t, err)
	})
}

func TestEntRepository_UpdateRequest(t *testing.T) {
	ctx := testutil.NewContext(t)
	client, storage, err := setup(t, ctx, "update_request")
	require.NoError(t, err)
	repo := NewEntRepository(client, storage)
	client2, storage2, err := setup(t, ctx, "update_request2")
	require.NoError(t, err)
	repo2 := NewEntRepository(client2, storage2)
	client3, storage3, err := setup(t, ctx, "update_request3")
	require.NoError(t, err)
	repo3 := NewEntRepository(client3, storage3)
	client4, storage4, err := setup(t, ctx, "update_request4")
	require.NoError(t, err)
	repo4 := NewEntRepository(client4, storage4)
	client5, storage5, err := setup(t, ctx, "update_request5")
	require.NoError(t, err)
	repo5 := NewEntRepository(client5, storage5)

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
		budget := random.Numeric(t, 10000)
		group, err := repo.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		request, err := repo.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{tag}, []*RequestTarget{target},
			group, user.ID)
		require.NoError(t, err)
		// CreatedAt の差を1秒以内に収めるためにここで time.Now を取る
		expTarget := &RequestTargetDetail{
			Target:    target.Target,
			Amount:    target.Amount,
			CreatedAt: time.Now(),
		}

		updatedRequest, err := repo.UpdateRequest(
			ctx,
			request.ID, request.Title, request.Content,
			[]*Tag{tag}, []*RequestTarget{target},
			group)
		require.NoError(t, err)
		exp := &RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			Title:     request.Title,
			Content:   request.Content,
			Comments:  request.Comments,
			Files:     request.Files,
			Tags:      []*Tag{tag},
			Targets:   []*RequestTargetDetail{expTarget},
			Statuses:  request.Statuses,
			Group:     group,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: request.CreatedBy,
		}
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(RequestTargetDetail{}, "ID", "PaidAt"))
		testutil.AssertEqual(t, exp, updatedRequest, opts...)
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
		target := &RequestTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		budget := random.Numeric(t, 10000)
		group, err := repo2.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)
		request, err := repo2.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{tag}, []*RequestTarget{target},
			group, user.ID)
		require.NoError(t, err)
		// CreatedAt の差を1秒以内に収めるためにここで time.Now を取る
		expTarget := &RequestTargetDetail{
			Target:    target.Target,
			Amount:    target.Amount,
			CreatedAt: time.Now(),
		}

		title := random.AlphaNumeric(t, 40)
		updatedRequest, err := repo2.UpdateRequest(
			ctx,
			request.ID, title, request.Content,
			[]*Tag{tag}, []*RequestTarget{target},
			group)
		require.NoError(t, err)
		exp := &RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			Title:     title,
			Content:   request.Content,
			Comments:  request.Comments,
			Files:     request.Files,
			Tags:      []*Tag{tag},
			Targets:   []*RequestTargetDetail{expTarget},
			Statuses:  request.Statuses,
			Group:     group,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: request.CreatedBy,
		}
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(RequestTargetDetail{}, "ID", "PaidAt"))
		testutil.AssertEqual(t, exp, updatedRequest, opts...)
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
		target := &RequestTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		budget := random.Numeric(t, 10000)
		group, err := repo3.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)
		request, err := repo3.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{tag}, []*RequestTarget{target},
			group, user.ID)
		require.NoError(t, err)
		// CreatedAt の差を1秒以内に収めるためにここで time.Now を取る
		expTarget := &RequestTargetDetail{
			Target:    target.Target,
			Amount:    target.Amount,
			CreatedAt: time.Now(),
		}
		content := random.AlphaNumeric(t, 100)
		updatedRequest, err := repo3.UpdateRequest(
			ctx,
			request.ID, request.Title, content,
			[]*Tag{tag}, []*RequestTarget{target},
			group)
		require.NoError(t, err)
		exp := &RequestDetail{
			ID:        request.ID,
			Status:    request.Status,
			Title:     request.Title,
			Content:   content,
			Comments:  request.Comments,
			Files:     request.Files,
			Tags:      []*Tag{tag},
			Targets:   []*RequestTargetDetail{expTarget},
			Statuses:  request.Statuses,
			Group:     group,
			CreatedAt: request.CreatedAt,
			UpdatedAt: time.Now(),
			CreatedBy: request.CreatedBy,
		}
		opts := testutil.ApproxEqualOptions()
		opts = append(opts,
			cmpopts.IgnoreFields(RequestTargetDetail{}, "ID", "PaidAt"))
		testutil.AssertEqual(t, exp, updatedRequest, opts...)
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
		target := &RequestTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		budget := random.Numeric(t, 10000)
		group, err := repo4.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)
		request, err := repo4.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{tag}, []*RequestTarget{target},
			group, user.ID)
		require.NoError(t, err)

		date := time.Now()
		unknownTag := &Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		_, err = repo4.UpdateRequest(
			ctx,
			request.ID, request.Title, request.Content,
			[]*Tag{unknownTag}, []*RequestTarget{target},
			group)
		require.Error(t, err)
	})

	t.Run("UnknownGroup", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		user, err := repo5.CreateUser(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			true)
		require.NoError(t, err)
		tag, err := repo5.CreateTag(ctx, random.AlphaNumeric(t, 20))
		require.NoError(t, err)
		target := &RequestTarget{
			Target: user.ID,
			Amount: random.Numeric(t, 10000),
		}
		budget := random.Numeric(t, 10000)
		group, err := repo5.CreateGroup(
			ctx,
			random.AlphaNumeric(t, 20),
			random.AlphaNumeric(t, 30),
			&budget)
		require.NoError(t, err)
		request, err := repo5.CreateRequest(
			ctx,
			random.AlphaNumeric(t, 40),
			random.AlphaNumeric(t, 100),
			[]*Tag{tag}, []*RequestTarget{target},
			group, user.ID)
		require.NoError(t, err)

		date := time.Now()
		unknownBudget := random.Numeric(t, 100000)
		unknownGroup := &Group{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 30),
			Budget:      &unknownBudget,
			CreatedAt:   date,
			UpdatedAt:   date,
		}
		_, err = repo5.UpdateRequest(
			ctx,
			request.ID, request.Title, request.Content,
			[]*Tag{tag}, []*RequestTarget{target},
			unknownGroup)
		require.Error(t, err)
	})
}
