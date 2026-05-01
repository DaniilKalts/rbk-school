package auth

import (
	"context"
	"fmt"
	"time"

	redisclient "github.com/redis/go-redis/v9"
)

const tokenBlacklistPrefix = "token_blacklist:"

type TokenBlacklist struct {
	client *redisclient.Client
}

func NewRepository(client *redisclient.Client) *TokenBlacklist {
	return &TokenBlacklist{client: client}
}

func (r *TokenBlacklist) Revoke(ctx context.Context, tokenHash string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}

	err := r.client.Set(ctx, tokenBlacklistPrefix+tokenHash, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("set revoked token: %w", err)
	}

	return nil
}

func (r *TokenBlacklist) Contains(ctx context.Context, tokenHash string) (bool, error) {
	count, err := r.client.Exists(ctx, tokenBlacklistPrefix+tokenHash).Result()
	if err != nil {
		return false, fmt.Errorf("check revoked token: %w", err)
	}

	return count > 0, nil
}
