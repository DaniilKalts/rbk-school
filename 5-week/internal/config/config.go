package config

import (
	"fmt"

	"github.com/DaniilKalts/rbk-school/5-week/internal/config/jwt"
	"github.com/DaniilKalts/rbk-school/5-week/internal/config/postgres"
	"github.com/DaniilKalts/rbk-school/5-week/internal/config/redis"
	"github.com/DaniilKalts/rbk-school/5-week/internal/config/server"
)

type Config struct {
	Server   server.Config   `envPrefix:"SERVER_"`
	Postgres postgres.Config `envPrefix:"POSTGRES_"`
	Redis    redis.Config    `envPrefix:"REDIS_"`
	JWT      jwt.Config      `envPrefix:"JWT_"`
}

func (c Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("конфигурация сервера: %w", err)
	}

	if err := c.Postgres.Validate(); err != nil {
		return fmt.Errorf("конфигурация postgres: %w", err)
	}

	if err := c.Redis.Validate(); err != nil {
		return fmt.Errorf("конфигурация redis: %w", err)
	}

	if err := c.JWT.Validate(); err != nil {
		return fmt.Errorf("конфигурация jwt: %w", err)
	}

	return nil
}
