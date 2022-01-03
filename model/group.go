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
	CreateGroup(ctx context.Context, name string, description string, budget *int, owners *[]User) (*Group, error)
	GetOwners(ctx context.Context, groupID uuid.UUID) ([]*Owner, error)
	CreateOwner(ctx context.Context, groupID uuid.UUID, ownerID uuid.UUID) (*Owner, error)
	DeleteOwner(ctx context.Context, groupID uuid.UUID, ownerID uuid.UUID) error
	GetMembers(ctx context.Context, groupID uuid.UUID) ([]*User, error)
	CreateMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) (*Member, error)
	DeleteMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) error
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
