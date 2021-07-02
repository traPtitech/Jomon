package model

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/enttest"
	"github.com/traPtitech/Jomon/ent/migrate"
	"github.com/traPtitech/Jomon/testutil"
)

func SetupTestEntClient(t *testing.T) (*ent.Client, error) {
	entOptions := []enttest.Option{
		enttest.WithOptions(ent.Log(t.Log)),
	}
	// checks
	dbUser := testutil.GetEnvOrDefault("MARIADB_USERNAME", "root")
	dbPass := testutil.GetEnvOrDefault("MARIADB_PASSWORD", "password")
	dbHost := testutil.GetEnvOrDefault("MARIADB_HOSTNAME", "test-db")
	dbName := testutil.GetEnvOrDefault("MARIADB_DATABASE", "jomon")
	dbPort := testutil.GetEnvOrDefault("MARIADB_PORT", "3307")

	dbDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)

	client := enttest.Open(t, "mysql", dbDsn, entOptions...)

	if err := client.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		return nil, err
	}

	return client, nil
}
