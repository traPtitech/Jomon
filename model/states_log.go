package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type StatesLog struct {
	ID               int       `gorm:"type:int(11) AUTO_INCREMENT;not null;primary_key"`
	ApplicationID    uuid.UUID `gorm:"type:char(36);not null"`
	UpdateUserTrapID string    `gorm:"type:varchar(32);not null;index"`
	ToState          int       `gorm:"type:tinyint(4);not null;default:0"`
	Reason           string    `gorm:"type:text;null"`
	CreatedAt        time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}
