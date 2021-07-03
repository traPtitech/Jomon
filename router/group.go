package router

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/ent"
)

type Group struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Budget      *int         `json:"budget"`
	Owners      []*uuid.UUID `json:"owners"`
}

type GroupResponse struct {
	Members []*Group `json:"member"`
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

type MemberResponse struct {
	ID []uuid.UUID `json:"members"`
}

type Member struct {
	ID uuid.UUID `json:"id"`
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
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := context.Background()
	members, err := h.Repository.GetMembers(ctx, groupID)
	if err != nil {
		if ent.IsNotFound(err) {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var res []uuid.UUID
	for _, member := range members {
		res = append(res, member.ID)
	}

	return c.JSON(http.StatusOK, &MemberResponse{res})
}

func (h *Handlers) PostMember(c echo.Context) error {
	var member Member
	if err := c.Bind(&member); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := context.Background()
	created, err := h.Repository.CreateMember(ctx, groupID, member.ID)
	if err != nil {
		if ent.IsConstraintError(err) {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := created.ID

	return c.JSON(http.StatusOK, &Member{res})
}

func (h *Handlers) DeleteMember(c echo.Context) error {
	var member Member
	if err := c.Bind(&member); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := context.Background()
	err = h.Repository.DeleteMember(ctx, groupID, member.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
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
