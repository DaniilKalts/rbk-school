package redis

import (
	"context"
	"fmt"

	redisclient "github.com/redis/go-redis/v9"

	configredis "github.com/DaniilKalts/rbk-school/3-week/internal/config/redis"
)

func New(ctx context.Context, cfg *configredis.Config) (*redisclient.Client, error) {
	client := redisclient.NewClient(&redisclient.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("ping redis at %s: %w", cfg.Addr, err)
	}

	return client, nil
}
