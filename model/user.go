package model

type User struct {
	TrapId  string `gorm:"type:varchar(32);not null;unique;" json:"trap_id"`
	IsAdmin bool   `gorm:"-" json:"is_admin"`
}
