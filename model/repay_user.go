package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type RepayUser struct {
	ID                 int       `gorm:"type:int(11);not null;primary_key;AUTO_INCREMENT"`
	ApplicationID      uuid.UUID `gorm:"type:char(36);not null"`
	RepaidToUserTrapID string    `gorm:"type:varchar(32);not null;index"`
	RepaidByUserTrapID string    `gorm:"type:varchar(32);null;index"`
	RepaidAt           *time.Time
}
