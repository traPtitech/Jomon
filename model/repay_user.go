package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/gofrs/uuid"
)

var (
	ErrAlreadyRepaid = errors.New("alreadyRepaid")
)

type RepayUser struct {
	ID                 int        `gorm:"type:int(11) AUTO_INCREMENT;primary_key" json:"-"`
	ApplicationID      uuid.UUID  `gorm:"type:char(36);not null" json:"-"`
	RepaidToUserTrapID User       `gorm:"embedded;embedded_prefix:repaid_to_user_" json:"repaid_to_user"`
	RepaidByUserTrapID User       `gorm:"embedded;embedded_prefix:repaid_by_user_" json:"repaid_by_user"`
	RepaidAt           *time.Time `gorm:"type:date" json:"repaid_at"`
	CreatedAt          time.Time  `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (ru *RepayUser) GiveIsUserAdmin(admins []string) {
	if ru == nil {
		return
	}

	ru.RepaidToUserTrapID.GiveIsUserAdmin(admins)
	ru.RepaidByUserTrapID.GiveIsUserAdmin(admins)
}

func (_ *applicationRepository) createRepayUser(db_ *gorm.DB, applicationId uuid.UUID, repaidToUserTrapID string) error {
	ru := RepayUser{
		ApplicationID: applicationId,
		RepaidToUserTrapID: User{
			TrapId: repaidToUserTrapID,
		},
	}

	if !db_.Where(&ru).First(&RepayUser{}).RecordNotFound() {
		return fmt.Errorf("already exists")
	}

	return db_.Create(&ru).Error
}

func (_ *applicationRepository) deleteRepayUserByApplicationID(db_ *gorm.DB, applicationId uuid.UUID) error {
	return db_.Where(&RepayUser{
		ApplicationID: applicationId,
	}).Delete(&RepayUser{}).Error
}

func (repo *applicationRepository) UpdateRepayUser(applicationId uuid.UUID, repaidToUserTrapID string, repaidByUserTrapID string, repaidAt time.Time) (RepayUser, bool, error) {
	dt := time.Now()
	ru := RepayUser{
		ApplicationID: applicationId,
		RepaidToUserTrapID: User{
			TrapId: repaidToUserTrapID,
		},
		RepaidByUserTrapID: User{
			TrapId: repaidByUserTrapID,
		},
		RepaidAt: &repaidAt,
	}

	var repaidUser RepayUser
	err := db.Where("application_id = ?", applicationId).Where("repaid_to_user_trap_id = ?", repaidToUserTrapID).First(&repaidUser).Error
	if err != nil {
		return RepayUser{}, false, err
	}
	if repaidUser.RepaidAt != nil {
		return RepayUser{}, false, ErrAlreadyRepaid
	}

	log := StatesLog{
		ApplicationID: applicationId,
		UpdateUserTrapID: User{
			TrapId: repaidByUserTrapID,
		},
		ToState: StateType{Type: FullyRepaid},
		Reason:  "",
	}
	allUsersRepaidCheck := true
	err = db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&RepayUser{}).Where("application_id = ?", applicationId).Where("repaid_to_user_trap_id = ?", repaidToUserTrapID).Update(RepayUser{
			RepaidByUserTrapID: User{
				TrapId: repaidByUserTrapID,
			},
			RepaidAt: &dt,
		}).Error
		if err != nil {
			return err
		}
		var rus []RepayUser
		if err := db.Where("application_id = ?", applicationId).Find(&rus).Error; err != nil {
			return err
		}
		for _, user := range rus {
			if user.RepaidToUserTrapID.TrapId == repaidToUserTrapID { // Transaction内で必ずUpdate
				continue
			}
			allUsersRepaidCheck = allUsersRepaidCheck && user.RepaidAt != nil
		}
		if allUsersRepaidCheck {
			err := repo.updateStatesLogTransaction(tx, applicationId, log)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return RepayUser{}, false, err
	}

	return ru, allUsersRepaidCheck, nil
}
