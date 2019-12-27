package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Application struct {
	ID                       uuid.UUID `gorm:"type:char(36);primary_key"`
	LatestApplicationsDetail ApplicationsDetail
	ApplicationsDetailsID    int `gorm:"type:int(11);not null"`
	LatestStatesLog          StatesLog
	StatesLogsID             int       `gorm:"type:int(11);not null"`
	CreateUserTrapID         string    `gorm:"type:varchar(32);not null;index"`
	CreatedAt                time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	ApplicationsDetails      []ApplicationsDetail
	StatesLogs               []StatesLog
	ApplicationsImages       []ApplicationsImage
	Comments                 []Comment
	RepayUsers               []RepayUser
}
