//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/internal/logging"
	"github.com/traPtitech/Jomon/internal/nulltime"
	"go.uber.org/zap"
)

type User struct {
	ID             uuid.UUID
	Name           string
	DisplayName    string
	AccountManager bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      nulltime.NullTime
}

type UserRepository interface {
	CreateUser(ctx context.Context, name string, dn string, accountManager bool) (*User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	GetUserByName(ctx context.Context, name string) (*User, error)
	GetUsers(ctx context.Context) ([]*User, error)
	UpdateUser(
		ctx context.Context, userID uuid.UUID, name string, dn string, accountManager bool,
	) (*User, error)
}

func (s *Service) GetUsers(ctx context.Context) ([]*User, error) {
	logger := logging.GetLogger(ctx)
	users, err := s.repository.GetUsers(ctx)
	if err != nil {
		logger.Error("failed to get users from repository", zap.Error(err))
		return nil, err
	}
	return users, nil
}

func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	logger := logging.GetLogger(ctx)
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		logger.Error("failed to get user by id from repository", zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (s *Service) GetUserByName(ctx context.Context, name string) (*User, error) {
	logger := logging.GetLogger(ctx)
	user, err := s.repository.GetUserByName(ctx, name)
	if err != nil {
		logger.Error("failed to get user by name from repository", zap.Error(err))
		return nil, err
	}
	return user, nil
}

type CreateUserInputs struct {
	Name        string
	DisplayName string
}

func (s *Service) CreateUser(ctx context.Context, inputs CreateUserInputs) (*User, error) {
	logger := logging.GetLogger(ctx)
	user, err := s.repository.CreateUser(
		ctx, inputs.Name, inputs.DisplayName, false,
	)
	if err != nil {
		logger.Error("failed to create user in repository", zap.Error(err))
		return nil, err
	}
	return user, nil
}

type UpdateUserInputs struct {
	Name           string
	DisplayName    string
	AccountManager bool
}

func (s *Service) UpdateUser(
	ctx context.Context, userID uuid.UUID, inputs UpdateUserInputs,
) (*User, error) {
	logger := logging.GetLogger(ctx)
	user, err := s.repository.UpdateUser(
		ctx, userID, inputs.Name, inputs.DisplayName, inputs.AccountManager,
	)
	if err != nil {
		logger.Error("failed to update user in repository", zap.Error(err))
		return nil, err
	}
	return user, nil
}
