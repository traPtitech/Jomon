package router

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/internal/logging"
	"github.com/traPtitech/Jomon/internal/service"
	"go.uber.org/zap"
)

func (h Handlers) GetAccountManagers(c *echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)
	loginUser, _ := c.Get(loginUserKey).(*service.User)
	accountManagers, err := h.Service.GetAccountManagers(ctx, loginUser)
	if err != nil {
		logger.Error("failed to get accountManagers from service", zap.Error(err))
		return err
	}

	res := lo.Map(accountManagers, func(accountManager *service.AccountManager, _ int) *uuid.UUID {
		return &accountManager.ID
	})

	return c.JSON(http.StatusOK, res)
}

func (h Handlers) PostAccountManagers(c *echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)
	loginUser, _ := c.Get(loginUserKey).(*service.User)

	var accountManager []uuid.UUID
	if err := c.Bind(&accountManager); err != nil {
		logger.Info("failed to get accountManager IDs from request", zap.Error(err))
		return service.NewBadInputError("failed to get accountManager IDs from request").
			WithInternal(err)
	}

	err := h.Service.AddAccountManagers(ctx, loginUser, accountManager)
	if err != nil {
		logger.Error("failed to add accountManager in service", zap.Error(err))
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h Handlers) DeleteAccountManagers(c *echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)
	loginUser, _ := c.Get(loginUserKey).(*service.User)

	var accountManager []uuid.UUID
	if err := c.Bind(&accountManager); err != nil {
		logger.Info("failed to get accountManager IDs from request", zap.Error(err))
		return service.NewBadInputError("failed to get accountManager IDs from request").
			WithInternal(err)
	}

	err := h.Service.DeleteAccountManagers(ctx, loginUser, accountManager)
	if err != nil {
		logger.Error("failed to delete accountManager from service", zap.Error(err))
		return err
	}

	return c.NoContent(http.StatusOK)
}
