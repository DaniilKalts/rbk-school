package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/rbk-school/7-week/internal/cache/blacklist"
	"github.com/DaniilKalts/rbk-school/7-week/internal/cache/weather"

	domainweather "github.com/DaniilKalts/rbk-school/7-week/internal/domain/weather"
)

type TokenBlacklist interface {
	Revoke(ctx context.Context, token string, expiresAt time.Time) error
	IsRevoked(ctx context.Context, token string) (bool, error)
}

type WeatherCache interface {
	Get(ctx context.Context, city string) (domainweather.Weather, bool, error)
	Set(ctx context.Context, city string, weather domainweather.Weather) error
}

type Caches struct {
	TokenBlacklist TokenBlacklist
	Weather        WeatherCache
}

func NewCaches(redisClient *redis.Client, weatherTTL time.Duration) *Caches {
	return &Caches{
		TokenBlacklist: blacklist.NewCache(redisClient),
		Weather:        weather.NewCache(redisClient, weatherTTL),
	}
}
