//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"

	"github.com/google/uuid"
)

type AccountManager struct {
	ID uuid.UUID `json:"id"`
}

type AccountManagerRepository interface {
	GetAccountManagers(ctx context.Context) ([]*AccountManager, error)
	AddAccountManagers(ctx context.Context, userIDs []uuid.UUID) error
	DeleteAccountManagers(ctx context.Context, userIDs []uuid.UUID) error
}
