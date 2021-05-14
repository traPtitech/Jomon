package model

import (
	"time"

	"github.com/google/uuid"
)

type GroupOwnerRepository interface {
}

type GroupOwner struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	GroupID   uuid.UUID `gorm:"type:char(36);not null;index"`
	Group     *Group    `gorm:"foreignKey:GroupID"`
	Owner     string    `gorm:"varchar(32);not null"`
	CreatedAt time.Time
}
