package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

type StatesLog struct {
	ID               int       `gorm:"type:int(11) AUTO_INCREMENT;primary_key" json:"-"`
	ApplicationID    uuid.UUID `gorm:"type:char(36);not null" json:"-"`
	UpdateUserTrapID User      `gorm:"embedded;embedded_prefix:update_user_" json:"update_user"`
	ToState          StateType `gorm:"embedded" json:"to_state"`
	Reason           string    `gorm:"type:text;not null" json:"reason"`
	CreatedAt        time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type StateType struct {
	Type int `gorm:"type:tinyint(4);not null;default:0"`
}

func (ty StateType) MarshalJSON() ([]byte, error) {
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
	return nil, errors.New("unknown state type")
}
