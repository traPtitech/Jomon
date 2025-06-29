package router

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/logging"
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
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Budget      *int        `json:"budget"`
	Owners      []uuid.UUID `json:"owners"`
	Members     []uuid.UUID `json:"members"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type MemberResponse struct {
	ID []uuid.UUID `json:"members"`
}

type Member struct {
	ID uuid.UUID `json:"id"`
}

var errUserIsNotAdminOrGroupOwner error = errors.New("user is not admin or group owner")

// GetGroups GET /groups
func (h Handlers) GetGroups(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	groups, err := h.Repository.GetGroups(ctx)
	if err != nil {
		logger.Error("failed to get groups from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := lo.Map(groups, func(group *model.Group, _ int) *GroupOverview {
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
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	var group Group
	if err := c.Bind(&group); err != nil {
		logger.Info("failed to get group from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	created, err := h.Repository.CreateGroup(ctx, group.Name, group.Description, group.Budget)
	if err != nil {
		logger.Error("failed to create group in repository", zap.Error(err))
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
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		logger.Info("could not parse query parameter `groupID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		logger.Info("received invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}

	group, err := h.Repository.GetGroup(ctx, groupID)
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Info(
				"could not fin group in repository",
				zap.String("ID", groupID.String()),
				zap.Error(err))
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		logger.Error("failed to get group from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	res := GroupDetail{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		Owners:      []uuid.UUID{},
		Members:     []uuid.UUID{},
		Budget:      group.Budget,
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
	}
	owners, err := h.Repository.GetOwners(ctx, groupID)
	if err != nil {
		logger.Error("failed to get owners from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	res.Owners = lo.Map(owners, func(owner *model.Owner, _ int) uuid.UUID {
		return owner.ID
	})
	members, err := h.Repository.GetMembers(ctx, groupID)
	if err != nil {
		logger.Error("failed to get members from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	res.Members = lo.Map(members, func(member *model.Member, indec int) uuid.UUID {
		return member.ID
	})

	return c.JSON(http.StatusOK, res)
}

// PutGroup PUT /groups/:groupID
func (h Handlers) PutGroup(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	loginUser, _ := c.Get(loginUserKey).(User)
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		logger.Info("could not parse query parameter `groupID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	if err := h.filterAdminOrGroupOwner(ctx, &loginUser, groupID); err != nil {
		return err
	}

	var group Group
	if err := c.Bind(&group); err != nil {
		logger.Info("could not get group from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	updated, err := h.Repository.UpdateGroup(
		ctx,
		groupID, group.Name, group.Description, group.Budget)
	if err != nil {
		logger.Error("failed to update group in repository", zap.Error(err))
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
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	loginUser, _ := c.Get(loginUserKey).(User)
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		logger.Info("could not parse query parameter `groupID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		logger.Info("received invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	if err := h.filterAdminOrGroupOwner(ctx, &loginUser, groupID); err != nil {
		return err
	}

	err = h.Repository.DeleteGroup(ctx, groupID)
	if err != nil {
		logger.Error("failed to delete group from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

// PostMember POST /groups/:groupID/members
func (h Handlers) PostMember(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	loginUser, _ := c.Get(loginUserKey).(User)
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		logger.Info("could not parse query parameter `groupID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		logger.Info("received invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	if err := h.filterAdminOrGroupOwner(ctx, &loginUser, groupID); err != nil {
		return err
	}
	var member []uuid.UUID
	if err := c.Bind(&member); err != nil {
		logger.Info("could not get member id from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	added, err := h.Repository.AddMembers(ctx, groupID, member)
	if err != nil {
		logger.Error("failed to add member in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	res := lo.Map(added, func(m *model.Member, _ int) uuid.UUID {
		return m.ID
	})
	return c.JSON(http.StatusOK, res)
}

// DeleteMember DELETE /groups/:groupID/members
func (h Handlers) DeleteMember(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	loginUser, _ := c.Get(loginUserKey).(User)
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		logger.Info("could not parse query parameter `groupID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		logger.Info("received invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	if err := h.filterAdminOrGroupOwner(ctx, &loginUser, groupID); err != nil {
		return err
	}
	var member []uuid.UUID
	if err := c.Bind(&member); err != nil {
		logger.Info("could not get member id from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	err = h.Repository.DeleteMembers(ctx, groupID, member)
	if err != nil {
		logger.Error("failed to delete member from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

// PostOwner POST /groups/:groupID/owners
func (h Handlers) PostOwner(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	loginUser, _ := c.Get(loginUserKey).(User)
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		logger.Info("could not parse query parameter `groupID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		logger.Info("received invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	if err := h.filterAdminOrGroupOwner(ctx, &loginUser, groupID); err != nil {
		return err
	}
	var owners []uuid.UUID
	if err := c.Bind(&owners); err != nil {
		logger.Info("could not get owner id from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	added, err := h.Repository.AddOwners(ctx, groupID, owners)
	if err != nil {
		if ent.IsConstraintError(err) {
			logger.Info("constraint error while adding owner in repository", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		logger.Error("failed to add owner in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := lo.Map(added, func(owner *model.Owner, _ int) uuid.UUID {
		return owner.ID
	})

	return c.JSON(http.StatusOK, res)
}

// DeleteOwner DELETE /groups/:groupID/owners
func (h Handlers) DeleteOwner(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	loginUser, _ := c.Get(loginUserKey).(User)
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		logger.Info("could not parse query parameter `groupID` as UUID", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if groupID == uuid.Nil {
		logger.Info("received invalid UUID")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid UUID"))
	}
	if err := h.filterAdminOrGroupOwner(ctx, &loginUser, groupID); err != nil {
		return err
	}
	var ownerIDs []uuid.UUID
	if err := c.Bind(&ownerIDs); err != nil {
		logger.Info("could not get owner id from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = h.Repository.DeleteOwners(ctx, groupID, ownerIDs)
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Info("could not find owner in repository", zap.Error(err))
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		logger.Error("failed to delete owner from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

// isGroupOwner checks if the user is an owner of the group.
func (h Handlers) isGroupOwner(
	ctx context.Context, userID, groupID uuid.UUID,
) (bool, error) {
	owners, err := h.Repository.GetOwners(ctx, groupID)
	if err != nil {
		return false, err
	}
	isOwner := lo.ContainsBy(owners, func(owner *model.Owner) bool {
		return owner.ID == userID
	})

	return isOwner, nil
}

// filterAdminOrGroupOwner 与えられたIDのユーザーが管理者またはグループのオーナーであるかを確認します
func (h Handlers) filterAdminOrGroupOwner(
	ctx context.Context, user *User, groupID uuid.UUID,
) *echo.HTTPError {
	logger := logging.GetLogger(ctx)
	if user.Admin {
		return nil
	}
	isOwner, err := h.isGroupOwner(ctx, user.ID, groupID)
	if err != nil {
		// NOTE: ent.IsNotFound(err)は起こり得ないと仮定
		logger.Error("failed to check if user is group owner", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if isOwner {
		return nil
	}
	return echo.NewHTTPError(http.StatusForbidden, errUserIsNotAdminOrGroupOwner)
}
