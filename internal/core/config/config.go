package config

import (
	"context"
	"time"

	"github.com/Polilo-User/test-task-hitalent/internal/core/errors"
	"github.com/Polilo-User/test-task-hitalent/internal/core/logging"
	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const (
	ErrRead      = errors.Error("failed to read file")
	ErrUnmarshal = errors.Error("failed to unmarshal file")
)

type Config struct {
	AppName  string        `env:"APP_NAME" validate:"required"`
	Env      string        `env:"ENV" validate:"required,oneof=local docker prod"`
	HTTPPort string        `env:"HTTP_PORT"`
	Timeout  time.Duration `env:"TIMEOUT" envDefault:"30s"`
}

func Load(ctx context.Context) (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	logging.From(ctx).Info("Loaded configuration from environment",
		zap.String("env", cfg.Env),
		zap.String("http_port", cfg.HTTPPort),
	)

	return cfg, nil
}
