package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type StatesLog struct {
	ID               int       `gorm:"type:int(11) AUTO_INCREMENT;primary_key" json:"-"`
	ApplicationID    uuid.UUID `gorm:"type:char(36);not null" json:"-"`
	UpdateUserTrapID User      `gorm:"embedded;embedded_prefix:update_user_" json:"update_user"`
	ToState          int       `gorm:"type:tinyint(4);not null;default:0" json:"to_state"`
	Reason           string    `gorm:"type:text;not null" json:"reason"`
	CreatedAt        time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}
