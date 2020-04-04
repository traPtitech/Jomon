package model

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"time"

	"github.com/gofrs/uuid"
)

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

func (ty *StateType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	var err error
	switch s {
	case "submitted":
		ty.Type = Submitted
	case "fix_required":
		ty.Type = FixRequired
	case "accepted":
		ty.Type = Accepted
	case "fully_repaid":
		ty.Type = FullyRepaid
	case "rejected":
		ty.Type = Rejected
	default:
		err = errors.New("unknown state type")
		return err
	}
	return nil
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

func (repo *applicationRepository) UpdateStatesLog(applicationId uuid.UUID, updateUserTrapId string, reason string, toState StateType) (StatesLog, error) {
	log := StatesLog{
		ApplicationID: applicationId,
		UpdateUserTrapID: User{
			TrapId: updateUserTrapId,
		},
		ToState: toState,
		Reason:  reason,
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		return repo.updateStatesLogTransaction(tx, applicationId, log)
	})
	if err != nil {
		return StatesLog{}, err
	}

	return log, nil
}

func (_ *applicationRepository) updateStatesLogTransaction(db *gorm.DB, applicationId uuid.UUID, log StatesLog) error {
	if err := db.Create(&log).Error; err != nil {
		return err
	}

	return db.Model(&Application{ID: applicationId}).Updates(Application{
		StatesLogsID: log.ID,
	}).Error
}

func (_ *applicationRepository) createStatesLog(db_ *gorm.DB, applicationId uuid.UUID, updateUserTrapId string) (StatesLog, error) {
	log := StatesLog{
		ApplicationID: applicationId,
		UpdateUserTrapID: User{
			TrapId: updateUserTrapId,
		},
		ToState: StateType{Type: Submitted},
		Reason:  "",
	}

	err := db_.Create(&log).Error

	if err != nil {
		return StatesLog{}, err
	}

	return log, nil
}
