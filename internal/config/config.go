package config

import (
	"context"

	"github.com/Polilo-User/test-task-hitalent/internal/core/errors"

	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
)

const (
	ErrInvalidEnvironment = errors.Error("ENV is not set")
	ErrValidation         = errors.Error("invalid configuration")
	ErrEnvVars            = errors.Error("failed parsing env vars")
)

type Config struct {
	HTTP_PORT string `env:"HTTP_PORT"`
	PSQL      string `env:"POSTGRES_DSN"`
}

func Load(ctx context.Context) (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, ErrEnvVars.Wrap(err)
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, ErrValidation.Wrap(err)
	}

	return cfg, nil
}
