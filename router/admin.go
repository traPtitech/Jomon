package router

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/logging"
	"github.com/traPtitech/Jomon/model"
	"go.uber.org/zap"
)

func (h Handlers) GetAdmins(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)
	admins, err := h.Repository.GetAdmins(ctx)
	if err != nil {
		logger.Error("failed to get admins from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	res := lo.Map(admins, func(admin *model.Admin, _ int) *uuid.UUID {
		return &admin.ID
	})

	return c.JSON(http.StatusOK, res)
}

func (h Handlers) PostAdmins(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	var admin []uuid.UUID
	if err := c.Bind(&admin); err != nil {
		logger.Info("failed to get admin id from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err := h.Repository.AddAdmins(ctx, admin)
	if err != nil {
		if ent.IsConstraintError(err) {
			logger.Info("constraint error while adding admin in repository", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		logger.Error("failed to add admin in repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

func (h Handlers) DeleteAdmins(c echo.Context) error {
	ctx := c.Request().Context()
	logger := logging.GetLogger(ctx)

	var admin []uuid.UUID
	if err := c.Bind(&admin); err != nil {
		logger.Info("failed to get admin id from request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err := h.Repository.DeleteAdmins(ctx, admin)
	if err != nil {
		logger.Error("failed to delete admin from repository", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

func (h Handlers) isAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	logger := logging.GetLogger(ctx)
	user, err := h.Repository.GetUserByID(ctx, userID)
	if err != nil {
		if ent.IsNotFound(err) {
			return false, nil
		}
		logger.Error("failed to get user by id", zap.Error(err))
		return false, err
	}
	return user.Admin, nil
}
