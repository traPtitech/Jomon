package model

import (
	"time"

	"github.com/google/uuid"
)

type RequestFileRepository interface {
}

type RequestFile struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	RequestID uuid.UUID `gorm:"type:char(36);not null;index"`
	Request   *Request  `gorm:"foeignKey:RequestID"`
	FileID    uuid.UUID `gorm:"type:char(36);not null"`
	File      *File     `gorm:"foeignKey:FileID"`
	CreatedAt time.Time
}
