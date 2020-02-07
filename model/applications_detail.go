package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

type ApplicationsDetail struct {
	ID               int             `gorm:"type:int(11) AUTO_INCREMENT;primary_key" json:"-"`
	ApplicationID    uuid.UUID       `gorm:"type:char(36);not null" json:"-"`
	UpdateUserTrapID User            `gorm:"embedded;embedded_prefix:update_user_" json:"update_user"`
	Type             ApplicationType `gorm:"embedded" json:"type"`
	Title            string          `gorm:"type:text;not null" json:"title"`
	Remarks          string          `gorm:"type:text;not null" json:"remarks"`
	Amount           int             `gorm:"type:int(11);not null" json:"amount"`
	PaidAt           time.Time       `gorm:"type:timestamp" json:"paid_at"`
	CreatedAt        time.Time       `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type ApplicationType struct {
	Type int `gorm:"tinyint(4);not null"`
}

func (ty ApplicationType) MarshalJSON() ([]byte, error) {
	switch ty.Type {
	case 0:
		return json.Marshal("submitted")
	case 1:
		return json.Marshal("fix_required")
	case 2:
		return json.Marshal("accepted")
	case 3:
		return json.Marshal("fully_repaid")
	case 4:
		return json.Marshal("rejected")
	}
	return nil, errors.New("unknown application type")
}
