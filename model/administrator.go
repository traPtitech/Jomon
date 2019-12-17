package model

type Administrator struct {
	TrapID string `gorm:"type:varchar(32);not null;primary_key"`
}
