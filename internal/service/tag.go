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

type Tag struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt nulltime.NullTime
}

type TagRepository interface {
	GetTags(ctx context.Context) ([]*Tag, error)
	GetTag(ctx context.Context, tagID uuid.UUID) (*Tag, error)
	CreateTag(ctx context.Context, name string) (*Tag, error)
	UpdateTag(ctx context.Context, tagID uuid.UUID, name string) (*Tag, error)
	DeleteTag(ctx context.Context, tagID uuid.UUID) error
}

func (s *Service) GetTags(ctx context.Context) ([]*Tag, error) {
	logger := logging.GetLogger(ctx)
	tags, err := s.repository.GetTags(ctx)
	if err != nil {
		logger.Error("failed to get tags from repository", zap.Error(err))
		return nil, err
	}
	return tags, nil
}

func (s *Service) GetTag(ctx context.Context, tagID uuid.UUID) (*Tag, error) {
	logger := logging.GetLogger(ctx)
	tag, err := s.repository.GetTag(ctx, tagID)
	if err != nil {
		logger.Error("failed to get tag from repository", zap.Error(err))
		return nil, err
	}
	return tag, nil
}

func (s *Service) CreateTag(ctx context.Context, name string) (*Tag, error) {
	logger := logging.GetLogger(ctx)
	tag, err := s.repository.CreateTag(ctx, name)
	if err != nil {
		logger.Error("failed to create tag in repository", zap.Error(err))
		return nil, err
	}
	return tag, nil
}

func (s *Service) UpdateTag(ctx context.Context, tagID uuid.UUID, name string) (*Tag, error) {
	logger := logging.GetLogger(ctx)
	tag, err := s.repository.UpdateTag(ctx, tagID, name)
	if err != nil {
		logger.Error("failed to update tag in repository", zap.Error(err))
		return nil, err
	}
	return tag, nil
}

func (s *Service) DeleteTag(ctx context.Context, tagID uuid.UUID) error {
	logger := logging.GetLogger(ctx)
	if err := s.repository.DeleteTag(ctx, tagID); err != nil {
		logger.Error("failed to delete tag in repository", zap.Error(err))
		return err
	}
	return nil
}
