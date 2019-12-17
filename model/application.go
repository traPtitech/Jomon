package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Application struct {
	ID                   uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	AppicationsDetailsID int       `gorm:"type:int(11);not null;index"`
	StatesLogsID         int       `gorm:"type:int(11);not null;index"`
	CreateUserTrapID     string    `gorm:"type:varchar(32);not null;index"`
	CreatedAt            time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}
