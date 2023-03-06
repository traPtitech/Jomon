//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type GroupRepository interface {
	GetGroups(ctx context.Context) ([]*Group, error)
	GetGroup(ctx context.Context, groupID uuid.UUID) (*Group, error)
	CreateGroup(ctx context.Context, name string, description string, budget *int) (*Group, error)
	UpdateGroup(ctx context.Context, groupID uuid.UUID, name string, description string, budget *int) (*Group, error)
	DeleteGroup(ctx context.Context, groupID uuid.UUID) error
	GetOwners(ctx context.Context, groupID uuid.UUID) ([]*Owner, error)
	AddOwners(ctx context.Context, groupID uuid.UUID, ownerIDs []uuid.UUID) ([]*Owner, error)
	DeleteOwners(ctx context.Context, groupID uuid.UUID, ownerIDs []uuid.UUID) error
	GetMembers(ctx context.Context, groupID uuid.UUID) ([]*Member, error)
	AddMembers(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) ([]*Member, error)
	DeleteMembers(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) error
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
	ID uuid.UUID
}

type Member struct {
	ID uuid.UUID
}
