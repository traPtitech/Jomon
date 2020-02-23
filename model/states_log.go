package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

const (
	submitted   int = 1
	fixRequired int = 2
	accepted    int = 3
	fullyRepaid int = 4
	rejected    int = 5
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
	case submitted:
		return json.Marshal("submitted")
	case fixRequired:
		return json.Marshal("fix_required")
	case accepted:
		return json.Marshal("accepted")
	case fullyRepaid:
		return json.Marshal("fully_repaid")
	case rejected:
		return json.Marshal("rejected")
	}
	return nil, errors.New("unknown state type")
}

func GetStateType(str string) (StateType, error) {
	var result StateType
	var err error
	switch str {
	case "submitted":
		result.Type = submitted
	case "fix_required":
		result.Type = fixRequired
	case "accepted":
		result.Type = accepted
	case "fully_repaid":
		result.Type = fullyRepaid
	case "rejected":
		result.Type = rejected
	default:
		err = errors.New("unknown state type")
	}

	return result, err
}

func (st *StatesLog) GiveIsUserAdmin(admins []string) {
	if st == nil {
		return
	}

	st.UpdateUserTrapID.GiveIsUserAdmin(admins)
}

func CreateStatesLog(applicationId uuid.UUID, updateUserTrapId string) (StatesLog, error) {
	log := StatesLog{
		ApplicationID: applicationId,
		UpdateUserTrapID: User{
			TrapId: updateUserTrapId,
		},
		ToState: StateType{Type: 1},
		Reason:  "",
	}

	err := db.Create(&log).Error

	if err != nil {
		return StatesLog{}, err
	}

	return log, nil
}
