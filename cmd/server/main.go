package main

import (
	"context"
	"fmt"

	"github.com/Polilo-User/test-task-hitalent/internal/chats"
	chatStore "github.com/Polilo-User/test-task-hitalent/internal/chats/store"
	"github.com/Polilo-User/test-task-hitalent/internal/config"
	"github.com/Polilo-User/test-task-hitalent/internal/core/app"
	psql "github.com/Polilo-User/test-task-hitalent/internal/core/drivers/gorm"
	"github.com/Polilo-User/test-task-hitalent/internal/core/listeners/http"
	"github.com/Polilo-User/test-task-hitalent/internal/core/logging"
	"github.com/Polilo-User/test-task-hitalent/internal/messages"
	messageStore "github.com/Polilo-User/test-task-hitalent/internal/messages/store"
	httptransport "github.com/Polilo-User/test-task-hitalent/internal/transport/http"

	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
)

func main() {
	app.Start(appStart)
}

func appStart(ctx context.Context, a *app.App) ([]app.Listener, error) {
	cfg, err := config.Load(ctx)
	if err != nil {
		return nil, err
	}

	db, err := initDatabase(ctx, cfg, a)
	if err != nil {
		return nil, err
	}

	if err := migrateDatabase(ctx, db); err != nil {
		return nil, err
	}
	a.OnShutdown(func() {
		rollbackMigrations(ctx, db)
	})

	cs := chatStore.New(db.GetDB())
	ms := messageStore.New(db.GetDB())
	c := chats.New(cs, ms)
	m := messages.New(ms, c)

	httpServer := httptransport.New(c, m, db.GetDB())

	h, err := http.New(httpServer, cfg.HTTP_PORT)
	if err != nil {
		return nil, err
	}

	return []app.Listener{
		h,
	}, nil
}

func initDatabase(ctx context.Context, cfg *config.Config, a *app.App) (*psql.Driver, error) {
	db := psql.New(cfg.PSQL)

	err := backoff.Retry(func() error {
		return db.Connect(ctx)
	}, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, err
	}

	a.OnShutdown(func() {
		logging.From(ctx).Info("shutting down db connection")
		if err := db.Close(ctx); err != nil {
			logging.From(ctx).Error("failed to close db connection", zap.Error(err))
		}
	})

	return db, nil
}

func migrateDatabase(ctx context.Context, db *psql.Driver) error {
	migrationsPath := "./migrations"

	if err := db.Migrate(ctx, migrationsPath); err != nil {
		return fmt.Errorf("migrateDatabase: %w", err)
	}
	return nil
}

func rollbackMigrations(ctx context.Context, db *psql.Driver) error {
	migrationsPath := "./migrations"

	if err := db.RollbackAll(ctx, migrationsPath); err != nil {
		return fmt.Errorf("rollbackMigrations: %w", err)
	}
	return nil
}
