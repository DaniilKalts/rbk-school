package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/rbk-school/5-week/internal/repository/auth"
)

type TokenBlacklist interface {
	Revoke(ctx context.Context, token string, expiresAt time.Time) error
	IsRevoked(ctx context.Context, token string) (bool, error)
}

type Caches struct {
	TokenBlacklist TokenBlacklist
}

func NewCaches(redisClient *redis.Client) *Caches {
	return &Caches{
		TokenBlacklist: auth.NewRepository(redisClient),
	}
}
