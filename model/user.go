//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/service"
)

type User struct {
	ID          uuid.UUID
	Name        string
	DisplayName string
	Admin       bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type UserRepository interface {
	CreateUser(ctx context.Context, name string, dn string, admin bool) (*User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	GetUserByName(ctx context.Context, name string) (*User, error)
	GetUsers(ctx context.Context) ([]*User, error)
	UpdateUser(
		ctx context.Context, userID uuid.UUID, name string, dn string, admin bool,
	) (*User, error)
}
