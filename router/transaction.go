package router

import (
	"context"
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
	Tags      []*TagOverview `json:"tags"`
	Group     *GroupOverview `json:"group"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type TransactionOverview struct {
	Amount  int          `json:"amount"`
	Target  string       `json:"target"`
	Tags    []*uuid.UUID `json:"tags"`
	Group   *uuid.UUID   `json:"group"`
	Request *uuid.UUID   `json:"request"`
}

func (h *Handlers) GetTransactions(c echo.Context) error {
	ctx := context.Background()
	var sort *string = nil
	if c.Param("sort") != "" {
		s := c.Param("sort")
		sort = &s
	}
	var target *string = nil
	if c.Param("target") != "" {
		t := c.Param("target")
		target = &t
	}
	var since *time.Time = nil
	if c.Param("since") != "" {
		var err error
		s, err := service.StrToDate(c.Param("since"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		since = &s
	}
	var until *time.Time = nil
	if c.Param("until") != "" {
		var err error
		u, err := service.StrToDate(c.Param("until"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		until = &u
	}
	var tag *string = nil
	if c.Param("tag") != "" {
		t := c.Param("tag")
		tag = &t
	}
	var group *string = nil
	if c.Param("group") != "" {
		g := c.Param("group")
		group = &g
	}
	query := model.TransactionQuery{
		Sort:   sort,
		Target: target,
		Since:  since,
		Until:  until,
		Tag:    tag,
		Group:  group,
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
				ID:          tag.ID,
				Name:        tag.Name,
				Description: tag.Description,
				CreatedAt:   tag.CreatedAt,
				UpdatedAt:   tag.UpdatedAt,
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
		res = append(res, &Transaction{
			ID:        tx.ID,
			Amount:    tx.Amount,
			Target:    tx.Target,
			Tags:      tags,
			Group:     group,
			CreatedAt: tx.CreatedAt,
			UpdatedAt: tx.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) PostTransaction(c echo.Context) error {
	var tx *TransactionOverview
	if err := c.Bind(&tx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := context.Background()
	created, err := h.Repository.CreateTransaction(ctx, tx.Amount, tx.Target, tx.Tags, tx.Group, tx.Request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tags := []*TagOverview{}
	for _, tag := range created.Tags {
		tags = append(tags, &TagOverview{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
			CreatedAt:   tag.CreatedAt,
			UpdatedAt:   tag.UpdatedAt,
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
		Tags:      tags,
		Group:     group,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}

	return c.JSON(http.StatusOK, &res)
}

func (h *Handlers) GetTransaction(c echo.Context) error {
	txID, err := uuid.Parse(c.Param("transactionID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := context.Background()
	tx, err := h.Repository.GetTransaction(ctx, txID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tags := []*TagOverview{}
	for _, tag := range tx.Tags {
		tags = append(tags, &TagOverview{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
			CreatedAt:   tag.CreatedAt,
			UpdatedAt:   tag.UpdatedAt,
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

	var tx *TransactionOverview
	if err := c.Bind(&tx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := context.Background()
	updated, err := h.Repository.UpdateTransaction(ctx, txID, tx.Amount, tx.Target, tx.Tags, tx.Group, tx.Request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tags := []*TagOverview{}
	for _, tag := range updated.Tags {
		tags = append(tags, &TagOverview{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
			CreatedAt:   tag.CreatedAt,
			UpdatedAt:   tag.UpdatedAt,
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
		Tags:      tags,
		Group:     group,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
	}

	return c.JSON(http.StatusOK, &res)
}
