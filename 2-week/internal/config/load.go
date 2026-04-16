package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func Load(path string) (*Config, error) {
	if err := godotenv.Load(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("config: load env file: %w", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config: parse env: %w", err)
	}

	return cfg, nil
}
