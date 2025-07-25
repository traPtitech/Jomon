package router

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/logging"
	"github.com/traPtitech/Jomon/model"
	"github.com/traPtitech/Jomon/service"
	"go.uber.org/zap"
)

type User struct {
	ID             uuid.UUID        `json:"id"`
	Name           string           `json:"name"`
	DisplayName    string           `json:"display_name"`
	AccountManager bool             `json:"account_manager"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      service.NullTime `json:"deleted_at"`
}

func (h Handlers) GetUsers(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	users, err := h.Repository.GetUsers(ctx)
	if err != nil {
		logger.Error("failed to get users from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := lo.Map(users, func(user *model.User, _ int) User {
		return User{
			ID:             user.ID,
			Name:           user.Name,
			DisplayName:    user.DisplayName,
			AccountManager: user.AccountManager,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			DeletedAt:      service.TimeToNullTime(&user.DeletedAt),
		}
	})

	return c.JSON(http.StatusOK, res)
}

type PutUserRequest struct {
	Name           string `json:"name"`
	DisplayName    string `json:"display_name"`
	AccountManager bool   `json:"account_manager"`
}

func (h Handlers) UpdateUserInfo(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	var newUser PutUserRequest
	if err := c.Bind(&newUser); err != nil {
		logger.Info("could not get user info from request", zap.Error(err))
		return c.NoContent(http.StatusBadRequest)
	}

	user, err := h.Repository.GetUserByName(ctx, newUser.Name)
	if err != nil {
		logger.Error("failed to get user from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	updated, err := h.Repository.UpdateUser(
		ctx, user.ID, newUser.Name, newUser.DisplayName, newUser.AccountManager)
	if err != nil {
		logger.Error("failed to update user in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, User{
		ID:             user.ID,
		Name:           updated.Name,
		DisplayName:    updated.DisplayName,
		AccountManager: updated.AccountManager,
	})
}

func userFromModelUser(u model.User) User {
	return User{
		ID:             u.ID,
		Name:           u.Name,
		DisplayName:    u.DisplayName,
		AccountManager: u.AccountManager,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		DeletedAt:      service.TimeToNullTime(&u.DeletedAt),
	}
}

func (h Handlers) GetMe(c echo.Context) error {
	loginUser, _ := c.Get(loginUserKey).(User)
	return c.JSON(http.StatusOK, loginUser)
}
