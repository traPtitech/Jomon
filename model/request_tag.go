package model

import (
	"time"

	"github.com/google/uuid"
)

type RequestTagRepository interface {
}

type RequestTag struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	RequestID uuid.UUID `gorm:"type:char(36);not null;index"`
	Request   *Request  `gorm:"foeignKey:RequestID"`
	TagID     uuid.UUID `gorm:"type:char(36);index"`
	Tag       *Tag      `gorm:"foeignKey:TagID"`
	CreatedAt time.Time
}
