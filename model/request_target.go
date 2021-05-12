package model

import (
	"time"

	"github.com/google/uuid"
)

type RequestTarget struct {
	ID        uuid.UUID  `gorm:"type:char(36);primaryKey"`
	RequestID string     `gorm:"type:char(36);not null"`
	Request   *Request   `gorm:"foeignKey:RequestID"`
	Target    string     `gorm:"type:varchar(64);not null"`
	PaidAt    *time.Time `gorm:"type:date"`
	CreatedAt time.Time
}
