package model

import (
	"time"

	"github.com/google/uuid"
)

type GroupUserRepository interface {
}

type GroupUser struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	GroupID   uuid.UUID `gorm:"type:char(36);not null;index"`
	Group     *Group    `gorm:"foreignKey:GroupID"`
	UserID    string    `gorm:"type:varchar(32);not null"`
	CreatedAt time.Time
}
