//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
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
	GetUserByName(ctx context.Context, name string) (*User, error)
	GetUsers(ctx context.Context) ([]*User, error)
	CreateUser(ctx context.Context, name string, dn string, admin bool) (*User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, name string, dn string, admin bool) (*User, error)
}
