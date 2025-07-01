//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"

	"github.com/google/uuid"
)

type Admin struct {
	ID uuid.UUID `json:"id"`
}

type AdminRepository interface {
	GetAdmins(ctx context.Context) ([]*Admin, error)
	AddAdmins(ctx context.Context, userIDs []uuid.UUID) error
	DeleteAdmins(ctx context.Context, userIDs []uuid.UUID) error
}
