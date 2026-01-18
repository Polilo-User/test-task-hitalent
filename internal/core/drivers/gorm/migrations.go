package psql

import (
	"context"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pressly/goose/v3"
)

func (d *Driver) Migrate(ctx context.Context, migrationsDir string) error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}

	goose.SetDialect("postgres")

	return goose.UpContext(ctx, sqlDB, migrationsDir)
}

func (d *Driver) RollbackAll(ctx context.Context, migrationsDir string) error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}

	goose.SetDialect("postgres")

	return goose.DownToContext(ctx, sqlDB, migrationsDir, 0)
}
