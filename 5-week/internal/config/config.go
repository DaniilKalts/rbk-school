package config

import "fmt"

type Config struct {
	Server   Server   `envPrefix:"SERVER_"`
	Postgres Postgres `envPrefix:"POSTGRES_"`
	Redis    Redis    `envPrefix:"REDIS_"`
	JWT      JWT      `envPrefix:"JWT_"`
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
