package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionDetailRepository interface {
}

type TransactionDetail struct {
	ID        uuid.UUID
	Amount    int
	Target    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
