package jwt

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Secret         string        `env:"SECRET"`
	AccessTokenTTL time.Duration `env:"ACCESS_TOKEN_TTL" envDefault:"15m"`
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.Secret) == "" {
		return fmt.Errorf("секрет обязателен")
	}

	if c.AccessTokenTTL <= 0 {
		return fmt.Errorf("ttl access token должен быть положительным")
	}

	return nil
}
