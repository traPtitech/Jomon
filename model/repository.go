package model

import (
	"gorm.io/gorm"
)

type Repository interface {
}

type GormRepository struct {
	db *gorm.DB
}
