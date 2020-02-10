package model

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type Application struct {
	ID                       uuid.UUID            `gorm:"type:char(36);primary_key" json:"application_id"`
	LatestApplicationsDetail ApplicationsDetail   `gorm:"foreignkey:ApplicationsDetailsID" json:"current_detail"`
	ApplicationsDetailsID    int                  `gorm:"type:int(11);not null" json:"-"`
	LatestStatesLog          StatesLog            `gorm:"foreignkey:StatesLogsID" json:"-"`
	LatestStatus             StateType            `gorm:"-" json:"current_state"`
	StatesLogsID             int                  `gorm:"type:int(11);not null" json:"-"`
	CreateUserTrapID         User                 `gorm:"embedded;embedded_prefix:create_user_;" json:"applicant"`
	CreatedAt                time.Time            `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	ApplicationsDetails      []ApplicationsDetail `json:"application_detail_logs,omitempty"`
	StatesLogs               []StatesLog          `json:"state_logs,omitempty"`
	ApplicationsImages       []ApplicationsImage  `json:"images,omitempty"`
	Comments                 []Comment            `json:"comments,omitempty"`
	RepayUsers               []RepayUser          `json:"repayment_logs,omitempty" `
}

func (app Application) MarshalJSON() ([]byte, error) {
	app.LatestStatus = app.LatestStatesLog.ToState
	return json.Marshal(app)
}
