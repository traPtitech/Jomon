package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"

	"github.com/gofrs/uuid"
)

type RepayUser struct {
	ID                 int        `gorm:"type:int(11) AUTO_INCREMENT;primary_key" json:"-"`
	ApplicationID      uuid.UUID  `gorm:"type:char(36);not null" json:"-"`
	RepaidToUserTrapID User       `gorm:"embedded;embedded_prefix:repaid_to_user_;not null" json:"repaid_to_user"`
	RepaidByUserTrapID *User      `gorm:"embedded;embedded_prefix:repaid_by_user_" json:"repaid_by_user"`
	RepaidAt           *time.Time `gorm:"type:date" json:"repaid_at"`
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
