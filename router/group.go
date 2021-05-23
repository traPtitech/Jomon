package router

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Group struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Budget      *int         `json:"budget"`
	Owners      []*uuid.UUID `json:"owners"`
}

type GroupOverview struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Budget      *int      `json:"budget"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GroupDetail struct {
	ID           uuid.UUID    `json:"id"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Budget       *int         `json:"budget"`
	Owners       []*uuid.UUID `json:"owners"`
	Transactions []*uuid.UUID `json:"transactions"`
	Requests     []*uuid.UUID `json:"requests"`
	Users        []*uuid.UUID `json:"users"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func (h *Handlers) GetGroups(c echo.Context) error {
	ctx := context.Background()
	groups, err := h.Repository.GetGroups(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	res := []*GroupOverview{}
	for _, group := range groups {
		res = append(res, &GroupOverview{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			Budget:      group.Budget,
			CreatedAt:   group.CreatedAt,
			UpdatedAt:   group.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) PostGroup(c echo.Context) error {
	var group Group
	if err := c.Bind(&group); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	ctx := context.Background()
	created, err := h.Repository.CreateGroup(ctx, group.Name, group.Description, group.Budget)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	res := GroupDetail{
		ID:           created.ID,
		Name:         created.Name,
		Description:  created.Description,
		Budget:       created.Budget,
		Owners:       created.Owners,
		Transactions: created.Transactions,
		Requests:     created.Requests,
		Users:        created.Users,
		CreatedAt:    created.CreatedAt,
		UpdatedAt:    created.UpdatedAt,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handlers) GetGroup(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) PostGroupUser(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) PutGroup(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) DeleteGroup(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}
