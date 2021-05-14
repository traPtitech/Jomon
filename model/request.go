package model

import (
	"time"

	"github.com/google/uuid"
)

type RequestRepository interface {
}

type Request struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	CreatedBy uuid.UUID `gorm:"type:varchar(36);not null"`
	Amount    int       `gorm:"type:int(11);not null"`
	CreatedAt time.Time
}
