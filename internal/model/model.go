package model

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/traPtitech/Jomon/internal/ent"
)

func getenvOrDefault(key, fallback string) string {
	// TODO: これは testutil.GetEnvOrDefault のコピペ
	//       適切な場所に動かしたい
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func loadDsn() string {
	username := getenvOrDefault("MARIADB_USERNAME", "root")
	password := getenvOrDefault("MARIADB_PASSWORD", "password")
	host := getenvOrDefault("MARIADB_HOSTNAME", "db")
	port := getenvOrDefault("MARIADB_PORT", "3306")
	database := getenvOrDefault("MARIADB_DATABASE", "jomon")
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, database,
	)
}

func Connect() (*ent.Client, error) {
	dsn := loadDsn()
	entOptions := []ent.Option{}
	if os.Getenv("IS_DEBUG_MODE") != "" {
		entOptions = append(entOptions, ent.Debug())
	}
	client, err := ent.Open("mysql", dsn, entOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return client, nil
}
