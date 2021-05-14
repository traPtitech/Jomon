package model

import "github.com/google/uuid"

type AdministratorRepository interface {
}

type Administrator struct {
	TrapID uuid.UUID `gorm:"type:varchar(32);not null;primaryKey"`
}
