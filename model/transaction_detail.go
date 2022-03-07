//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
)

type TransactionDetailRepository interface {
	CreateTransactionDetail(ctx context.Context, amount int, target string) (*ent.TransactionDetail, error)
}

type TransactionDetail struct {
	ID        uuid.UUID
	Amount    int
	Target    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
