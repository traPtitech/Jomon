package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
	"github.com/traPtitech/Jomon/testutil"
	"github.com/traPtitech/Jomon/testutil/random"
	"go.uber.org/mock/gomock"
)

// FIXME: 同様の処理がtransaction.goにもある
func modelTransactionResponseToTransaction(tx *model.TransactionResponse) *TransactionResponse {
	tag := lo.Map(tx.Tags, func(modelTag *model.Tag, _ int) *TagOverview {
		return &TagOverview{
			ID:        modelTag.ID,
			Name:      modelTag.Name,
			CreatedAt: modelTag.CreatedAt,
			UpdatedAt: modelTag.UpdatedAt,
		}
	})

	group := &GroupOverview{
		ID:          tx.Group.ID,
		Name:        tx.Group.Name,
		Description: tx.Group.Description,
		Budget:      tx.Group.Budget,
		CreatedAt:   tx.Group.CreatedAt,
		UpdatedAt:   tx.Group.UpdatedAt,
	}
	return &TransactionResponse{
		ID:        tx.ID,
		Title:     tx.Title,
		Amount:    tx.Amount,
		Target:    tx.Target,
		Tags:      tag,
		Group:     group,
		CreatedAt: tx.CreatedAt,
		UpdatedAt: tx.UpdatedAt,
	}
}

