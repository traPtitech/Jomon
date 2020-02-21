package model

type Administrator struct {
	TrapID string `gorm:"type:varchar(32);primary_key"`
}

func GetAdministratorList() ([]string, error) {
	var admin []string
	err := db.Model(&Administrator{}).Pluck("trap_id", &admin).Error
	if err != nil {
		return nil, err
	}

	return admin, nil
}
