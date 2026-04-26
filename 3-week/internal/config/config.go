package config

import (
	"fmt"

	"github.com/DaniilKalts/rbk-school/3-week/internal/config/postgres"
	"github.com/DaniilKalts/rbk-school/3-week/internal/config/redis"
	"github.com/DaniilKalts/rbk-school/3-week/internal/config/server"
)

type Config struct {
	Server           server.Config           `envPrefix:"SERVER_"`
	Postgres         postgres.Config         `envPrefix:"POSTGRES_"`
	Redis            redis.Config            `envPrefix:"REDIS_"`
}

func (c Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	if err := c.Postgres.Validate(); err != nil {
		return fmt.Errorf("postgres config: %w", err)
	}

	if err := c.Redis.Validate(); err != nil {
		return fmt.Errorf("redis config: %w", err)
	}

	return nil
}
