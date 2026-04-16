package config

import "time"

type Config struct {
	CountryStateCity CountryStateCityConfig `envPrefix:"COUNTRY_STATE_CITY_"`
	Server           ServerConfig           `envPrefix:"SERVER_"`
}

type CountryStateCityConfig struct {
	APIKey string `env:"API_KEY"`
}

type ServerConfig struct {
	Addr        string        `env:"ADDR" envDefault:":8080"`
	HTTPTimeout time.Duration `env:"HTTP_TIMEOUT" envDefault:"15s"`
}
