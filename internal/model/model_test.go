package model

// FIXME: package model_test にする

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/traPtitech/Jomon/internal/ent"
	"github.com/traPtitech/Jomon/internal/ent/enttest"
	"github.com/traPtitech/Jomon/internal/testutil"
)

func SetupTestEntClient(t *testing.T, ctx context.Context, dbName string) (*ent.Client, error) {
	t.Helper()
	entOptions := []enttest.Option{
		enttest.WithOptions(ent.Log(t.Log)),
	}
	// FIXME: model.go にある `Connect` を流用したい
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

	// nolint:contextcheck
	client := enttest.Open(t, "mysql", dsn, entOptions...).Debug()

	// model/ ディレクトリをPWDとしてテストが実行されるため, migrations ディレクトリのパスを揃える
	if err := MigrateApply(ctx, client, MigrationsDir("../../migrations")); err != nil {
		return nil, fmt.Errorf("failed to apply migration: %w", err)
	}

	return client, nil
}
