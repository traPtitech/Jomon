package model

import "github.com/jinzhu/gorm"

type Administrator struct {
	TrapID string `gorm:"type:varchar(32);primary_key"`
}

type AdministratorRepository interface {
	IsAdmin(userId string) (bool, error)
	GetAdministratorList() ([]string, error)
	AddAdministrator(userId string) error
	RemoveAdministrator(userId string) error
}

type administratorRepository struct{}

func NewAdministratorRepository() AdministratorRepository {
	return &administratorRepository{}
}

func (_ administratorRepository) IsAdmin(userId string) (bool, error) {
	var tmp *Administrator
	err := db.Where(&Administrator{TrapID: userId}).First(tmp).Error
	if err != nil {
		return true, nil
	} else if gorm.IsRecordNotFoundError(err) {
		return false, nil
	} else {
		return false, err
	}
}

func (_ administratorRepository) GetAdministratorList() ([]string, error) {
	var admin []string

	err := db.Model(&Administrator{}).Pluck("trap_id", &admin).Error
	if err != nil {
		return nil, err
	}

	return admin, nil
}

func (_ administratorRepository) AddAdministrator(userId string) error {
	admin := Administrator{TrapID: userId}
	return db.Create(&admin).Error
}

func (_ administratorRepository) RemoveAdministrator(userId string) error {
	admin := Administrator{TrapID: userId}
	return db.Delete(&admin).Error
}
