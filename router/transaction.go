package router

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
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

func (h *Handlers) GetTransactions(c echo.Context) error {
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
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		since = &s
	}
	var until *time.Time
	if c.QueryParam("until") != "" {
		var err error
		u, err := service.StrToDate(c.QueryParam("until"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		until = &u
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
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		request = &r
	}
	query := model.TransactionQuery{
		Sort:    sort,
		Target:  target,
		Since:   since,
		Until:   until,
		Tag:     tag,
		Group:   group,
		Request: request,
	}
	txs, err := h.Repository.GetTransactions(ctx, query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := []*Transaction{}
	for _, tx := range txs {
		tags := []*TagOverview{}
		for _, tag := range tx.Tags {
			tags = append(tags, &TagOverview{
				ID:        tag.ID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt,
				UpdatedAt: tag.UpdatedAt,
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
		tx := &Transaction{
			ID:        tx.ID,
			Amount:    tx.Amount,
			Target:    tx.Target,
			Request:   tx.Request,
			Tags:      tags,
			Group:     group,
			CreatedAt: tx.CreatedAt,
			UpdatedAt: tx.UpdatedAt,
		}
		res = append(res, tx)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) PostTransaction(c echo.Context) error {
	var tx *TransactionOverview
	if err := c.Bind(&tx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	transactions := []*Transaction{}
	ctx := c.Request().Context()
	for _, target := range tx.Targets {
		if target == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "target is nil")
		}
		created, err := h.Repository.CreateTransaction(ctx, tx.Amount, *target, tx.Tags, tx.Group, tx.Request)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		tags := []*TagOverview{}
		for _, tag := range created.Tags {
			tags = append(tags, &TagOverview{
				ID:        tag.ID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt,
				UpdatedAt: tag.UpdatedAt,
			})
		}
		group := &GroupOverview{
			ID:          created.Group.ID,
			Name:        created.Group.Name,
			Description: created.Group.Description,
			Budget:      created.Group.Budget,
			CreatedAt:   created.Group.CreatedAt,
			UpdatedAt:   created.Group.UpdatedAt,
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

func (h *Handlers) GetTransaction(c echo.Context) error {
	txID, err := uuid.Parse(c.Param("transactionID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	tx, err := h.Repository.GetTransaction(ctx, txID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tags := []*TagOverview{}
	for _, tag := range tx.Tags {
		tags = append(tags, &TagOverview{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
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

func (h *Handlers) PutTransaction(c echo.Context) error {
	txID, err := uuid.Parse(c.Param("transactionID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var tx *TransactionOverviewWithOneTarget
	if err := c.Bind(&tx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	updated, err := h.Repository.UpdateTransaction(ctx, txID, tx.Amount, tx.Target, tx.Tags, tx.Group, tx.Request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tags := []*TagOverview{}
	for _, tag := range updated.Tags {
		tags = append(tags, &TagOverview{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}
	group := &GroupOverview{
		ID:          updated.Group.ID,
		Name:        updated.Group.Name,
		Description: updated.Group.Description,
		Budget:      updated.Group.Budget,
		CreatedAt:   updated.Group.CreatedAt,
		UpdatedAt:   updated.Group.UpdatedAt,
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
