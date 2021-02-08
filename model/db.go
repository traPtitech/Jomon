package model

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

var allTables = []interface{}{
	Administrator{},
	Request{},
	RequestTarget{},
	File{},
	RequestStatus{},
	Comment{},
}

// EstablishConnection Establish Connection
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
	_db.BlockGlobalUpdate(true)
	db = _db
	return db, err
}

// Migrate Migration
func Migrate() error {
	if err := db.AutoMigrate(allTables...).Error; err != nil {
		return err
	}

	db.Model(&RequestTarget{}).AddForeignKey("application_id", "applications(id)", "RESTRICT", "RESTRICT")
	db.Model(&File{}).AddForeignKey("application_id", "applications(id)", "RESTRICT", "RESTRICT")
	db.Model(&RequestStatus{}).AddForeignKey("application_id", "applications(id)", "RESTRICT", "RESTRICT")
	db.Model(&Comment{}).AddForeignKey("application_id", "applications(id)", "RESTRICT", "RESTRICT")

	return nil
}
