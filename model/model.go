package model

import (
	"context"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/testutil"
)

func SetupEntClient() (*ent.Client, error) {
	// Logging
	var entOptions []ent.Option
	if os.Getenv("IS_DEBUG_MODE") != "" {
		entOptions = []ent.Option{ent.Debug()}
	} else {
		entOptions = []ent.Option{}
	}
	dbUser := testutil.GetEnvOrDefault("MARIADB_USERNAME", "root")
	dbPass := testutil.GetEnvOrDefault("MARIADB_PASSWORD", "password")
	dbHost := testutil.GetEnvOrDefault("MARIADB_HOSTNAME", "db")
	dbName := testutil.GetEnvOrDefault("MARIADB_DATABASE", "jomon")
	dbPort := testutil.GetEnvOrDefault("MARIADB_PORT", "3306")

	dbDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)

	client, err := ent.Open("mysql", dbDsn, entOptions...)
	if err != nil {
		return nil, fmt.Errorf("can't connect to DATABASE: %w", err)
	}

	if os.Getenv("IS_DEBUG_MODE") != "" {
		if err := client.Debug().Schema.Create(context.Background()); err != nil {
			return nil, err
		}
	} else {
		if err := client.Schema.Create(context.Background()); err != nil {
			return nil, err
		}
	}

	return client, nil
}
