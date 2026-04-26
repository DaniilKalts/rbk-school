package redis

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Addr            string        `env:"ADDR" envDefault:"localhost:6379"`
	Password        string        `env:"PASSWORD"`
	DB              int           `env:"DB" envDefault:"0"`
	DialTimeout     time.Duration `env:"DIAL_TIMEOUT" envDefault:"5s"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT" envDefault:"3s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" envDefault:"3s"`
	WeatherCacheTTL time.Duration `env:"WEATHER_CACHE_TTL" envDefault:"10m"`
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.Addr) == "" {
		return fmt.Errorf("addr is required")
	}

	if c.DB < 0 {
		return fmt.Errorf("db must not be negative")
	}

	if c.DialTimeout <= 0 {
		return fmt.Errorf("dial timeout must be positive")
	}

	if c.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout must be positive")
	}

	if c.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout must be positive")
	}

	if c.WeatherCacheTTL <= 0 {
		return fmt.Errorf("weather cache ttl must be positive")
	}

	return nil
}
