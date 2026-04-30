package server

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Addr        string        `env:"ADDR" envDefault:":8080"`
	HTTPTimeout time.Duration `env:"HTTP_TIMEOUT" envDefault:"15s"`
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.Addr) == "" {
		return fmt.Errorf("addr is required")
	}

	if c.HTTPTimeout <= 0 {
		return fmt.Errorf("http timeout must be positive")
	}

	return nil
}
