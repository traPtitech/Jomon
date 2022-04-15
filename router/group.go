package router

import (
	"context"
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

// GetGroups GET /groups
func (h *Handlers) GetGroups(c echo.Context) error {
	ctx := context.Background()
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

	ctx := context.Background()
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

	ctx := context.Background()
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

	ctx := context.Background()
	err = h.Repository.DeleteGroup(ctx, groupID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

// GetMembers GET /groups/:groupID/members
func (h *Handlers) GetMembers(c echo.Context) error {
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}

	ctx := context.Background()
	members, err := h.Repository.GetMembers(ctx, groupID)
	if err != nil {
		if ent.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var res []uuid.UUID
	for _, member := range members {
		res = append(res, member.ID)
	}

	return c.JSON(http.StatusOK, &MemberResponse{res})
}

// PostMember POST /groups/:groupID/members
func (h *Handlers) PostMember(c echo.Context) error {
	var member Member
	if err := c.Bind(&member); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}

	ctx := context.Background()
	added, err := h.Repository.AddMember(ctx, groupID, member.ID)
	if err != nil {
		if ent.IsConstraintError(err) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := added.ID

	return c.JSON(http.StatusOK, &Member{res})
}

// DeleteMember DELETE /groups/:groupID/members
func (h *Handlers) DeleteMember(c echo.Context) error {
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}

	memberID, err := uuid.Parse(c.Param("memberID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if memberID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}

	ctx := context.Background()
	err = h.Repository.DeleteMember(ctx, groupID, memberID)
	if err != nil {
		if ent.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

// GetOwners GET /groups/:groupID/owners
func (h *Handlers) GetOwners(c echo.Context) error {
	ctx := context.Background()
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	owners, err := h.Repository.GetOwners(ctx, groupID)
	if err != nil {
		if ent.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var res []uuid.UUID
	for _, owner := range owners {
		res = append(res, owner.ID)
	}

	return c.JSON(http.StatusOK, &OwnerResponse{res})
}

// PostOwner POST /groups/:groupID/owners
func (h *Handlers) PostOwner(c echo.Context) error {
	ctx := context.Background()
	var owner Owner
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	if err := c.Bind(&owner); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	added, err := h.Repository.AddOwner(ctx, groupID, owner.ID)
	if err != nil {
		if ent.IsConstraintError(err) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := &Owner{
		ID: added.ID,
	}

	return c.JSON(http.StatusOK, res)
}

// DeleteOwner DELETE /groups/:groupID/owners
func (h *Handlers) DeleteOwner(c echo.Context) error {
	ctx := context.Background()
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	ownerID, err := uuid.Parse(c.Param("ownerID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if ownerID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}

	err = h.Repository.DeleteOwner(ctx, groupID, ownerID)
	if err != nil {
		if ent.IsNotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}
