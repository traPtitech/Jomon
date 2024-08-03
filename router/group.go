package router

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/model"
	"go.uber.org/zap"
)

type Group struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Budget      *int   `json:"budget"`
}

type Owner struct {
	ID uuid.UUID `json:"id"`
}

type OwnerResponse struct {
	Owners []uuid.UUID `json:"owners"`
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
	Members     []*uuid.UUID `json:"members"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type MemberResponse struct {
	ID []uuid.UUID `json:"members"`
}

type Member struct {
	ID uuid.UUID `json:"id"`
}

// GetGroups GET /groups
func (h Handlers) GetGroups(c echo.Context) error {
	ctx := c.Request().Context()
	groups, err := h.Repository.GetGroups(ctx)
	if err != nil {
		h.Logger.Error("failed to get groups from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := lo.Map(groups, func(group *model.Group, index int) *GroupOverview {
		return &GroupOverview{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			Budget:      group.Budget,
			CreatedAt:   group.CreatedAt,
			UpdatedAt:   group.UpdatedAt,
		}
	})

	return c.JSON(http.StatusOK, res)
}

// PostGroup POST /groups
func (h Handlers) PostGroup(c echo.Context) error {
	var group Group
	if err := c.Bind(&group); err != nil {
		h.Logger.Info("failed to get group from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	created, err := h.Repository.CreateGroup(ctx, group.Name, group.Description, group.Budget)
	if err != nil {
		h.Logger.Error("failed to create group in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := GroupOverview{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
		Budget:      created.Budget,
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}

	return c.JSON(http.StatusOK, res)
}

// GetGroupDetail GET /groups/:groupID
func (h Handlers) GetGroupDetail(c echo.Context) error {
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		h.Logger.Info("could not parse query parameter `groupID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		h.Logger.Info("received invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}

	ctx := c.Request().Context()
	group, err := h.Repository.GetGroup(ctx, groupID)
	if err != nil {
		if ent.IsNotFound(err) {
			h.Logger.Info(
				"could not fin group in repository",
				zap.String("ID", groupID.String()),
				zap.Error(err))
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		h.Logger.Error("failed to get group from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	res := GroupDetail{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		Owners:      []*uuid.UUID{},
		Members:     []*uuid.UUID{},
		Budget:      group.Budget,
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
	}
	owners, err := h.Repository.GetOwners(ctx, groupID)
	if err != nil {
		h.Logger.Error("failed to get owners from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	res.Owners = lo.Map(owners, func(owner *model.Owner, index int) *uuid.UUID {
		return &owner.ID
	})
	members, err := h.Repository.GetMembers(ctx, groupID)
	if err != nil {
		h.Logger.Error("failed to get members from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	res.Members = lo.Map(members, func(member *model.Member, indec int) *uuid.UUID {
		return &member.ID
	})

	return c.JSON(http.StatusOK, res)
}

// PutGroup PUT /groups/:groupID
func (h Handlers) PutGroup(c echo.Context) error {
	var group Group
	if err := c.Bind(&group); err != nil {
		h.Logger.Info("could not get group from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		h.Logger.Info("could not parse query parameter `groupID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}

	ctx := c.Request().Context()
	updated, err := h.Repository.UpdateGroup(
		ctx,
		groupID, group.Name, group.Description, group.Budget)
	if err != nil {
		h.Logger.Error("failed to update group in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := GroupOverview{
		ID:          updated.ID,
		Name:        updated.Name,
		Description: updated.Description,
		Budget:      updated.Budget,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteGroup DELETE /groups/:groupID
func (h Handlers) DeleteGroup(c echo.Context) error {
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		h.Logger.Info("could not parse query parameter `groupID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		h.Logger.Info("received invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}

	ctx := c.Request().Context()
	err = h.Repository.DeleteGroup(ctx, groupID)
	if err != nil {
		h.Logger.Error("failed to delete group from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

// PostMember POST /groups/:groupID/members
func (h Handlers) PostMember(c echo.Context) error {
	ctx := c.Request().Context()
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		h.Logger.Info("could not parse query parameter `groupID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		h.Logger.Info("received invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	var member []uuid.UUID
	if err := c.Bind(&member); err != nil {
		h.Logger.Info("could not get member id from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	added, err := h.Repository.AddMembers(ctx, groupID, member)
	if err != nil {
		h.Logger.Error("failed to add member in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	res := lo.Map(added, func(m *model.Member, index int) *uuid.UUID {
		return &m.ID
	})
	return c.JSON(http.StatusOK, res)
}

// DeleteMember DELETE /groups/:groupID/members
func (h Handlers) DeleteMember(c echo.Context) error {
	ctx := c.Request().Context()
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		h.Logger.Info("could not parse query parameter `groupID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		h.Logger.Info("received invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	var member []uuid.UUID
	if err := c.Bind(&member); err != nil {
		h.Logger.Info("could not get member id from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err = h.Repository.DeleteMembers(ctx, groupID, member)
	if err != nil {
		h.Logger.Error("failed to delete member from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

// PostOwner POST /groups/:groupID/owners
func (h Handlers) PostOwner(c echo.Context) error {
	ctx := c.Request().Context()
	var owners []uuid.UUID
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		h.Logger.Info("could not parse query parameter `groupID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		h.Logger.Info("received invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	if err := c.Bind(&owners); err != nil {
		h.Logger.Info("could not get owner id from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	added, err := h.Repository.AddOwners(ctx, groupID, owners)
	if err != nil {
		if ent.IsConstraintError(err) {
			h.Logger.Info("constraint error while adding owner in repository", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		h.Logger.Error("failed to add owner in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := make([]uuid.UUID, len(added))
	for i, owner := range added {
		res[i] = owner.ID
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteOwner DELETE /groups/:groupID/owners
func (h Handlers) DeleteOwner(c echo.Context) error {
	ctx := c.Request().Context()
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		h.Logger.Info("could not parse query parameter `groupID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		h.Logger.Info("received invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	var ownerIDs []uuid.UUID
	if err := c.Bind(&ownerIDs); err != nil {
		h.Logger.Info("could not get owner id from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = h.Repository.DeleteOwners(ctx, groupID, ownerIDs)
	if err != nil {
		if ent.IsNotFound(err) {
			h.Logger.Info("could not find owner in repository", zap.Error(err))
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		h.Logger.Error("failed to delete owner from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
