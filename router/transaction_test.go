package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
	"github.com/traPtitech/Jomon/testutil/random"
	"go.uber.org/mock/gomock"
)

func TestHandlers_GetTransactions(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
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
		req, err := http.NewRequest(http.MethodGet, "/api/transactions", nil)
		require.NoError(t, err)
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

		res := lo.Map(txs, func(tx *model.TransactionResponse, _ int) *TransactionNewCreate {
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
			return &TransactionNewCreate{
				ID:        tx.ID,
				Amount:    tx.Amount,
				Target:    tx.Target,
				Tags:      tag,
				Group:     group,
				CreatedAt: tx.CreatedAt,
				UpdatedAt: tx.UpdatedAt,
			}
		})

		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetTransactions(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessWithSort", func(t *testing.T) {
		t.Parallel()
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
		req, err := http.NewRequest(http.MethodGet, "/api/transactions?sort=created_at", nil)
		require.NoError(t, err)
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

		res := lo.Map(txs, func(tx *model.TransactionResponse, _ int) *TransactionNewCreate {
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
			return &TransactionNewCreate{
				ID:        tx.ID,
				Amount:    tx.Amount,
				Target:    tx.Target,
				Tags:      tag,
				Group:     group,
				CreatedAt: tx.CreatedAt,
				UpdatedAt: tx.UpdatedAt,
			}
		})

		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetTransactions(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessWithAscSort", func(t *testing.T) {
		t.Parallel()
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
		req, err := http.NewRequest(http.MethodGet, "/api/transactions?sort=-created_at", nil)
		require.NoError(t, err)
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

		res := lo.Map(txs, func(tx *model.TransactionResponse, _ int) *TransactionNewCreate {
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
			return &TransactionNewCreate{
				ID:        tx.ID,
				Amount:    tx.Amount,
				Target:    tx.Target,
				Tags:      tag,
				Group:     group,
				CreatedAt: tx.CreatedAt,
				UpdatedAt: tx.UpdatedAt,
			}
		})

		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetTransactions(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessWithTarget", func(t *testing.T) {
		t.Parallel()
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
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/transactions?target=%s", target1),
			nil)
		require.NoError(t, err)
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

		res := lo.Map(txs, func(tx *model.TransactionResponse, _ int) *TransactionNewCreate {
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
			return &TransactionNewCreate{
				ID:        tx.ID,
				Amount:    tx.Amount,
				Target:    tx.Target,
				Tags:      tag,
				Group:     group,
				CreatedAt: tx.CreatedAt,
				UpdatedAt: tx.UpdatedAt,
			}
		})

		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetTransactions(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessWithSinceUntil", func(t *testing.T) {
		t.Parallel()
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
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/api/transactions?since=%s&until=%s", "2020-01-01", "2020-01-02"),
			nil)
		require.NoError(t, err)
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

		res := lo.Map(txs, func(tx *model.TransactionResponse, _ int) *TransactionNewCreate {
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
			return &TransactionNewCreate{
				ID:        tx.ID,
				Amount:    tx.Amount,
				Target:    tx.Target,
				Tags:      tag,
				Group:     group,
				CreatedAt: tx.CreatedAt,
				UpdatedAt: tx.UpdatedAt,
			}
		})

		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetTransactions(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	// TODO: SuccessWithLimit, SuccessWithOffset
}

func TestHandlers_PostTransaction(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
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
		reqBody := fmt.Sprintf(
			`{"amount": %d, "targets": ["%s"], "tags": ["%s"], "group": "%s"}`,
			tx1.Amount, tx1.Target, tag.ID, group)
		req, err := http.NewRequest(
			http.MethodPost,
			"/api/transactions",
			strings.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		h.Repository.MockTransactionRepository.
			EXPECT().
			CreateTransaction(c.Request().Context(), tx1.Amount, tx1.Target, tags, &group, nil).
			Return(tx1, nil)

		res := lo.Map(txs, func(tx *model.TransactionResponse, _ int) *TransactionNewCreate {
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
			return &TransactionNewCreate{
				ID:        tx.ID,
				Amount:    tx.Amount,
				Target:    tx.Target,
				Tags:      tag,
				Group:     group,
				CreatedAt: tx.CreatedAt,
				UpdatedAt: tx.UpdatedAt,
			}
		})

		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PostTransaction(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})

	t.Run("SuccessWithRequest", func(t *testing.T) {
		t.Parallel()
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
		reqBody := fmt.Sprintf(
			`{"amount": %d, "targets": ["%s"], "tags": ["%s"], "group": "%s", "request": "%s"}`,
			tx.Amount, tx.Target, tag.ID, group, request.ID)
		req, err := http.NewRequest(
			http.MethodPost,
			"/api/transactions",
			strings.NewReader(reqBody))
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		h.Repository.MockTransactionRepository.
			EXPECT().
			CreateTransaction(
				c.Request().Context(),
				tx.Amount, tx.Target,
				tags, &group, &request.ID).
			Return(tx, nil)

		to := lo.Map(tx.Tags, func(modelTag *model.Tag, _ int) *TagOverview {
			return &TagOverview{
				ID:        modelTag.ID,
				Name:      modelTag.Name,
				CreatedAt: modelTag.CreatedAt,
				UpdatedAt: modelTag.UpdatedAt,
			}
		})

		grov := &GroupOverview{
			ID:          tx.Group.ID,
			Name:        tx.Group.Name,
			Description: tx.Group.Description,
			Budget:      tx.Group.Budget,
			CreatedAt:   tx.Group.CreatedAt,
			UpdatedAt:   tx.Group.UpdatedAt,
		}
		res := []*TransactionNewCreate{
			{
				ID:        tx.ID,
				Amount:    tx.Amount,
				Target:    tx.Target,
				Tags:      to,
				Group:     grov,
				CreatedAt: tx.CreatedAt,
				UpdatedAt: tx.UpdatedAt,
			},
		}
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PostTransaction(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})
}

func TestHandlers_GetTransaction(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
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
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/transactions/%s", tx.ID), nil)
		require.NoError(t, err)
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

		var resOverview TransactionCorrestion
		to := lo.Map(tx.Tags, func(modelTag *model.Tag, _ int) *TagOverview {
			return &TagOverview{
				ID:        modelTag.ID,
				Name:      modelTag.Name,
				CreatedAt: modelTag.CreatedAt,
				UpdatedAt: modelTag.UpdatedAt,
			}
		})

		grov := &GroupOverview{
			ID:          tx.Group.ID,
			Name:        tx.Group.Name,
			Description: tx.Group.Description,
			Budget:      tx.Group.Budget,
			CreatedAt:   tx.Group.CreatedAt,
			UpdatedAt:   tx.Group.UpdatedAt,
		}
		resOverview = TransactionCorrestion{
			ID:        tx.ID,
			Amount:    tx.Amount,
			Target:    tx.Target,
			Tags:      to,
			Group:     grov,
			CreatedAt: tx.CreatedAt,
			UpdatedAt: tx.UpdatedAt,
		}
		res := resOverview
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetTransaction(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})
}

func TestHandlers_PutTransaction(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
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
		reqBody := fmt.Sprintf(
			`{"amount": %d, "target": "%s", "tags": ["%s"]}`,
			updated.Amount, updated.Target, updatedTag.ID)
		req, err := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("/api/transactions/%s", tx.ID),
			strings.NewReader(reqBody))
		require.NoError(t, err)
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
				tx.ID, updated.Amount, updated.Target,
				updatedTags, nil, nil).
			Return(updated, nil)

		var resOverview TransactionCorrestion
		to := lo.Map(updated.Tags, func(modelTag *model.Tag, _ int) *TagOverview {
			return &TagOverview{
				ID:        modelTag.ID,
				Name:      modelTag.Name,
				CreatedAt: modelTag.CreatedAt,
				UpdatedAt: modelTag.UpdatedAt,
			}
		})
		grov := &GroupOverview{
			ID:          updated.Group.ID,
			Name:        updated.Group.Name,
			Description: updated.Group.Description,
			Budget:      updated.Group.Budget,
			CreatedAt:   updated.Group.CreatedAt,
			UpdatedAt:   updated.Group.UpdatedAt,
		}
		resOverview = TransactionCorrestion{
			ID:        tx.ID,
			Amount:    updated.Amount,
			Target:    updated.Target,
			Tags:      to,
			Group:     grov,
			CreatedAt: updated.CreatedAt,
			UpdatedAt: updated.UpdatedAt,
		}
		res := resOverview
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.PutTransaction(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})
}
