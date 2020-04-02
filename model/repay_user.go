package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
	"errors"

	"github.com/gofrs/uuid"
)

var (
    ErrAlreadyRepaid = errors.New("alreadyRepaid")
)

type RepayUser struct {
	ID                 int        `gorm:"type:int(11) AUTO_INCREMENT;primary_key" json:"-"`
	ApplicationID      uuid.UUID  `gorm:"type:char(36);not null" json:"-"`
	RepaidToUserTrapID User       `gorm:"embedded;embedded_prefix:repaid_to_user_;not null" json:"repaid_to_user"`
	RepaidByUserTrapID *User      `gorm:"embedded;embedded_prefix:repaid_by_user_" json:"repaid_by_user"`
	RepaidAt           *time.Time `gorm:"type:timestamp" json:"repaid_at"`
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

func (repo *applicationRepository) UpdateRepayUser(applicationId uuid.UUID, repaidToUserTrapID string, repaidByUserTrapID string) (RepayUser, bool, error) {
	dt := time.Now()
	ru := RepayUser{
		ApplicationID: applicationId,
		RepaidToUserTrapID: User{
			TrapId: repaidToUserTrapID,
		},
		RepaidByUserTrapID: &User{
			TrapId: repaidByUserTrapID,
		},
		RepaidAt: &dt,
	}
	var repaidUser RepayUser
	err := db.Where("ApplicationID = ?", applicationId).Where("RepaidToUserTrapID = ?", repaidToUserTrapID).First(&repaidUser).Error
	if err != nil {
		return RepayUser{}, false, err
	}
	if repaidUser.RepaidByUserTrapID != nil || repaidUser.RepaidAt != nil {
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
		var rus []RepayUser
		if err := db.Where("ApplicationID = ?", applicationId).Find(&rus).Error; err != nil {
			return err
		}
		for _, user := range rus {
			allUsersRepaidCheck = allUsersRepaidCheck && (user.RepaidByUserTrapID != nil) && (user.RepaidAt != nil)
		}
		if allUsersRepaidCheck {
			err := repo.updateStatesLogTransaction(tx, applicationId, log)
			if err != nil {
				return err
			}
		}
		return tx.Model(&RepayUser{ApplicationID: applicationId, RepaidToUserTrapID: User{TrapId: repaidToUserTrapID}}).Updates(RepayUser{
			RepaidByUserTrapID: &User{
				TrapId: repaidByUserTrapID,
			},
		}).Error
	})
	if err != nil {
		return RepayUser{}, false, err
	}

	return ru, allUsersRepaidCheck, nil
}