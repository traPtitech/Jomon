package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type GroupRepository interface {
	GetGroups(ctx context.Context) ([]*Group, error)
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
