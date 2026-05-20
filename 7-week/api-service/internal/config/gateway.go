package config

import (
	"fmt"
	"strings"
	"time"
)

type Gateway struct {
	BaseURL string        `env:"BASE_URL" envDefault:"http://gateway-service:8081"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"15s"`
}

func (c Gateway) Validate() error {
	if strings.TrimSpace(c.BaseURL) == "" {
		return fmt.Errorf("base url обязателен")
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("timeout должен быть положительным")
	}

	return nil
}
