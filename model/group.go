package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type GroupRepository interface {
	GetGroups(ctx context.Context) ([]*Group, error)
	CreateGroup(ctx context.Context, name string, description string, budget *int, owners *[]User) (*Group, error)
	GetOwners(ctx context.Context, GroupID uuid.UUID) (*Owners, error)
	CreateOwners(ctx context.Context, GroupID uuid.UUID, OwnerID uuid.UUID) (*Owners, error)
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

type Owners struct {
	Owners uuid.UUID `json:"owners"`
}
