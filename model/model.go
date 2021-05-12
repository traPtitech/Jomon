package model

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func EstablishConnection() (*gorm.DB, error) {
	dbUser := os.Getenv("MYSQL_USERNAME")
	if dbUser == "" {
		dbUser = "root"
	}

	dbPass := os.Getenv("MYSQL_PASSWORD")
	if dbPass == "" {
		dbPass = "password"
	}

	dbHost := os.Getenv("MYSQL_HOSTNAME")
	if dbHost == "" {
		dbHost = "db"
	}

	dbPort := os.Getenv("MYSQL_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	dbName := os.Getenv("MYSQL_DATABASE")
	if dbName == "" {
		dbName = "jomon"
	}

	config := &gorm.Config{}
	if os.Getenv("GORM_DEBUG") != "" {
		config = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("can't connect to DATABASE: %w", err)
	}

	return db, nil
}

var AllTables = []interface{}{
	&Administrator{},
	&Transaction{},
	&TransactionDetail{},
	&TransactionTag{},
	&Request{},
	&RequestStatus{},
	&RequestTarget{},
	&RequestTag{},
	&RequestFile{},
	&File{},
	&Comment{},
	&Group{},
	&GroupBudget{},
	&GroupUser{},
	&GroupOwner{},
	&Tag{},
}

func Migrate(db *gorm.DB) error {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(AllTables...)
	if err != nil {
		return err
	}

	return nil
}
