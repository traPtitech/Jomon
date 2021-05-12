package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	MimeType  string    `gorm:"type:text;not null"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}
