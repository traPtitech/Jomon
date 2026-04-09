package router

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/internal/logging"
	"github.com/traPtitech/Jomon/internal/nulltime"
	"github.com/traPtitech/Jomon/internal/service"
	"go.uber.org/zap"
)

type User struct {
	ID             uuid.UUID         `json:"id"`
	Name           string            `json:"name"`
	DisplayName    string            `json:"display_name"`
	AccountManager bool              `json:"account_manager"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	DeletedAt      nulltime.NullTime `json:"deleted_at"`
}

func (h Handlers) GetUsers(c *echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	users, err := h.Service.GetUsers(ctx)
	if err != nil {
		logger.Error("failed to get users from service", zap.Error(err))
		return err
	}

	res := lo.Map(users, func(user *service.User, _ int) User {
		return User{
			ID:             user.ID,
			Name:           user.Name,
			DisplayName:    user.DisplayName,
			AccountManager: user.AccountManager,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			DeletedAt:      user.DeletedAt,
		}
	})

	return c.JSON(http.StatusOK, res)
}

type PutUserRequest struct {
	Name           string `json:"name"`
	DisplayName    string `json:"display_name"`
	AccountManager bool   `json:"account_manager"`
}

func (h Handlers) UpdateUserInfo(c *echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)
	loginUser, _ := c.Get(loginUserKey).(*service.User)

	var newUser PutUserRequest
	if err := c.Bind(&newUser); err != nil {
		logger.Info("could not get user info from request", zap.Error(err))
		return service.NewBadInputError("could not get user info from request").
			WithInternal(err)
	}

	user, err := h.Service.GetUserByName(ctx, newUser.Name)
	if err != nil {
		logger.Error("failed to get user from service", zap.Error(err))
		return err
	}

	inputs := service.UpdateUserInputs{
		Name:           newUser.Name,
		DisplayName:    newUser.DisplayName,
		AccountManager: newUser.AccountManager,
	}
	updated, err := h.Service.UpdateUser(ctx, loginUser, user.ID, inputs)
	if err != nil {
		logger.Error("failed to update user in service", zap.Error(err))
		return err
	}

	return c.JSON(http.StatusOK, User{
		ID:             user.ID,
		Name:           updated.Name,
		DisplayName:    updated.DisplayName,
		AccountManager: updated.AccountManager,
	})
}

func userFromModelUser(u service.User) User {
	return User{
		ID:             u.ID,
		Name:           u.Name,
		DisplayName:    u.DisplayName,
		AccountManager: u.AccountManager,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		DeletedAt:      u.DeletedAt,
	}
}

func (h Handlers) GetMe(c *echo.Context) error {
	loginUser, _ := c.Get(loginUserKey).(*service.User)
	user := userFromModelUser(*loginUser)
	return c.JSON(http.StatusOK, user)
}
