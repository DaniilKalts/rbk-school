package config

import (
	"fmt"
	"strings"
	"time"
)

type Server struct {
	Addr            string        `env:"ADDR" envDefault:":8081"`
	HTTPTimeout     time.Duration `env:"HTTP_TIMEOUT" envDefault:"15s"`
	HandlerTimeout  time.Duration `env:"HANDLER_TIMEOUT" envDefault:"10s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"15s"`
}

func (c Server) Validate() error {
	if strings.TrimSpace(c.Addr) == "" {
		return fmt.Errorf("адрес обязателен")
	}

	if c.HTTPTimeout <= 0 {
		return fmt.Errorf("http timeout должен быть положительным")
	}

	if c.HandlerTimeout <= 0 {
		return fmt.Errorf("handler timeout должен быть положительным")
	}

	if c.ShutdownTimeout <= 0 {
		return fmt.Errorf("shutdown timeout должен быть положительным")
	}

	return nil
}
