package model

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/testutil"
)

func SetupEntClient() (*ent.Client, error) {
	// Logging
	entOptions := []ent.Option{ent.Debug()}

	dbUser := testutil.GetEnvOrDefault("MYSQL_USERNAME", "root")
	dbPass := testutil.GetEnvOrDefault("MYSQL_PASSWORD", "password")
	dbHost := testutil.GetEnvOrDefault("MYSQL_HOSTNAME", "db")
	dbName := testutil.GetEnvOrDefault("MYSQL_DATABASE", "jomon")
	dbPort := testutil.GetEnvOrDefault("MYSQL_PORT", "3306")

	dbDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)

	client, err := ent.Open("mysql", dbDsn, entOptions...)
	if err != nil {
		return nil, fmt.Errorf("can't connect to DATABASE: %w", err)
	}

	return client, nil
}
