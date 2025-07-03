package router

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/logging"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
	"go.uber.org/zap"
)

type TransactionResponse struct {
	ID        uuid.UUID      `json:"id"`
	Title     string         `json:"title"`
	Amount    int            `json:"amount"`
	Target    string         `json:"target"`
	Request   uuid.NullUUID      `json:"request"`
	Tags      []*TagOverview `json:"tags"`
	Group     *GroupResponse `json:"group"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type TransactionOverview struct {
	Title   string      `json:"title"`
	Amount  int         `json:"amount"`
	Targets []*string   `json:"targets"`
	Tags    []uuid.UUID `json:"tags"`
	Group   uuid.NullUUID   `json:"group"`
	Request uuid.NullUUID   `json:"request"`
}

type TransactionOverviewWithOneTarget struct {
	Title   string      `json:"title"`
	Amount  int         `json:"amount"`
	Target  string      `json:"target"`
	Tags    []uuid.UUID `json:"tags"`
	Group   uuid.NullUUID   `json:"group"`
	Request uuid.NullUUID   `json:"request"`
}

func (h Handlers) GetTransactions(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

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
	var since time.Time
	if c.QueryParam("since") != "" {
		var err error
		s, err := service.StrToDate(c.QueryParam("since"))
		if err != nil {
			logger.Info("could not parse since as time.Time", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		since = s
	}
	var until time.Time
	if c.QueryParam("until") != "" {
		var err error
		u, err := service.StrToDate(c.QueryParam("until"))
		if err != nil {
			logger.Info("could not parse until as time.Time", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		until = u
	}
	limit := 100
	if limitQuery := c.QueryParam("limit"); limitQuery != "" {
		limitI, err := strconv.Atoi(limitQuery)
		if err != nil {
			logger.Info("could not parse limit as integer", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		if limitI < 0 {
			logger.Info("received negative limit", zap.Int("limit", limitI))
			err := fmt.Errorf("negative limit(=%d) is invalid", limitI)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		limit = limitI
	}
	offset := 0
	if offsetQuery := c.QueryParam("offset"); offsetQuery != "" {
		offsetI, err := strconv.Atoi(offsetQuery)
		if err != nil {
			logger.Info("could not parse offset as integer", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		if offsetI < 0 {
			logger.Info("received negative offset", zap.Int("offset", offsetI))
			err := fmt.Errorf("negative offset(=%d) is invalid", offsetI)
			return echo.NewHTTPError(http.StatusBadRequest, err)
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
	var request uuid.UUID
	if c.QueryParam("request") != "" {
		var r uuid.UUID
		r, err := uuid.Parse(c.QueryParam("request"))
		if err != nil {
			logger.Info("could not parse request as uuid.UUID", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		request = r
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

	res := lo.Map(txs, func(tx *model.TransactionResponse, _ int) *TransactionResponse {
		tags := lo.Map(tx.Tags, func(tag *model.Tag, _ int) *TagOverview {
			return &TagOverview{
				ID:        tag.ID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt,
				UpdatedAt: tag.UpdatedAt,
			}
		})

		var group *GroupResponse
		if tx.Group != nil {
			group = &GroupResponse{
				ID:          tx.Group.ID,
				Name:        tx.Group.Name,
				Description: tx.Group.Description,
				Budget:      tx.Group.Budget,
				CreatedAt:   tx.Group.CreatedAt,
				UpdatedAt:   tx.Group.UpdatedAt,
			}
		}
		return &TransactionResponse{
			ID:        tx.ID,
			Title:     tx.Title,
			Amount:    tx.Amount,
			Target:    tx.Target,
			Request:   uuid.NullUUID{UUID:tx.Request,Valid: true},
			Tags:      tags,
			Group:     group,
			CreatedAt: tx.CreatedAt,
			UpdatedAt: tx.UpdatedAt,
		}
	})

	return c.JSON(http.StatusOK, res)
}

func (h Handlers) PostTransaction(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	var tx *PostTransactionsRequest
	// TODO: validate
	if err := c.Bind(&tx); err != nil {
		logger.Info("could not get transaction overview from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	transactions := []*TransactionResponse{}
	for _, target := range tx.Targets {
		if target == nil {
			logger.Info("target is nil")
			return echo.NewHTTPError(http.StatusBadRequest, "target is nil")
		}
		created, err := h.Repository.CreateTransaction(
			ctx,
			tx.Title, tx.Amount, *target, tx.Tags, tx.Group.UUID, tx.Request.UUID)
		if err != nil {
			logger.Error("failed to create transaction in repository", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		tags := lo.Map(created.Tags, func(tag *model.Tag, _ int) *TagOverview {
			return &TagOverview{
				ID:        tag.ID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt,
				UpdatedAt: tag.UpdatedAt,
			}
		})

		var group *GroupResponse
		if created.Group != nil {
			group = &GroupResponse{
				ID:          created.Group.ID,
				Name:        created.Group.Name,
				Description: created.Group.Description,
				Budget:      created.Group.Budget,
				CreatedAt:   created.Group.CreatedAt,
				UpdatedAt:   created.Group.UpdatedAt,
			}
		}
		res := TransactionResponse{
			ID:        created.ID,
			Title:     created.Title,
			Amount:    created.Amount,
			Target:    created.Target,
			Request:   uuid.NullUUID{UUID:created.Request,Valid: true},
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
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	txID, err := uuid.Parse(c.Param("transactionID"))
	if err != nil {
		logger.Info("could not parse query parameter `transactionID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	tx, err := h.Repository.GetTransaction(ctx, txID)
	if err != nil {
		logger.Error("failed to get transaction from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tags := lo.Map(tx.Tags, func(tag *model.Tag, _ int) *TagOverview {
		return &TagOverview{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}
	})

	var group *GroupResponse
	if tx.Group != nil {
		group = &GroupResponse{
			ID:          tx.Group.ID,
			Name:        tx.Group.Name,
			Description: tx.Group.Description,
			Budget:      tx.Group.Budget,
			CreatedAt:   tx.Group.CreatedAt,
			UpdatedAt:   tx.Group.UpdatedAt,
		}
	}
	res := TransactionResponse{
		ID:        tx.ID,
		Title:     tx.Title,
		Amount:    tx.Amount,
		Target:    tx.Target,
		Request:   uuid.NullUUID{UUID:tx.Request,Valid: true},
		Tags:      tags,
		Group:     group,
		CreatedAt: tx.CreatedAt,
		UpdatedAt: tx.UpdatedAt,
	}

	return c.JSON(http.StatusOK, &res)
}

