package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type GroupRepository interface {
	GetGroups(ctx context.Context) ([]*Group, error)
	CreateGroup(ctx context.Context, name string, description string, budget *int, owners *[]User) (*Group, error)
	GetMembers(ctx context.Context, groupID uuid.UUID) ([]*User, error)
	PostMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) (*User, error)
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
