package model

import (
	"github.com/gofrs/uuid"
)

type ApplicationImage struct {
	ID            uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	ApplicationID uuid.UUID `gorm:"type:char(36);not null;index"`
}
