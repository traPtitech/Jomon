package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionDetail struct {
	ID            uuid.UUID `gorm:"type:char(36);primaryKey"`
	TransactionID uuid.UUID `gorm:"type:char(36);not null;index`
	CreatedAt     time.Time
}
