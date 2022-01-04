package router

import (
	"encoding/json"
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
		assert.NoError(t, err)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h, err := NewTestHandlers(t, ctrl)
		assert.NoError(t, err)
		h.Repository.MockTransactionRepository.
			EXPECT().
			GetTransactions(c.Request().Context(), model.TransactionQuery{
				Sort:   nil,
				Target: nil,
				Year:   nil,
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

}
