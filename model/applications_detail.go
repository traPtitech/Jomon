package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type ApplicationsDetail struct {
	ID               int       `gorm:"type:int(11);not null;primary_key;AUTO_INCREMENT"`
	ApplicationID    uuid.UUID `gorm:"type:char(36);not null"`
	UpdateUserTrapID string    `gorm:"type:varchar(32);not null;index"`
	Type             int       `gorm:"type:tinyint(4);not null"`
	Title            string    `gorm:"type:text;not null"`
	Remarks          string    `gorm:"type:text"`
	Amount           string    `gorm:"type:int(11);not null"`
	PaidAt           time.Time `gorm:"type:timestamp"`
	CreatedAt        time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}
