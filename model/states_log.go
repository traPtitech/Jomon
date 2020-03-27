package model

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"time"

	"github.com/gofrs/uuid"
)

type StateRepository interface {
	CreateStatesLog(applicationId uuid.UUID, updateUserTrapId string, reason string, toState StateType) (StatesLog, error)
}

type stateRepository struct{}

func NewStateRepository() StateRepository {
	return &stateRepository{}
}

const (
	Submitted   int = 1
	FixRequired int = 2
	Accepted    int = 3
	FullyRepaid int = 4
	Rejected    int = 5
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
	Type int `gorm:"column:to_state;type:tinyint(4);not null;default:0"`
}

func (ty StateType) MarshalJSON() ([]byte, error) {
	switch ty.Type {
	case Submitted:
		return json.Marshal("submitted")
	case FixRequired:
		return json.Marshal("fix_required")
	case Accepted:
		return json.Marshal("accepted")
	case FullyRepaid:
		return json.Marshal("fully_repaid")
	case Rejected:
		return json.Marshal("rejected")
	}
	return nil, errors.New("unknown state type")
}

func GetStateType(str string) (StateType, error) {
	var result StateType
	var err error
	switch str {
	case "submitted":
		result.Type = Submitted
	case "fix_required":
		result.Type = FixRequired
	case "accepted":
		result.Type = Accepted
	case "fully_repaid":
		result.Type = FullyRepaid
	case "rejected":
		result.Type = Rejected
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

func (_ *stateRepository) CreateStatesLog(applicationId uuid.UUID, updateUserTrapId string, reason string, toState StateType) (StatesLog, error) {
	log := StatesLog{
		ApplicationID: applicationId,
		UpdateUserTrapID: User{
			TrapId: updateUserTrapId,
		},
		ToState: toState,
		Reason:  reason,
	}

	err := db.Create(&log).Error

	if err != nil {
		return StatesLog{}, err
	}

	return log, nil
}

func (_ *applicationRepository) createStatesLog(db_ *gorm.DB, applicationId uuid.UUID, updateUserTrapId string) (StatesLog, error) {
	log := StatesLog{
		ApplicationID: applicationId,
		UpdateUserTrapID: User{
			TrapId: updateUserTrapId,
		},
		ToState: StateType{Type: 1},
		Reason:  "",
	}

	err := db_.Create(&log).Error

	if err != nil {
		return StatesLog{}, err
	}

	return log, nil
}
