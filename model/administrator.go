package model

import "github.com/jinzhu/gorm"

// Administrator struct of Administrator
type Administrator struct {
	TrapID string `gorm:"type:varchar(32);primary_key"`
}

// AdministratorRepository Repo of Administrator
type AdministratorRepository interface {
	IsAdmin(userID string) (bool, error)
	GetAdministratorList() ([]string, error)
	AddAdministrator(userID string) error
	RemoveAdministrator(userID string) error
}

type administratorRepository struct{}

// NewAdministratorRepository Make AdministratorRepository
func NewAdministratorRepository() AdministratorRepository {
	return &administratorRepository{}
}

func (administratorRepository) IsAdmin(userID string) (bool, error) {
	ad := &Administrator{TrapID: userID}
	err := db.First(ad).Error
	if err == nil {
		return true, nil
	} else if gorm.IsRecordNotFoundError(err) {
		return false, nil
	} else {
		return false, err
	}
}

func (administratorRepository) GetAdministratorList() ([]string, error) {
	var admin []string

	err := db.Model(&Administrator{}).Pluck("trap_id", &admin).Error
	if err != nil {
		return nil, err
	}

	return admin, nil
}

func (administratorRepository) AddAdministrator(userID string) error {
	admin := Administrator{TrapID: userID}
	return db.FirstOrCreate(&admin).Error
}

func (administratorRepository) RemoveAdministrator(userID string) error {
	admin := Administrator{TrapID: userID}
	return db.Delete(&admin).Error
}
