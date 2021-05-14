package model

import (
	"time"

	"github.com/google/uuid"
)

type GroupBudgetRepository interface {
}

type GroupBudget struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	GroupID   uuid.UUID `gorm:"type:char(36);not null;index"`
	Group     *Group    `gorm:"foreignKey:GroupID"`
	Amount    int       `gorm:"type:int(11);not null"`
	CreatedAt time.Time
}
