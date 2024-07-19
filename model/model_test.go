package model

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/enttest"
	"github.com/traPtitech/Jomon/testutil"
)

func SetupTestEntClient(t *testing.T, dbName string) (*ent.Client, error) {
	t.Helper()
	entOptions := []enttest.Option{
		enttest.WithOptions(ent.Log(t.Log)),
	}
	dbUser := testutil.GetEnvOrDefault("MARIADB_USERNAME", "root")
	dbPass := testutil.GetEnvOrDefault("MARIADB_PASSWORD", "password")
	dbHost := testutil.GetEnvOrDefault("MARIADB_HOSTNAME", "db")
	dbPort := testutil.GetEnvOrDefault("MARIADB_PORT", "3306")

	dbDsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort)
	conn, err := sql.Open("mysql", dbDsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	client := enttest.Open(t, "mysql", dsn, entOptions...).Debug()

	ctx := context.Background()

	if err := client.Schema.Create(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
