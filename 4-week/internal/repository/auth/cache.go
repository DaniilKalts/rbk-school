package auth

import (
	"context"
	"fmt"
	"time"

	redisclient "github.com/redis/go-redis/v9"
)

const tokenBlacklistPrefix = "token_blacklist:"

type Cache struct {
	client *redisclient.Client
}

func NewRepository(client *redisclient.Client) *Cache {
	return &Cache{client: client}
}

func (r *Cache) Revoke(ctx context.Context, token string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}

	err := r.client.Set(ctx, tokenBlacklistPrefix+token, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("добавление отозванного токена: %w", err)
	}

	return nil
}

func (r *Cache) IsRevoked(ctx context.Context, token string) (bool, error) {
	count, err := r.client.Exists(ctx, tokenBlacklistPrefix+token).Result()
	if err != nil {
		return false, fmt.Errorf("проверка отозванного токена: %w", err)
	}

	return count > 0, nil
}
