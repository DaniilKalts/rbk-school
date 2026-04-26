package countrystatecity

import (
	"fmt"
	"strings"
)

type Config struct {
	APIKey string `env:"API_KEY"`
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.APIKey) == "" {
		return fmt.Errorf("api key is required")
	}

	return nil
}
