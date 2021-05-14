package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionDetailRepository interface {
}

type TransactionDetail struct {
	ID            uuid.UUID    `gorm:"type:char(36);primaryKey"`
	TransactionID uuid.UUID    `gorm:"type:char(36);not null;index"`
	Transaction   *Transaction `gorm:"foeignKey:TransactionID"`
	Amount        int          `gorm:"type:int(11);not null"`
	Target        string       `gorm:"type:varchar(64);not null"`
	RequestID     uuid.UUID    `gorm:"type:char(36);index"`
	Request       *Request
	GroupID       uuid.UUID `gorm:"type:char(36);index"`
	Group         *Group
	CreatedAt     time.Time
}
