package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type RepayUser struct {
	ID                 int        `gorm:"type:int(11) AUTO_INCREMENT;primary_key" json:"-"`
	ApplicationID      uuid.UUID  `gorm:"type:char(36);not null" json:"-"`
	RepaidToUserTrapID User       `gorm:"embedded;embedded_prefix:repaid_to_user_" json:"repaid_to_user"`
	RepaidByUserTrapID User       `gorm:"embedded;embedded_prefix:repaid_by_user_" json:"repaid_by_user"`
	RepaidAt           *time.Time `gorm:"type:timestamp;null;" json:"repaid_at"`
}