func TestHandlers_GetTransactions(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		budget := random.Numeric(t, 1000000)
		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		tx1 := &model.TransactionResponse{
			ID:     uuid.New(),
			Title:  random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 1000000),
			Target: random.AlphaNumeric(t, 20),
			Tags: []*model.Tag{
				tag,
			},
			Group: &model.Group{
				ID:          uuid.New(),
				Name:        random.AlphaNumeric(t, 20),
				Description: random.AlphaNumeric(t, 50),
				Budget:      &budget,
				CreatedAt:   date,
				UpdatedAt:   date,
			},
			CreatedAt: date,
			UpdatedAt: date,
		}
		tx2 := &model.TransactionResponse{
			ID:     uuid.New(),
			Title:  random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 1000000),
			Target: random.AlphaNumeric(t, 20),
			Tags: []*model.Tag{
				tag,
			},
			Group: &model.Group{
				ID:          uuid.New(),
				Name:        random.AlphaNumeric(t, 20),
				Description: random.AlphaNumeric(t, 50),
				Budget:      &budget,
				CreatedAt:   date,
				UpdatedAt:   date,
			},
			CreatedAt: date,
			UpdatedAt: date,
		}
		txs := []*model.TransactionResponse{tx1, tx2}

		e := echo.New()
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/transactions", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockTransactionRepository.
			EXPECT().
			GetTransactions(c.Request().Context(), model.TransactionQuery{
				Sort:   nil,
				Target: nil,
				Since:  nil,
				Until:  nil,
				Limit:  100,
				Offset: 0,
				Tag:    nil,
				Group:  nil,
			}).
			Return(txs, nil)

		require.NoError(t, h.Handlers.GetTransactions(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*TransactionResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := lo.Map(txs, func(tx *model.TransactionResponse, _ int) *TransactionResponse {
			return modelTransactionResponseToTransaction(tx)
		})
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithSort", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		budget := random.Numeric(t, 1000000)
		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		tx1 := &model.TransactionResponse{
			ID:     uuid.New(),
			Title:  random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 1000000),
			Target: random.AlphaNumeric(t, 20),
			Tags: []*model.Tag{
				tag,
			},
			Group: &model.Group{
				ID:          uuid.New(),
				Name:        random.AlphaNumeric(t, 20),
				Description: random.AlphaNumeric(t, 50),
				Budget:      &budget,
				CreatedAt:   date,
				UpdatedAt:   date,
			},
			CreatedAt: date,
			UpdatedAt: date,
		}
		tx2 := &model.TransactionResponse{
			ID:     uuid.New(),
			Title:  random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 1000000),
			Target: random.AlphaNumeric(t, 20),
			Tags: []*model.Tag{
				tag,
			},
			Group: &model.Group{
				ID:          uuid.New(),
				Name:        random.AlphaNumeric(t, 20),
				Description: random.AlphaNumeric(t, 50),
				Budget:      &budget,
				CreatedAt:   date.Add(time.Hour),
				UpdatedAt:   date.Add(time.Hour),
			},

			CreatedAt: date.Add(time.Hour),
			UpdatedAt: date.Add(time.Hour),
		}
		txs := []*model.TransactionResponse{tx1, tx2}

		e := echo.New()
		path := "/api/transactions?sort=created_at"
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		sortQuery := "created_at"
		h.Repository.MockTransactionRepository.
			EXPECT().
			GetTransactions(c.Request().Context(), model.TransactionQuery{
				Sort:   &sortQuery,
				Limit:  100,
				Offset: 0,
			}).
			Return(txs, nil)

		require.NoError(t, h.Handlers.GetTransactions(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*TransactionResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := lo.Map(txs, func(tx *model.TransactionResponse, _ int) *TransactionResponse {
			return modelTransactionResponseToTransaction(tx)
		})
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithAscSort", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		budget := random.Numeric(t, 1000000)
		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		tx1 := &model.TransactionResponse{
			ID:     uuid.New(),
			Title:  random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 1000000),
			Target: random.AlphaNumeric(t, 20),
			Tags: []*model.Tag{
				tag,
			},
			Group: &model.Group{
				ID:          uuid.New(),
				Name:        random.AlphaNumeric(t, 20),
				Description: random.AlphaNumeric(t, 50),
				Budget:      &budget,
				CreatedAt:   date,
				UpdatedAt:   date,
			},
			CreatedAt: date,
			UpdatedAt: date,
		}
		tx2 := &model.TransactionResponse{
			ID:     uuid.New(),
			Title:  random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 1000000),
			Target: random.AlphaNumeric(t, 20),
			Tags: []*model.Tag{
				tag,
			},
			Group: &model.Group{
				ID:          uuid.New(),
				Name:        random.AlphaNumeric(t, 20),
				Description: random.AlphaNumeric(t, 50),
				Budget:      &budget,
				CreatedAt:   date.Add(time.Hour),
				UpdatedAt:   date.Add(time.Hour),
			},

			CreatedAt: date.Add(time.Hour),
			UpdatedAt: date.Add(time.Hour),
		}
		// Reverse
		txs := []*model.TransactionResponse{tx2, tx1}

		e := echo.New()
		path := "/api/transactions?sort=-created_at"
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		sortQuery := "-created_at"
		h.Repository.MockTransactionRepository.
			EXPECT().
			GetTransactions(c.Request().Context(), model.TransactionQuery{
				Sort:   &sortQuery,
				Limit:  100,
				Offset: 0,
			}).
			Return(txs, nil)

		require.NoError(t, h.Handlers.GetTransactions(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*TransactionResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := lo.Map(txs, func(tx *model.TransactionResponse, _ int) *TransactionResponse {
			return modelTransactionResponseToTransaction(tx)
		})
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithTarget", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		budget := random.Numeric(t, 1000000)
		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		target1 := random.AlphaNumeric(t, 20)
		tx1 := &model.TransactionResponse{
			ID:     uuid.New(),
			Title:  random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 1000000),
			Target: target1,
			Tags: []*model.Tag{
				tag,
			},
			Group: &model.Group{
				ID:          uuid.New(),
				Name:        random.AlphaNumeric(t, 20),
				Description: random.AlphaNumeric(t, 50),
				Budget:      &budget,
				CreatedAt:   date,
				UpdatedAt:   date,
			},
			CreatedAt: date,
			UpdatedAt: date,
		}
		tx2 := &model.TransactionResponse{
			ID:     uuid.New(),
			Title:  random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 1000000),
			Target: random.AlphaNumeric(t, 20),
			Tags: []*model.Tag{
				tag,
			},
			Group: &model.Group{
				ID:          uuid.New(),
				Name:        random.AlphaNumeric(t, 20),
				Description: random.AlphaNumeric(t, 50),
				Budget:      &budget,
				CreatedAt:   date.Add(time.Hour),
				UpdatedAt:   date.Add(time.Hour),
			},

			CreatedAt: date.Add(time.Hour),
			UpdatedAt: date.Add(time.Hour),
		}
		txs := []*model.TransactionResponse{tx1, tx2}

		e := echo.New()
		path := fmt.Sprintf("/api/transactions?target=%s", target1)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		targetQuery := target1
		h.Repository.MockTransactionRepository.
			EXPECT().
			GetTransactions(c.Request().Context(), model.TransactionQuery{
				Target: &targetQuery,
				Limit:  100,
				Offset: 0,
			}).
			Return(txs, nil)

		require.NoError(t, h.Handlers.GetTransactions(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*TransactionResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := lo.Map(txs, func(tx *model.TransactionResponse, _ int) *TransactionResponse {
			return modelTransactionResponseToTransaction(tx)
		})
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithSinceUntil", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		budget := random.Numeric(t, 1000000)
		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		target1 := random.AlphaNumeric(t, 20)
		tx1 := &model.TransactionResponse{
			ID:     uuid.New(),
			Title:  random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 1000000),
			Target: target1,
			Tags: []*model.Tag{
				tag,
			},
			Group: &model.Group{
				ID:          uuid.New(),
				Name:        random.AlphaNumeric(t, 20),
				Description: random.AlphaNumeric(t, 50),
				Budget:      &budget,
				CreatedAt:   date,
				UpdatedAt:   date,
			},
			CreatedAt: date,
			UpdatedAt: date,
		}
		txs := []*model.TransactionResponse{tx1}
		since, err := service.StrToDate("2020-01-01")
		require.NoError(t, err)
		until, err := service.StrToDate("2020-01-02")
		require.NoError(t, err)

		e := echo.New()
		path := fmt.Sprintf("/api/transactions?since=%s&until=%s", "2020-01-01", "2020-01-02")
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)
		h.Repository.MockTransactionRepository.
			EXPECT().
			GetTransactions(c.Request().Context(), model.TransactionQuery{
				Since:  &since,
				Until:  &until,
				Limit:  100,
				Offset: 0,
			}).
			Return(txs, nil)

		require.NoError(t, h.Handlers.GetTransactions(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*TransactionResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := lo.Map(txs, func(tx *model.TransactionResponse, _ int) *TransactionResponse {
			return modelTransactionResponseToTransaction(tx)
		})
		testutil.RequireEqual(t, exp, got, opts...)
	})

	// TODO: SuccessWithLimit, SuccessWithOffset
}

