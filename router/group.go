package router

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/Jomon/ent"
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
func (h *Handlers) GetGroups(c echo.Context) error {
	ctx := c.Request().Context()
	groups, err := h.Repository.GetGroups(ctx)
	if err != nil {
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

// PostGroup POST /groups
func (h *Handlers) PostGroup(c echo.Context) error {
	var group Group
	if err := c.Bind(&group); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	created, err := h.Repository.CreateGroup(ctx, group.Name, group.Description, group.Budget)
	if err != nil {
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
func (h *Handlers) GetGroupDetail(c echo.Context) error {
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}

	ctx := c.Request().Context()
	group, err := h.Repository.GetGroup(ctx, groupID)
	if err != nil {
		if ent.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
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
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	for _, owner := range owners {
		res.Owners = append(res.Owners, &owner.ID)
	}
	members, err := h.Repository.GetMembers(ctx, groupID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	for _, member := range members {
		res.Members = append(res.Members, &member.ID)
	}

	return c.JSON(http.StatusOK, res)
}

// PutGroup PUT /groups/:groupID
func (h *Handlers) PutGroup(c echo.Context) error {
	var group Group
	if err := c.Bind(&group); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}

	ctx := c.Request().Context()
	updated, err := h.Repository.UpdateGroup(ctx, groupID, group.Name, group.Description, group.Budget)
	if err != nil {
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
func (h *Handlers) DeleteGroup(c echo.Context) error {
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}

	ctx := c.Request().Context()
	err = h.Repository.DeleteGroup(ctx, groupID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

// PostMember POST /groups/:groupID/members
func (h *Handlers) PostMember(c echo.Context) error {
	ctx := c.Request().Context()
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	var member []uuid.UUID
	if err := c.Bind(&member); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	added, err := h.Repository.AddMembers(ctx, groupID, member)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	res := []*uuid.UUID{}
	for _, m := range added {
		res = append(res, &m.ID)
	}
	return c.JSON(http.StatusOK, res)
}

// DeleteMember DELETE /groups/:groupID/members
func (h *Handlers) DeleteMember(c echo.Context) error {
	ctx := c.Request().Context()
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	var member []uuid.UUID
	if err := c.Bind(&member); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err = h.Repository.DeleteMembers(ctx, groupID, member)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

// PostOwner POST /groups/:groupID/owners
func (h *Handlers) PostOwner(c echo.Context) error {
	ctx := c.Request().Context()
	var owners []uuid.UUID
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	if err := c.Bind(&owners); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	added, err := h.Repository.AddOwners(ctx, groupID, owners)
	if err != nil {
		if ent.IsConstraintError(err) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := make([]uuid.UUID, len(added))
	for i, owner := range added {
		res[i] = owner.ID
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteOwner DELETE /groups/:groupID/owners
func (h *Handlers) DeleteOwner(c echo.Context) error {
	ctx := c.Request().Context()
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	var ownerIDs []uuid.UUID
	if err := c.Bind(&ownerIDs); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = h.Repository.DeleteOwners(ctx, groupID, ownerIDs)
	if err != nil {
		if ent.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
