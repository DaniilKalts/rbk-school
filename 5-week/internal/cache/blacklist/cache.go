package blacklist

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const tokenBlacklistPrefix = "token_blacklist:"

type Blacklist struct {
	client *redis.Client
}

func NewBlacklist(client *redis.Client) *Blacklist {
	return &Blacklist{client: client}
}

func (b *Blacklist) Revoke(ctx context.Context, token string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}

	if err := b.client.Set(ctx, tokenBlacklistPrefix+token, "1", ttl).Err(); err != nil {
		return fmt.Errorf("добавление отозванного токена: %w", err)
	}

	return nil
}

func (b *Blacklist) IsRevoked(ctx context.Context, token string) (bool, error) {
	count, err := b.client.Exists(ctx, tokenBlacklistPrefix+token).Result()
	if err != nil {
		return false, fmt.Errorf("проверка отозванного токена: %w", err)
	}

	return count > 0, nil
}
