package model

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/gofrs/uuid"
)

// Status 依頼の状態
type Status string

// 依頼の状態のenum
const (
	Submitted   Status = "submitted"
	FixRequired Status = "fix_required"
	Accepted    Status = "accepted"
	FullyRepaid Status = "fully_repaid"
	Rejected    Status = "rejected"
)

// RequestStatus 依頼の状態
type RequestStatus struct {
	ID        int       `gorm:"type:int(11) AUTO_INCREMENT;primary_key" json:"-"`
	RequestID uuid.UUID `gorm:"type:char(36);not null" json:"-"`
	CreatedBy TrapUser  `gorm:"embedded;embedded_prefix:created_by_" json:"created_by_"`
	Status    Status    `gorm:"embedded" json:"status"`
	Reason    string    `gorm:"type:text;not null" json:"reason"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// GiveIsUserAdmin check whether create user is admin or not
func (rs *RequestStatus) GiveIsUserAdmin(admins []string) {
	if rs == nil {
		return
	}

	rs.CreatedBy.GiveIsUserAdmin(admins)
}

func (repo *requestRepository) UpdateRequestStatus(requestID uuid.UUID, createdBy string, reason string, status Status) (RequestStatus, error) {
	rs := RequestStatus{
		RequestID: requestID,
		CreatedBy: TrapUser{
			TrapID: createdBy,
		},
		Status: status,
		Reason: reason,
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		return repo.updateRequestStatusTransaction(tx, requestID, rs)
	})
	if err != nil {
		return RequestStatus{}, err
	}

	return rs, nil
}

func (*requestRepository) updateRequestStatusTransaction(db *gorm.DB, requestID uuid.UUID, rs RequestStatus) error {
	if err := db.Create(&rs).Error; err != nil {
		return err
	}

	return db.Model(&Request{ID: requestID}).Updates(Request{
		RequestStatusID: rs.ID,
	}).Error
}

func (*requestRepository) createRequestStatus(db *gorm.DB, requestID uuid.UUID, createdBy string) (RequestStatus, error) {
	rs := RequestStatus{
		RequestID: requestID,
		CreatedBy: TrapUser{
			TrapID: createdBy,
		},
		Status: Submitted,
		Reason: "",
	}

	err := db.Create(&rs).Error

	if err != nil {
		return RequestStatus{}, err
	}

	return rs, nil
}
