package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func Load(path string) (*Config, error) {
	if path != "" {
		if err := godotenv.Load(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("загрузка env-файла %q: %w", path, err)
		}
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("разбор конфигурации из окружения: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("валидация конфигурации: %w", err)
	}

	return &cfg, nil
}
