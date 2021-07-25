package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Name        string
	DisplayName string
	Admin       bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type UserRepository interface {
	GetMe(ctx context.Context, name string) (*User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
}