func TestHandlers_PostTransaction(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		budget := random.Numeric(t, 1000000)
		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		target1 := random.AlphaNumeric(t, 20)
		tx1 := &model.TransactionResponse{
			ID:     uuid.New(),
			Title:  random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 1000000),
			Target: target1,
			Tags: []*model.Tag{
				tag,
			},
			Group: &model.Group{
				ID:          uuid.New(),
				Name:        random.AlphaNumeric(t, 20),
				Description: random.AlphaNumeric(t, 50),
				Budget:      &budget,
				CreatedAt:   date,
				UpdatedAt:   date,
			},
			CreatedAt: date,
			UpdatedAt: date,
		}
		txs := []*model.TransactionResponse{tx1}
		tags := []*uuid.UUID{&tag.ID}
		group := tx1.Group.ID

		e := echo.New()
		reqBody, err := json.Marshal(&struct {
			Title   string      `json:"title"`
			Amount  int         `json:"amount"`
			Targets []string    `json:"targets"`
			Tags    []uuid.UUID `json:"tags"`
			Group   *uuid.UUID  `json:"group"`
		}{
			Title:   tx1.Title,
			Amount:  tx1.Amount,
			Targets: []string{tx1.Target},
			Tags:    []uuid.UUID{tag.ID},
			Group:   &group,
		})
		require.NoError(t, err)
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/transactions", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		h.Repository.MockTransactionRepository.
			EXPECT().
			CreateTransaction(
				c.Request().Context(),
				tx1.Title, tx1.Amount, tx1.Target, tags, &group, nil).
			Return(tx1, nil)

		require.NoError(t, h.Handlers.PostTransaction(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*TransactionResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := lo.Map(txs, func(tx *model.TransactionResponse, _ int) *TransactionResponse {
			return modelTransactionResponseToTransaction(tx)
		})
		testutil.RequireEqual(t, exp, got, opts...)
	})

	t.Run("SuccessWithRequest", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		date := time.Now()
		budget := random.Numeric(t, 1000000)
		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: date,
			UpdatedAt: date,
		}
		target1 := random.AlphaNumeric(t, 20)
		tx := &model.TransactionResponse{
			ID:     uuid.New(),
			Title:  random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 1000000),
			Target: target1,
			Tags: []*model.Tag{
				tag,
			},
			Group: &model.Group{
				ID:          uuid.New(),
				Name:        random.AlphaNumeric(t, 20),
				Description: random.AlphaNumeric(t, 50),
				Budget:      &budget,
				CreatedAt:   date,
				UpdatedAt:   date,
			},
			CreatedAt: date,
			UpdatedAt: date,
		}
		tags := []*uuid.UUID{&tag.ID}
		group := tx.Group.ID
		request := &model.RequestDetail{
			ID:        uuid.New(),
			Status:    model.Accepted,
			Title:     random.AlphaNumeric(t, 20),
			Content:   random.AlphaNumeric(t, 50),
			Comments:  []*model.Comment{},
			Files:     []*uuid.UUID{},
			Statuses:  []*model.RequestStatus{},
			Tags:      []*model.Tag{},
			Group:     nil,
			CreatedAt: date,
			UpdatedAt: date,
			CreatedBy: uuid.New(),
		}

		e := echo.New()
		reqBody, err := json.Marshal(&PostTransactionsRequest{
			Title:   tx.Title,
			Amount:  tx.Amount,
			Targets: []*string{&tx.Target},
			Tags:    tags,
			Group:   &group,
			Request: &request.ID,
		})
		require.NoError(t, err)
		req := httptest.NewRequestWithContext(
			ctx, http.MethodPost, "/api/transactions", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		h.Repository.MockTransactionRepository.
			EXPECT().
			CreateTransaction(
				c.Request().Context(),
				tx.Title, tx.Amount, tx.Target,
				tags, &group, &request.ID).
			Return(tx, nil)

		require.NoError(t, h.Handlers.PostTransaction(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got []*TransactionResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := lo.Map(
			[]*model.TransactionResponse{tx},
			func(tx *model.TransactionResponse, _ int) *TransactionResponse {
				return modelTransactionResponseToTransaction(tx)
			})
		testutil.RequireEqual(t, exp, got, opts...)
	})

	// TODO: FailWithoutTitle
	// PostTransactionにvalidationが入ってから
}

