package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type GroupRepository interface {
	GetGroups(ctx context.Context) ([]*Group, error)
	CreateGroup(ctx context.Context, name string, description string, budget *int, owners *[]User) (*Group, error)
	GetOwner(ctx context.Context, GroupID uuid.UUID) ([]*Owner, error)
	CreateOwner(ctx context.Context, GroupID uuid.UUID, OwnerID uuid.UUID) (*Owner, error)
	DeleteOwner(ctx context.Context, GroupID uuid.UUID, OwnerID uuid.UUID) error
}

type Group struct {
	ID          uuid.UUID
	Name        string
	Description string
	Budget      *int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type Owner struct {
	Owner uuid.UUID
}
