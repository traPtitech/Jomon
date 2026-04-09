//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package service

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/traPtitech/Jomon/internal/logging"
)

type AccountManager struct {
	ID uuid.UUID `json:"id"`
}

type AccountManagerRepository interface {
	GetAccountManagers(ctx context.Context) ([]*AccountManager, error)
	AddAccountManagers(ctx context.Context, userIDs []uuid.UUID) error
	DeleteAccountManagers(ctx context.Context, userIDs []uuid.UUID) error
}

func (s *Service) GetAccountManagers(ctx context.Context, user *User) ([]*AccountManager, error) {
	if !s.isAccountManager(user) {
		return nil, NewForbiddenError("user is not account manager")
	}
	logger := logging.GetLogger(ctx)
	accountManagers, err := s.repository.GetAccountManagers(ctx)
	if err != nil {
		logger.Error("failed to get accountManagers from repository", zap.Error(err))
		return nil, err
	}
	return accountManagers, nil
}

func (s *Service) AddAccountManagers(ctx context.Context, user *User, userIDs []uuid.UUID) error {
	if !s.isAccountManager(user) {
		return NewForbiddenError("user is not account manager")
	}
	logger := logging.GetLogger(ctx)
	err := s.repository.AddAccountManagers(ctx, userIDs)
	if err != nil {
		logger.Error("failed to add accountManager in repository", zap.Error(err))
		return err
	}
	return nil
}

func (s *Service) DeleteAccountManagers(
	ctx context.Context, user *User, userIDs []uuid.UUID,
) error {
	if !s.isAccountManager(user) {
		return NewForbiddenError("user is not account manager")
	}
	logger := logging.GetLogger(ctx)
	err := s.repository.DeleteAccountManagers(ctx, userIDs)
	if err != nil {
		logger.Error("failed to delete accountManager in repository", zap.Error(err))
		return err
	}
	return nil
}
