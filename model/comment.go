package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository interface {
}

type Comment struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	RequestID uuid.UUID `gorm:"type:varchar(36);not null;index"`
	Request   *Request  `gorm:"foeignKey:RequestID"`
	CreatedBy string    `gorm:"type:varchar(32);not null"`
	Comment   string    `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
