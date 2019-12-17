package model

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

var allTables = []interface{}{
	Administrator{},
	Application{},
	ApplicationsDetail{},
	RepayUser{},
	ApplicationsImage{},
	StatesLog{},
	Comment{},
}

func EstablishConnection() (*gorm.DB, error) {
	user := os.Getenv("MARIADB_USERNAME")
	if user == "" {
		user = "root"
	}

	pass := os.Getenv("MARIADB_PASSWORD")
	if pass == "" {
		pass = "password"
	}

	host := os.Getenv("MARIADB_HOSTNAME")
	if host == "" {
		host = "localhost"
	}

	dbname := os.Getenv("MARIADB_DATABASE")
	if dbname == "" {
		dbname = "jomon"
	}

	_db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", user, pass, host, dbname)+"?charset=utf8mb4&parseTime=True&loc=Local")

	db = _db
	return db, err
}

func Migrate() error {
	if err := db.AutoMigrate(allTables...).Error; err != nil {
		return err
	}

	return nil
}
