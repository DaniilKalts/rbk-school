package postgres

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	Host            string        `env:"HOST" envDefault:"localhost"`
	Port            int           `env:"PORT" envDefault:"5432"`
	User            string        `env:"USER" envDefault:"postgres"`
	Password        string        `env:"PASSWORD" envDefault:"postgres"`
	Database        string        `env:"DATABASE" envDefault:"weather_api"`
	SSLMode         string        `env:"SSL_MODE" envDefault:"disable"`
	MaxConns        int32         `env:"MAX_CONNS" envDefault:"10"`
	MinConns        int32         `env:"MIN_CONNS" envDefault:"1"`
	MaxConnLifetime time.Duration `env:"MAX_CONN_LIFETIME" envDefault:"1h"`
	MaxConnIdleTime time.Duration `env:"MAX_CONN_IDLE_TIME" envDefault:"30m"`
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.Host) == "" {
		return fmt.Errorf("host обязателен")
	}

	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port должен быть в диапазоне от 1 до 65535")
	}

	if strings.TrimSpace(c.User) == "" {
		return fmt.Errorf("user обязателен")
	}

	if strings.TrimSpace(c.Database) == "" {
		return fmt.Errorf("database обязателен")
	}

	if !isValidSSLMode(c.SSLMode) {
		return fmt.Errorf("ssl mode должен быть одним из: disable, allow, prefer, require, verify-ca, verify-full")
	}

	if c.MaxConns <= 0 {
		return fmt.Errorf("max conns должен быть положительным")
	}

	if c.MinConns < 0 {
		return fmt.Errorf("min conns не может быть отрицательным")
	}

	if c.MinConns > c.MaxConns {
		return fmt.Errorf("min conns не может быть больше max conns")
	}

	if c.MaxConnLifetime <= 0 {
		return fmt.Errorf("max conn lifetime должен быть положительным")
	}

	if c.MaxConnIdleTime <= 0 {
		return fmt.Errorf("max conn idle time должен быть положительным")
	}

	return nil
}

func (c Config) DSN() string {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:   c.Database,
	}

	query := dsn.Query()
	query.Set("sslmode", c.SSLMode)
	dsn.RawQuery = query.Encode()

	return dsn.String()
}

func isValidSSLMode(value string) bool {
	switch value {
	case "disable", "allow", "prefer", "require", "verify-ca", "verify-full":
		return true
	default:
		return false
	}
}
