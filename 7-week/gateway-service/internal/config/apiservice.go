package config

import (
	"fmt"
	"strings"
	"time"
)

type APIService struct {
	BaseURL string        `env:"BASE_URL" envDefault:"http://api-service:8080"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"5s"`
}

func (c APIService) Validate() error {
	if strings.TrimSpace(c.BaseURL) == "" {
		return fmt.Errorf("base url обязателен")
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("timeout должен быть положительным")
	}

	return nil
}
