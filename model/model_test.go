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
	dbUser := testutil.GetEnvOrDefault("MARIADB_USERNAME", "root")
	dbPass := testutil.GetEnvOrDefault("MARIADB_PASSWORD", "password")
	dbHost := testutil.GetEnvOrDefault("MARIADB_HOSTNAME", "db")
	dbName := testutil.GetEnvOrDefault("MARIADB_DATABASE", "jomon-test")
	dbPort := testutil.GetEnvOrDefault("MARIADB_PORT", "3306")

	dbDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)

	client := enttest.Open(t, "mysql", dbDsn, entOptions...)

	ctx := context.Background()

	if err := client.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		return nil, err
	}

	return client, nil
}
