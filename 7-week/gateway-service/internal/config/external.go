package config

import (
	"fmt"
	"strings"
	"time"
)

type External struct {
	Timeout      time.Duration `env:"TIMEOUT" envDefault:"15s"`
	GeocodingURL string        `env:"GEOCODING_URL" envDefault:"https://geocoding-api.open-meteo.com/v1/search"`
	OpenMeteoURL string        `env:"OPENMETEO_URL" envDefault:"https://api.open-meteo.com/v1/forecast"`
}

func (c External) Validate() error {
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout должен быть положительным")
	}

	if strings.TrimSpace(c.GeocodingURL) == "" {
		return fmt.Errorf("geocoding url обязателен")
	}

	if strings.TrimSpace(c.OpenMeteoURL) == "" {
		return fmt.Errorf("openmeteo url обязателен")
	}

	return nil
}
