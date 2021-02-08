package model

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/gofrs/uuid"
)

// ErrAlreadyPaid すでに取引がされている
// ErrAlreadyExists すでにレコードが存在している
var (
	ErrAlreadyPaid   = errors.New("already paid")
	ErrAlreadyExists = errors.New("already exists")
)

// RequestTarget target
type RequestTarget struct {
	ID        int        `gorm:"type:int(11) AUTO_INCREMENT;primary_key" json:"-"`
	RequestID uuid.UUID  `gorm:"type:char(36);not null" json:"-"`
	Target    string     `gorm:"embedded;embedded_prefix:repaid_to_user_" json:"repaid_to_user"`
	CreatedBy TrapUser   `gorm:"embedded;embedded_prefix:repaid_by_user_" json:"repaid_by_user"`
	PaidAt    *time.Time `gorm:"type:date" json:"repaid_at"`
	CreatedAt time.Time  `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (*requestRepository) createRequestTarget(db *gorm.DB, requestID uuid.UUID, target string) error {
	rt := RequestTarget{
		RequestID: requestID,
		Target:    target,
	}

	if !db.Where(&rt).First(&RequestTarget{}).RecordNotFound() {
		return ErrAlreadyExists
	}

	return db.Create(&rt).Error
}

func (*requestRepository) deleteRequestTargetByRequestID(db *gorm.DB, requestID uuid.UUID) error {
	return db.Where(&RequestTarget{
		RequestID: requestID,
	}).Delete(&RequestTarget{}).Error
}

func (repo *requestRepository) UpdateRequestTarget(requestID uuid.UUID, target string, createdBy string, paidAt time.Time) (RequestTarget, bool, error) {
	rt := RequestTarget{
		RequestID: requestID,
		Target:    target,
		CreatedBy: TrapUser{
			TrapID: createdBy,
		},
		PaidAt: &paidAt,
	}

	var requestTarget RequestTarget
	err := db.Where("request_id = ?", requestID).Where("target = ?", target).First(&requestTarget).Error
	if err != nil {
		return RequestTarget{}, false, err
	}
	if requestTarget.PaidAt != nil {
		return RequestTarget{}, false, ErrAlreadyPaid
	}

	rs := RequestStatus{
		RequestID: requestID,
		CreatedBy: TrapUser{
			TrapID: createdBy,
		},
		Status: FullyRepaid,
		Reason: "",
	}
	allUsersRepaidCheck := true
	err = db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&RequestTarget{}).Where("request_id = ?", requestID).Where("target = ?", target).Update(RequestTarget{
			Target: target,
			PaidAt: &paidAt,
		}).Error
		if err != nil {
			return err
		}
		var rts []RequestTarget
		if err := db.Where("request_id = ?", requestID).Find(&rts).Error; err != nil {
			return err
		}
		for _, user := range rts {
			if user.Target == target { // Transaction内で必ずUpdate
				continue
			}
			allUsersRepaidCheck = allUsersRepaidCheck && user.PaidAt != nil
		}
		if allUsersRepaidCheck {
			err := repo.updateRequestStatusTransaction(tx, requestID, rs)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return RequestTarget{}, false, err
	}
	return rt, allUsersRepaidCheck, nil
}
