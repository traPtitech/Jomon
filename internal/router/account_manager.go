package router

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/internal/logging"
	"github.com/traPtitech/Jomon/internal/model"
	"github.com/traPtitech/Jomon/internal/service"
	"go.uber.org/zap"
)

func (h Handlers) GetAccountManagers(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)
	accountManagers, err := h.Repository.GetAccountManagers(ctx)
	if err != nil {
		logger.Error("failed to get accountManagers from repository", zap.Error(err))
		return err
	}

	res := lo.Map(accountManagers, func(accountManager *model.AccountManager, _ int) *uuid.UUID {
		return &accountManager.ID
	})

	return c.JSON(http.StatusOK, res)
}

func (h Handlers) PostAccountManagers(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	var accountManager []uuid.UUID
	if err := c.Bind(&accountManager); err != nil {
		logger.Info("failed to get accountManager IDs from request", zap.Error(err))
		return service.NewBadInputError("failed to get accountManager IDs from request").
			WithInternal(err)
	}

	err := h.Repository.AddAccountManagers(ctx, accountManager)
	if err != nil {
		logger.Error("failed to add accountManager in repository", zap.Error(err))
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h Handlers) DeleteAccountManagers(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	var accountManager []uuid.UUID
	if err := c.Bind(&accountManager); err != nil {
		logger.Info("failed to get accountManager IDs from request", zap.Error(err))
		return service.NewBadInputError("failed to get accountManager IDs from request").
			WithInternal(err)
	}

	err := h.Repository.DeleteAccountManagers(ctx, accountManager)
	if err != nil {
		logger.Error("failed to delete accountManager from repository", zap.Error(err))
		return err
	}

	return c.NoContent(http.StatusOK)
}
