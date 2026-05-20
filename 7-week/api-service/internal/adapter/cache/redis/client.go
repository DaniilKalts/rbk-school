package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/rbk-school/7-week/internal/config"
)

func NewClient(ctx context.Context, cfg *config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("проверка подключения к redis %s: %w", cfg.Addr, err)
	}

	return client, nil
}
