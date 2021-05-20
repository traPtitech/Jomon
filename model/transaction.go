package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionRepository interface {
}

type Transaction struct {
	ID        uuid.UUID
	CreatedAt time.Time
}
