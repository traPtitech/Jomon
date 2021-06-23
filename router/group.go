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
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Budget      *int         `json:"budget"`
	Owners      []*uuid.UUID `json:"owners"`
	Users       []*uuid.UUID `json:"users"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (h *Handlers) GetGroups(c echo.Context) error {
	ctx := context.Background()
	groups, err := h.Repository.GetGroups(ctx)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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
	return c.NoContent(http.StatusOK)
	// TODO: Implement
	/*
		var group Group
		if err := c.Bind(&group); err != nil {
			return badRequest(err)
		}

		ctx := context.Background()
		created, err := h.Repository.CreateGroup(ctx, group.Name, group.Description, group.Budget, owners)
		if err != nil {
			return internalServerError(err)
		}

		res := GroupDetail{
			ID:          created.ID,
			Name:        created.Name,
			Description: created.Description,
			Budget:      created.Budget,
			//Owners:      created.Owners,
			//Users:       created.Users,
			CreatedAt: created.CreatedAt,
			UpdatedAt: created.UpdatedAt,
		}

		return c.JSON(http.StatusOK, res)
	*/
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

func (h *Handlers) GetMembers(c echo.Context) error {
	ctx := c.Request().Context()
	members, err := h.Repository.GetMembers(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK, members)
}

func (h *Handlers) PostMember(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) DeleteMember(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) GetOwners(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}

func (h *Handlers) PostOwner(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}
func (h *Handlers) DeleteOwner(c echo.Context) error {
	return c.NoContent(http.StatusOK)
	// TODO: Implement
}
