package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name        string    `gorm:"type:varchar(64);not null"`
	Description string    `gorm:"type:text;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}
