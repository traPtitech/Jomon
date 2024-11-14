//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	DisplayName string     `json:"display_name"`
	Admin       bool       `json:"admin"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
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
