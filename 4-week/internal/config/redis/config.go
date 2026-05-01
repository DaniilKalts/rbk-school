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
		return fmt.Errorf("адрес обязателен")
	}

	if c.DB < 0 {
		return fmt.Errorf("номер базы redis не может быть отрицательным")
	}

	if c.DialTimeout <= 0 {
		return fmt.Errorf("dial timeout должен быть положительным")
	}

	if c.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout должен быть положительным")
	}

	if c.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout должен быть положительным")
	}

	if c.WeatherCacheTTL <= 0 {
		return fmt.Errorf("ttl кеша погоды должен быть положительным")
	}

	return nil
}
