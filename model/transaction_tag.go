package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionTagRepository interface {
}

type TransactionTag struct {
	ID            uuid.UUID    `gorm:"type:char(36);primaryKey"`
	TransactionID uuid.UUID    `gorm:"type:char(36);not null;index"`
	Transaction   *Transaction `gorm:"foeignKey:TransactionID"`
	TagID         uuid.UUID    `gorm:"type:char(36);not null;index"`
	Tag           *Tag
	CreatedAt     time.Time
}
