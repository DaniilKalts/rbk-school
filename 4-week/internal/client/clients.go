package client

import (
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/rbk-school/4-week/internal/client/geocoding"
	"github.com/DaniilKalts/rbk-school/4-week/internal/client/openmeteo"
	"github.com/DaniilKalts/rbk-school/4-week/internal/repository/weather"
)

type Clients struct {
	Geocoding    *geocoding.Client
	OpenMeteo    *openmeteo.Client
	WeatherCache *weather.WeatherCache
}

func NewClients(httpClient *http.Client, redisClient *redis.Client, weatherCacheTTL time.Duration) *Clients {
	return &Clients{
		Geocoding:    geocoding.NewClient(httpClient),
		OpenMeteo:    openmeteo.NewClient(httpClient),
		WeatherCache: weather.NewWeatherCache(redisClient, weatherCacheTTL),
	}
}
