package config

import "fmt"

type Config struct {
	Server   Server   `envPrefix:"SERVER_"`
	Logger   Logger   `envPrefix:"LOG_"`
	External External `envPrefix:"EXTERNAL_"`
}

func (c Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("конфигурация сервера: %w", err)
	}

	if err := c.Logger.Validate(); err != nil {
		return fmt.Errorf("конфигурация логгера: %w", err)
	}

	if err := c.External.Validate(); err != nil {
		return fmt.Errorf("конфигурация внешних API: %w", err)
	}

	return nil
}
