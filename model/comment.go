package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Comment struct {
	ID            int       `gorm:"type:int(11) AUTO_INCREMENT;primary_key" json:"comment_id"`
	ApplicationID uuid.UUID `gorm:"type:char(36);not null" json:"-"`
	UserTrapID    User    `gorm:"embedded;embedded_prefix:user_" json:"user"`
	Comment       string    `gorm:"type:text;not null" json:"comment"`
	CreatedAt     time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"type:timestamp;not null;" json:"-"`
}
