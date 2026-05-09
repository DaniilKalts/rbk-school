package blacklist

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const tokenBlacklistPrefix = "token_blacklist:"

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Revoke(ctx context.Context, token string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}

	if err := c.client.Set(ctx, tokenBlacklistPrefix+token, "1", ttl).Err(); err != nil {
		return fmt.Errorf("добавление отозванного токена: %w", err)
	}

	return nil
}

func (c *Cache) IsRevoked(ctx context.Context, token string) (bool, error) {
	count, err := c.client.Exists(ctx, tokenBlacklistPrefix+token).Result()
	if err != nil {
		return false, fmt.Errorf("проверка отозванного токена: %w", err)
	}

	return count > 0, nil
}
