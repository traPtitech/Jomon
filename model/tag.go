package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type TagRepository interface {
	GetTags(ctx context.Context) ([]*Tag, error)
	GetTag(ctx context.Context, tagID uuid.UUID) (*Tag, error)
	CreateTag(ctx context.Context, name string, description string) (*Tag, error)
	UpdateTag(ctx context.Context, tagID uuid.UUID, name string, description string) (*Tag, error)
	DeleteTag(ctx context.Context, tagID uuid.UUID) error
}
