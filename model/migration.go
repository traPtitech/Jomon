package model

// - https://entgo.io/docs/versioned-migrations/

import (
	"context"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/traPtitech/Jomon/ent"
)

type MigrateConfig struct {
	// `migrations` ディレクトリへののパス
	migrationsDir string
}

type MigrateOption func(*MigrateConfig)

// `migrations` ディレクトリへのパスを設定します.
//
// このオプションを指定しなかった場合, デフォルトでは `migrations` が使われます.
// `MigrateDiff` ないし `MigrateApply` の実行時にここで指定されたディレクトリが存在している必要があります.
func MigrationsDir(dir string) MigrateOption {
	return func(c *MigrateConfig) {
		c.migrationsDir = dir
	}
}

func defaultMigrateConfig() *MigrateConfig {
	return &MigrateConfig{"migrations"}
}

func (c *MigrateConfig) applyOptions(options ...MigrateOption) {
	for _, o := range options {
		o(c)
	}
}

func defaultMigrateOptions() []schema.MigrateOption {
	return []schema.MigrateOption{
		schema.WithMigrationMode(schema.ModeReplay),
		schema.WithDialect(dialect.MySQL),
		schema.WithFormatter(atlas.DefaultFormatter),
	}
}

// `atlas migrate diff` へのエイリアスです.
//
// `MigrationsDir` のディレクトリを参照してdiffの計算が行われます.
func MigrateDiff(ctx context.Context, client *ent.Client, options ...MigrateOption) error {
	config := defaultMigrateConfig()
	config.applyOptions(options...)

	dir, err := atlas.NewLocalDir(config.migrationsDir)
	if err != nil {
		return err
	}
	// atlas migrate diff \
	//     --dev-url "${connection from ent.Client}" \
	//     --to "ent://ent/schema" \
	//     --dir "file://${config.migrationDir}"
	opts := append(defaultMigrateOptions(), schema.WithDir(dir))
	err = client.Schema.Diff(ctx, opts...)
	if err != nil {
		return err
	}
	return nil
}

// `atlas migrate apply` へのエイリアスです.
//
// `MigrationsDir` のディレクトリを参照してdiffの計算が行われます.
func MigrateApply(ctx context.Context, client *ent.Client, options ...MigrateOption) error {
	config := defaultMigrateConfig()
	config.applyOptions(options...)

	dir, err := atlas.NewLocalDir(config.migrationsDir)
	if err != nil {
		return err
	}
	// atlas migrate apply \
	//     --url "${connection from ent.Client}" \
	//     --dir "file://${MigrationDir}"
	opts := append(defaultMigrateOptions(), schema.WithDir(dir))
	err = client.Schema.Create(ctx, opts...)
	if err != nil {
		return err
	}
	return nil
}

// `atlas migrate diff` へのエイリアスです.
//
// `MigrationsDir` のディレクトリを参照してdiffの計算が行われます.
func (r *EntRepository) MigrateDiff(ctx context.Context, options ...MigrateOption) error {
	return MigrateDiff(ctx, r.client, options...)
}

// `atlas migrate apply` へのエイリアスです.
//
// `MigrationsDir` のディレクトリを参照してdiffの計算が行われます.
func (r *EntRepository) MigrateApply(ctx context.Context, options ...MigrateOption) error {
	return MigrateApply(ctx, r.client, options...)
}
