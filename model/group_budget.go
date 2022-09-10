//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"time"

	"github.com/google/uuid"
)

type GroupBudgetRepository interface {
}

type GroupBudget struct {
	ID        uuid.UUID
	Amount    int
	Comment   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
