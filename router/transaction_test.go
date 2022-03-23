package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
	"github.com/traPtitech/Jomon/testutil/random"
)

func TestHandlers_GetTransactions(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)

		date := time.Now()

		budget := random.Numeric(t, 1000000)

		tag := &model.Tag{
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
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
				Tag:    nil,
				Group:  nil,
			}).
			Return(txs, nil)

		resOverview := []*Transaction{}
		for _, tx := range txs {
			tag := []*TagOverview{}
			for _, modelTag := range tx.Tags {
				tag = append(tag, &TagOverview{
					ID:          modelTag.ID,
					Name:        modelTag.Name,
					Description: modelTag.Description,
					CreatedAt:   modelTag.CreatedAt,
					UpdatedAt:   modelTag.UpdatedAt,
				})
			}
			group := &GroupOverview{
				ID:          tx.Group.ID,
				Name:        tx.Group.Name,
				Description: tx.Group.Description,
				Budget:      tx.Group.Budget,
				CreatedAt:   tx.Group.CreatedAt,
				UpdatedAt:   tx.Group.UpdatedAt,
			}
			resOverview = append(resOverview, &Transaction{
				ID:        tx.ID,
				Amount:    tx.Amount,
				Target:    tx.Target,
				Tags:      tag,
				Group:     group,
				CreatedAt: tx.CreatedAt,
				UpdatedAt: tx.UpdatedAt,
			})
		}
		res := resOverview
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
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
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
		c.SetParamNames("sort")
		c.SetParamValues("created_at")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		sortQuery := "created_at"
		h.Repository.MockTransactionRepository.
			EXPECT().
			GetTransactions(c.Request().Context(), model.TransactionQuery{
				Sort: &sortQuery,
			}).
			Return(txs, nil)

		resOverview := []*Transaction{}
		for _, tx := range txs {
			tag := []*TagOverview{}
			for _, modelTag := range tx.Tags {
				tag = append(tag, &TagOverview{
					ID:          modelTag.ID,
					Name:        modelTag.Name,
					Description: modelTag.Description,
					CreatedAt:   modelTag.CreatedAt,
					UpdatedAt:   modelTag.UpdatedAt,
				})
			}
			group := &GroupOverview{
				ID:          tx.Group.ID,
				Name:        tx.Group.Name,
				Description: tx.Group.Description,
				Budget:      tx.Group.Budget,
				CreatedAt:   tx.Group.CreatedAt,
				UpdatedAt:   tx.Group.UpdatedAt,
			}
			resOverview = append(resOverview, &Transaction{
				ID:        tx.ID,
				Amount:    tx.Amount,
				Target:    tx.Target,
				Tags:      tag,
				Group:     group,
				CreatedAt: tx.CreatedAt,
				UpdatedAt: tx.UpdatedAt,
			})
		}
		res := resOverview
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
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
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
		c.SetParamNames("sort")
		c.SetParamValues("-created_at")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		sortQuery := "-created_at"
		h.Repository.MockTransactionRepository.
			EXPECT().
			GetTransactions(c.Request().Context(), model.TransactionQuery{
				Sort: &sortQuery,
			}).
			Return(txs, nil)

		resOverview := []*Transaction{}
		for _, tx := range txs {
			tag := []*TagOverview{}
			for _, modelTag := range tx.Tags {
				tag = append(tag, &TagOverview{
					ID:          modelTag.ID,
					Name:        modelTag.Name,
					Description: modelTag.Description,
					CreatedAt:   modelTag.CreatedAt,
					UpdatedAt:   modelTag.UpdatedAt,
				})
			}
			group := &GroupOverview{
				ID:          tx.Group.ID,
				Name:        tx.Group.Name,
				Description: tx.Group.Description,
				Budget:      tx.Group.Budget,
				CreatedAt:   tx.Group.CreatedAt,
				UpdatedAt:   tx.Group.UpdatedAt,
			}
			resOverview = append(resOverview, &Transaction{
				ID:        tx.ID,
				Amount:    tx.Amount,
				Target:    tx.Target,
				Tags:      tag,
				Group:     group,
				CreatedAt: tx.CreatedAt,
				UpdatedAt: tx.UpdatedAt,
			})
		}
		res := resOverview
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
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
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
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/transactions?target=%s", target1), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("target")
		c.SetParamValues(target1)

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		targetQuery := target1
		h.Repository.MockTransactionRepository.
			EXPECT().
			GetTransactions(c.Request().Context(), model.TransactionQuery{
				Target: &targetQuery,
			}).
			Return(txs, nil)

		resOverview := []*Transaction{}
		for _, tx := range txs {
			tag := []*TagOverview{}
			for _, modelTag := range tx.Tags {
				tag = append(tag, &TagOverview{
					ID:          modelTag.ID,
					Name:        modelTag.Name,
					Description: modelTag.Description,
					CreatedAt:   modelTag.CreatedAt,
					UpdatedAt:   modelTag.UpdatedAt,
				})
			}
			group := &GroupOverview{
				ID:          tx.Group.ID,
				Name:        tx.Group.Name,
				Description: tx.Group.Description,
				Budget:      tx.Group.Budget,
				CreatedAt:   tx.Group.CreatedAt,
				UpdatedAt:   tx.Group.UpdatedAt,
			}
			resOverview = append(resOverview, &Transaction{
				ID:        tx.ID,
				Amount:    tx.Amount,
				Target:    tx.Target,
				Tags:      tag,
				Group:     group,
				CreatedAt: tx.CreatedAt,
				UpdatedAt: tx.UpdatedAt,
			})
		}
		res := resOverview
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
			ID:          uuid.New(),
			Name:        random.AlphaNumeric(t, 20),
			Description: random.AlphaNumeric(t, 50),
			CreatedAt:   date,
			UpdatedAt:   date,
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
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/transactions?since=%s&until=%s", "2020-01-01", "2020-01-02"), nil)
		require.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("since", "until")
		c.SetParamValues("2020-01-01", "2020-01-02")

		h, err := NewTestHandlers(t, ctrl)
		require.NoError(t, err)

		h.Repository.MockTransactionRepository.
			EXPECT().
			GetTransactions(c.Request().Context(), model.TransactionQuery{
				Since: &since,
				Until: &until,
			}).
			Return(txs, nil)

		resOverview := []*Transaction{}
		for _, tx := range txs {
			tag := []*TagOverview{}
			for _, modelTag := range tx.Tags {
				tag = append(tag, &TagOverview{
					ID:          modelTag.ID,
					Name:        modelTag.Name,
					Description: modelTag.Description,
					CreatedAt:   modelTag.CreatedAt,
					UpdatedAt:   modelTag.UpdatedAt,
				})
			}
			group := &GroupOverview{
				ID:          tx.Group.ID,
				Name:        tx.Group.Name,
				Description: tx.Group.Description,
				Budget:      tx.Group.Budget,
				CreatedAt:   tx.Group.CreatedAt,
				UpdatedAt:   tx.Group.UpdatedAt,
			}
			resOverview = append(resOverview, &Transaction{
				ID:        tx.ID,
				Amount:    tx.Amount,
				Target:    tx.Target,
				Tags:      tag,
				Group:     group,
				CreatedAt: tx.CreatedAt,
				UpdatedAt: tx.UpdatedAt,
			})
		}
		res := resOverview
		resBody, err := json.Marshal(res)
		require.NoError(t, err)

		if assert.NoError(t, h.Handlers.GetTransactions(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, string(resBody), strings.TrimRight(rec.Body.String(), "\n"))
		}
	})
}
