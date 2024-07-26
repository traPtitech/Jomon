package router

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
	"go.uber.org/zap"
)

type Transaction struct {
	ID        uuid.UUID      `json:"id"`
	Amount    int            `json:"amount"`
	Target    string         `json:"target"`
	Request   *uuid.UUID     `json:"request"`
	Tags      []*TagOverview `json:"tags"`
	Group     *GroupOverview `json:"group"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type TransactionOverview struct {
	Amount  int          `json:"amount"`
	Targets []*string    `json:"targets"`
	Tags    []*uuid.UUID `json:"tags"`
	Group   *uuid.UUID   `json:"group"`
	Request *uuid.UUID   `json:"request"`
}

type TransactionOverviewWithOneTarget struct {
	Amount  int          `json:"amount"`
	Target  string       `json:"target"`
	Tags    []*uuid.UUID `json:"tags"`
	Group   *uuid.UUID   `json:"group"`
	Request *uuid.UUID   `json:"request"`
}

func (h Handlers) GetTransactions(c echo.Context) error {
	ctx := c.Request().Context()
	var sort *string
	if c.QueryParam("sort") != "" {
		s := c.QueryParam("sort")
		sort = &s
	}
	var target *string
	if c.QueryParam("target") != "" {
		t := c.QueryParam("target")
		target = &t
	}
	var since *time.Time
	if c.QueryParam("since") != "" {
		var err error
		s, err := service.StrToDate(c.QueryParam("since"))
		if err != nil {
			h.Logger.Info("could not parse since as time.Time", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		since = &s
	}
	var until *time.Time
	if c.QueryParam("until") != "" {
		var err error
		u, err := service.StrToDate(c.QueryParam("until"))
		if err != nil {
			h.Logger.Info("could not parse until as time.Time", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		until = &u
	}
	limit := 100
	if limitQuery := c.QueryParam("limit"); limitQuery != "" {
		limitI, err := strconv.Atoi(limitQuery)
		if err != nil {
			h.Logger.Info("could not parse limit as integer", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		if limitI < 0 {
			h.Logger.Info("received negative limit", zap.Int("limit", limitI))
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("negative limit(=%d) is invalid", limitI))
		}
		limit = limitI
	}
	offset := 0
	if offsetQuery := c.QueryParam("offset"); offsetQuery != "" {
		offsetI, err := strconv.Atoi(offsetQuery)
		if err != nil {
			h.Logger.Info("could not parse limit as integer", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		if offsetI < 0 {
			h.Logger.Info("received negative offset", zap.Int("offset", offsetI))
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("negative offset(=%d) is invalid", offsetI))
		}
		offset = offsetI
	}
	var tag *string
	if c.QueryParam("tag") != "" {
		t := c.QueryParam("tag")
		tag = &t
	}
	var group *string
	if c.QueryParam("group") != "" {
		g := c.QueryParam("group")
		group = &g
	}
	var request *uuid.UUID
	if c.QueryParam("request") != "" {
		var r uuid.UUID
		r, err := uuid.Parse(c.QueryParam("request"))
		if err != nil {
			h.Logger.Info("could not parse request as uuid.UUID", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		request = &r
	}
	query := model.TransactionQuery{
		Sort:    sort,
		Target:  target,
		Since:   since,
		Until:   until,
		Limit:   limit,
		Offset:  offset,
		Tag:     tag,
		Group:   group,
		Request: request,
	}
	txs, err := h.Repository.GetTransactions(ctx, query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := lo.Map(txs, func(tx *model.TransactionResponse, index int) *Transaction {
		tags := lo.Map(tx.Tags, func(tag *model.Tag, index int) *TagOverview {
			return &TagOverview{
				ID:        tag.ID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt,
				UpdatedAt: tag.UpdatedAt,
			}
		})

		var group *GroupOverview
		if tx.Group != nil {
			group = &GroupOverview{
				ID:          tx.Group.ID,
				Name:        tx.Group.Name,
				Description: tx.Group.Description,
				Budget:      tx.Group.Budget,
				CreatedAt:   tx.Group.CreatedAt,
				UpdatedAt:   tx.Group.UpdatedAt,
			}
		}
		return &Transaction{
			ID:        tx.ID,
			Amount:    tx.Amount,
			Target:    tx.Target,
			Request:   tx.Request,
			Tags:      tags,
			Group:     group,
			CreatedAt: tx.CreatedAt,
			UpdatedAt: tx.UpdatedAt,
		}
	})

	return c.JSON(http.StatusOK, res)
}

func (h Handlers) PostTransaction(c echo.Context) error {
	var tx *TransactionOverview
	if err := c.Bind(&tx); err != nil {
		h.Logger.Info("could not get transaction overview from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	transactions := []*Transaction{}
	ctx := c.Request().Context()
	for _, target := range tx.Targets {
		if target == nil {
			h.Logger.Info("target is nil")
			return echo.NewHTTPError(http.StatusBadRequest, "target is nil")
		}
		created, err := h.Repository.CreateTransaction(ctx, tx.Amount, *target, tx.Tags, tx.Group, tx.Request)
		if err != nil {
			h.Logger.Error("failed to create transaction in repository", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		tags := lo.Map(created.Tags, func(tag *model.Tag, index int) *TagOverview {
			return &TagOverview{
				ID:        tag.ID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt,
				UpdatedAt: tag.UpdatedAt,
			}
		})

		var group *GroupOverview
		if created.Group != nil {
			group = &GroupOverview{
				ID:          created.Group.ID,
				Name:        created.Group.Name,
				Description: created.Group.Description,
				Budget:      created.Group.Budget,
				CreatedAt:   created.Group.CreatedAt,
				UpdatedAt:   created.Group.UpdatedAt,
			}
		}
		res := Transaction{
			ID:        created.ID,
			Amount:    created.Amount,
			Target:    created.Target,
			Request:   created.Request,
			Tags:      tags,
			Group:     group,
			CreatedAt: created.CreatedAt,
			UpdatedAt: created.UpdatedAt,
		}
		transactions = append(transactions, &res)
	}

	return c.JSON(http.StatusOK, transactions)
}

func (h Handlers) GetTransaction(c echo.Context) error {
	txID, err := uuid.Parse(c.Param("transactionID"))
	if err != nil {
		h.Logger.Info("could not parse query parameter `transactionID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	tx, err := h.Repository.GetTransaction(ctx, txID)
	if err != nil {
		h.Logger.Error("failed to get transaction from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tags := lo.Map(tx.Tags, func(tag *model.Tag, index int) *TagOverview {
		return &TagOverview{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}
	})

	var group *GroupOverview
	if tx.Group != nil {
		group = &GroupOverview{
			ID:          tx.Group.ID,
			Name:        tx.Group.Name,
			Description: tx.Group.Description,
			Budget:      tx.Group.Budget,
			CreatedAt:   tx.Group.CreatedAt,
			UpdatedAt:   tx.Group.UpdatedAt,
		}
	}
	res := Transaction{
		ID:        tx.ID,
		Amount:    tx.Amount,
		Target:    tx.Target,
		Request:   tx.Request,
		Tags:      tags,
		Group:     group,
		CreatedAt: tx.CreatedAt,
		UpdatedAt: tx.UpdatedAt,
	}

	return c.JSON(http.StatusOK, &res)
}

