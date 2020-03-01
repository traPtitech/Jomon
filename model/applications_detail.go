package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"

	"github.com/gofrs/uuid"
)

const (
	Club    int = 1
	Contest int = 2
	Event   int = 3
	Public  int = 4
)

type ApplicationsDetail struct {
	ID               int             `gorm:"type:int(11) AUTO_INCREMENT;primary_key" json:"-"`
	ApplicationID    uuid.UUID       `gorm:"type:char(36);not null" json:"-"`
	UpdateUserTrapID User            `gorm:"embedded;embedded_prefix:update_user_" json:"update_user"`
	Type             ApplicationType `gorm:"embedded" json:"type"`
	Title            string          `gorm:"type:text;not null" json:"title"`
	Remarks          string          `gorm:"type:text;not null" json:"remarks"`
	Amount           int             `gorm:"type:int(11);not null" json:"amount"`
	PaidAt           PaidAt          `gorm:"embedded" json:"paid_at"`
	UpdatedAt        time.Time       `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type ApplicationType struct {
	Type int `gorm:"tinyint(4);not null"`
}

func (ty ApplicationType) MarshalJSON() ([]byte, error) {
	switch ty.Type {
	case Club:
		return json.Marshal("club")
	case Contest:
		return json.Marshal("contest")
	case Event:
		return json.Marshal("event")
	case Public:
		return json.Marshal("public")
	}
	return nil, fmt.Errorf("unknown application type: %d", ty.Type)
}

func (ty *ApplicationType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	typeInt, err := GetApplicationTypeFromString(str)
	if err != nil {
		return err
	}

	ty.Type = typeInt
	return nil
}

type PaidAt struct {
	PaidAt time.Time `gorm:"type:timestamp;not null"`
}

func (p PaidAt) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.PaidAt.Format("2006-01-02"))
}

func GetApplicationTypeFromString(str string) (int, error) {
	switch str {
	case "club":
		return Club, nil
	case "contest":
		return Contest, nil
	case "event":
		return Event, nil
	case "public":
		return Public, nil
	}

	return 0, errors.New("unknown application type")
}

func GetApplicationType(str string) (ApplicationType, error) {
	typeInt, err := GetApplicationTypeFromString(str)
	if err != nil {
		return ApplicationType{}, err
	}

	return ApplicationType{Type: typeInt}, nil
}

func (det *ApplicationsDetail) GiveIsUserAdmin(admins []string) {
	if det == nil {
		return
	}

	det.UpdateUserTrapID.GiveIsUserAdmin(admins)
}

func createApplicationsDetail(db_ *gorm.DB, applicationId uuid.UUID, updateUserTrapID string, typ ApplicationType, title string, remarks string, amount int, paidAt time.Time) (ApplicationsDetail, error) {
	detail := ApplicationsDetail{
		ApplicationID:    applicationId,
		UpdateUserTrapID: User{TrapId: updateUserTrapID},
		Type:             typ,
		Title:            title,
		Remarks:          remarks,
		Amount:           amount,
		PaidAt:           PaidAt{PaidAt: paidAt},
	}

	err := db_.Create(&detail).Error
	if err != nil {
		return ApplicationsDetail{}, err
	}

	return detail, nil
}

func putApplicationsDetail(db_ *gorm.DB, currentDetailId int, updateUserTrapID string, typ *ApplicationType, title string, remarks string, amount *int, paidAt *time.Time) (ApplicationsDetail, error) {
	var detail ApplicationsDetail
	err := db_.Find(&detail, ApplicationsDetail{ID: currentDetailId}).Error
	if err != nil {
		return ApplicationsDetail{}, err
	}

	detail.ID = 0 // zero value of int is 0
	detail.UpdateUserTrapID.TrapId = updateUserTrapID
	if typ != nil {
		detail.Type = *typ
	}

	if title != "" {
		detail.Title = title
	}

	if remarks != "" {
		detail.Remarks = remarks
	}

	if amount != nil {
		detail.Amount = *amount
	}

	if paidAt != nil {
		detail.PaidAt.PaidAt = *paidAt
	}

	err = db_.Create(&detail).Error
	if err != nil {
		return ApplicationsDetail{}, err
	}

	return detail, nil
}
