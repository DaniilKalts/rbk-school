package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	redisclient "github.com/redis/go-redis/v9"

	domainweather "github.com/DaniilKalts/rbk-school/4-week/internal/domain/weather"
)

const weatherKeyPrefix = "weather:"

type WeatherCache struct {
	client *redisclient.Client
	ttl    time.Duration
}

func NewWeatherCache(client *redisclient.Client, ttl time.Duration) *WeatherCache {
	return &WeatherCache{client: client, ttl: ttl}
}

func (c *WeatherCache) Get(ctx context.Context, city string) (domainweather.Weather, bool, error) {
	value, err := c.client.Get(ctx, weatherKey(city)).Result()
	if errors.Is(err, redisclient.Nil) {
		return domainweather.Weather{}, false, nil
	}
	if err != nil {
		return domainweather.Weather{}, false, fmt.Errorf("получение кеша погоды: %w", err)
	}

	var weather domainweather.Weather
	if err := json.Unmarshal([]byte(value), &weather); err != nil {
		return domainweather.Weather{}, false, fmt.Errorf("декодирование кеша погоды: %w", err)
	}

	return weather, true, nil
}

func (c *WeatherCache) Set(ctx context.Context, city string, weather domainweather.Weather) error {
	weather.RequestedAt = time.Time{}

	data, err := json.Marshal(weather)
	if err != nil {
		return fmt.Errorf("кодирование кеша погоды: %w", err)
	}

	if err := c.client.Set(ctx, weatherKey(city), data, c.ttl).Err(); err != nil {
		return fmt.Errorf("запись кеша погоды: %w", err)
	}

	return nil
}

func weatherKey(city string) string {
	return weatherKeyPrefix + strings.ToLower(domainweather.NormalizeCityName(city))
}
