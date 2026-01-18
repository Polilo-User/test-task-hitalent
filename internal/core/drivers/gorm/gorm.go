package psql

import (
	"context"
	"time"

	"github.com/Polilo-User/test-task-hitalent/internal/core/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	ErrConnect = errors.Error("failed to connect to postgres db")
	ErrClose   = errors.Error("failed to close postgres db connection")
)

type Driver struct {
	dsn string
	db  *gorm.DB
}

func New(dsn string) *Driver {
	return &Driver{
		dsn: dsn,
	}
}

func (d *Driver) Connect(ctx context.Context) error {
	gormDB, err := gorm.Open(postgres.Open(d.dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return ErrConnect.Wrap(err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return ErrConnect.Wrap(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.PingContext(ctx); err != nil {
		return ErrConnect.Wrap(err)
	}

	d.db = gormDB
	return nil
}

func (d *Driver) Close(ctx context.Context) error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return ErrClose.Wrap(err)
	}

	if err := sqlDB.Close(); err != nil {
		return ErrClose.Wrap(err)
	}

	return nil
}

func (d *Driver) GetDB() *gorm.DB {
	return d.db
}

func (d *Driver) PingContext(ctx context.Context) error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
