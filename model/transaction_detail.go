//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionDetailRepository interface {
}

type TransactionDetail struct {
	ID        uuid.UUID
	Title     string
	Amount    int
	Target    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
