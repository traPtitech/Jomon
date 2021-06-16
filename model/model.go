package model

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/testutil"
)

func SetupEntClient() (*ent.Client, error) {
	// Logging
	entOptions := []ent.Option{}

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

	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, err
	}

	return client, nil
}
