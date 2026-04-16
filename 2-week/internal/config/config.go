package config

type Config struct {
	CountryStateCity CountryStateCityConfig `envPrefix:"COUNTRY_STATE_CITY_"`
}

type CountryStateCityConfig struct {
	APIKey string `env:"API_KEY"`
}