func (h Handlers) PutTransaction(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	txID, err := uuid.Parse(c.Param("transactionID"))
	if err != nil {
		logger.Info("could not parse query parameter `transactionID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var tx *PutTransactionRequest
	// TODO: validate
	if err := c.Bind(&tx); err != nil {
		logger.Info(
			"could not get transaction overview with one target from request",
			zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	updated, err := h.Repository.UpdateTransaction(
		ctx,
		txID, tx.Title, tx.Amount, tx.Target, tx.Tags, tx.Group.UUID, tx.Request.UUID)
	if err != nil {
		logger.Error("failed to update transaction in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tags := lo.Map(updated.Tags, func(tag *model.Tag, _ int) *TagOverview {
		return &TagOverview{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}
	})

	var group *GroupResponse
	if updated.Group != nil {
		group = &GroupResponse{
			ID:          updated.Group.ID,
			Name:        updated.Group.Name,
			Description: updated.Group.Description,
			Budget:      updated.Group.Budget,
			CreatedAt:   updated.Group.CreatedAt,
			UpdatedAt:   updated.Group.UpdatedAt,
		}
	}
	res := TransactionResponse{
		ID:        updated.ID,
		Title:     updated.Title,
		Amount:    updated.Amount,
		Target:    updated.Target,
		Request:   uuid.NullUUID{UUID:updated.Request,Valid: true},
		Tags:      tags,
		Group:     group,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
	}

	return c.JSON(http.StatusOK, &res)
}
