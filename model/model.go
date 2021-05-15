package model

import (
	"fmt"
	"os"

	"github.com/traPtitech/Jomon/ent"

	"github.com/go-sql-driver/mysql"
)

func SetupEntClient() (*ent.Client, error) {
	entOptions := []ent.Option{}

	// 発行されるSQLをロギングするなら
	entOptions = append(entOptions, ent.Debug())
	dbUser := os.Getenv("MYSQL_USERNAME")
	if dbUser == "" {
		dbUser = "root"
	}

	dbPass := os.Getenv("MYSQL_PASSWORD")
	if dbPass == "" {
		dbPass = "root"
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
		dbName = "test_database"
	}
	mc := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPass,
		Net:                  "tcp",
		Addr:                 "localhost" + ":" + dbPort,
		DBName:               dbName,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	client, err := ent.Open("mysql", mc.FormatDSN(), entOptions...)
	if err != nil {
		return nil, fmt.Errorf("can't connect to DATABASE: %w", err)
	}

	return client, nil
}