func TestHandlers_GetTransaction(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		budget := random.Numeric(t, 1000000)
		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		tx := &model.TransactionResponse{
			ID:     uuid.New(),
			Title:  random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 1000000),
			Target: random.AlphaNumeric(t, 20),
			Tags: []*model.Tag{
				tag,
			},
			Group: &model.Group{
				ID:          uuid.New(),
				Name:        random.AlphaNumeric(t, 20),
				Description: random.AlphaNumeric(t, 50),
				Budget:      &budget,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		e := echo.New()
		path := fmt.Sprintf("/api/transactions/%s", tx.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("transactionID")
		c.SetParamValues(tx.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		h.Repository.MockTransactionRepository.
			EXPECT().
			GetTransaction(c.Request().Context(), tx.ID).
			Return(tx, nil)

		require.NoError(t, h.Handlers.GetTransaction(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *TransactionResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelTransactionResponseToTransaction(tx)
		testutil.RequireEqual(t, exp, got, opts...)
	})
}

func TestHandlers_PutTransaction(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.NewContext(t)
		ctrl := gomock.NewController(t)

		budget := random.Numeric(t, 1000000)
		tag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		updatedTag := &model.Tag{
			ID:        uuid.New(),
			Name:      random.AlphaNumeric(t, 20),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		tx := &model.TransactionResponse{
			ID:     uuid.New(),
			Title:  random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 1000000),
			Target: random.AlphaNumeric(t, 20),
			Tags: []*model.Tag{
				tag,
			},
			Group: &model.Group{
				ID:          uuid.New(),
				Name:        random.AlphaNumeric(t, 20),
				Description: random.AlphaNumeric(t, 50),
				Budget:      &budget,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		updated := &model.TransactionResponse{
			ID:     tx.ID,
			Title:  random.AlphaNumeric(t, 20),
			Amount: random.Numeric(t, 1000000),
			Target: random.AlphaNumeric(t, 20),
			Tags: []*model.Tag{
				updatedTag,
			},
			Group: &model.Group{
				ID:          uuid.New(),
				Name:        random.AlphaNumeric(t, 20),
				Description: random.AlphaNumeric(t, 50),
				Budget:      &budget,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		updatedTags := lo.Map(updated.Tags, func(modelTag *model.Tag, _ int) *uuid.UUID {
			return &modelTag.ID
		})

		e := echo.New()
		reqBody, err := json.Marshal(&struct {
			Title  string      `json:"title"`
			Amount int         `json:"amount"`
			Target string      `json:"target"`
			Tags   []uuid.UUID `json:"tags"`
		}{
			Title:  updated.Title,
			Amount: updated.Amount,
			Target: updated.Target,
			Tags:   []uuid.UUID{updatedTag.ID},
		})
		require.NoError(t, err)
		path := fmt.Sprintf("/api/transactions/%s", tx.ID)
		req := httptest.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("transactionID")
		c.SetParamValues(tx.ID.String())

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		h.Repository.MockTransactionRepository.
			EXPECT().
			UpdateTransaction(
				c.Request().Context(),
				tx.ID, updated.Title, updated.Amount, updated.Target,
				updatedTags, nil, nil).
			Return(updated, nil)
		require.NoError(t, h.Handlers.PutTransaction(c))
		require.Equal(t, http.StatusOK, rec.Code)
		var got *TransactionResponse
		err = json.Unmarshal(rec.Body.Bytes(), &got)
		require.NoError(t, err)
		opts := testutil.ApproxEqualOptions()
		exp := modelTransactionResponseToTransaction(updated)
		testutil.RequireEqual(t, exp, got, opts...)
	})
}