func (h Handlers) PutTransaction(c echo.Context) error {
	txID, err := uuid.Parse(c.Param("transactionID"))
	if err != nil {
		h.Logger.Info("could not parse query parameter `transactionID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var tx *TransactionOverviewWithOneTarget
	if err := c.Bind(&tx); err != nil {
		h.Logger.Info("could not get transaction overview with one target from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	updated, err := h.Repository.UpdateTransaction(ctx, txID, tx.Amount, tx.Target, tx.Tags, tx.Group, tx.Request)
	if err != nil {
		h.Logger.Error("failed to update transaction in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tags := lo.Map(updated.Tags, func(tag *model.Tag, index int) *TagOverview {
		return &TagOverview{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}
	})

	var group *GroupOverview
	if updated.Group != nil {
		group = &GroupOverview{
			ID:          updated.Group.ID,
			Name:        updated.Group.Name,
			Description: updated.Group.Description,
			Budget:      updated.Group.Budget,
			CreatedAt:   updated.Group.CreatedAt,
			UpdatedAt:   updated.Group.UpdatedAt,
		}
	}
	res := Transaction{
		ID:        updated.ID,
		Amount:    updated.Amount,
		Target:    updated.Target,
		Request:   updated.Request,
		Tags:      tags,
		Group:     group,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
	}

	return c.JSON(http.StatusOK, &res)
}
