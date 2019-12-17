package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Comment struct {
	ID           int       `gorm:"type:int(11);not null;primary_key;AUTO_INCREMENT"`
	AppicationID uuid.UUID `gorm:"type:char(36);not null;index"`
	UserTrapID   string    `gorm:"type:varchar(32);not null;index"`
	Comment      string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
