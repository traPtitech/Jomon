//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
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
	AddAdmin(ctx context.Context, userID uuid.UUID) (*Admin, error)
	DeleteAdmin(ctx context.Context, userID uuid.UUID) error
}
