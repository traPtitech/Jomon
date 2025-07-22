//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/service"
)

type Tag struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type TagRepository interface {
	GetTags(ctx context.Context) ([]*Tag, error)
	GetTag(ctx context.Context, tagID uuid.UUID) (*Tag, error)
	CreateTag(ctx context.Context, name string) (*Tag, error)
	UpdateTag(ctx context.Context, tagID uuid.UUID, name string) (*Tag, error)
	DeleteTag(ctx context.Context, tagID uuid.UUID) error
}
