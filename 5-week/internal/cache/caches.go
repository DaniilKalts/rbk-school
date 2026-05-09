package cache

import (
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/rbk-school/5-week/internal/cache/blacklist"
	"github.com/DaniilKalts/rbk-school/5-week/internal/cache/weather"
)

type Caches struct {
	TokenBlacklist *blacklist.Blacklist
	Weather        *weather.Cache
}

func NewCaches(redisClient *redis.Client, weatherTTL time.Duration) *Caches {
	return &Caches{
		TokenBlacklist: blacklist.NewBlacklist(redisClient),
		Weather:        weather.NewCache(redisClient, weatherTTL),
	}
}
