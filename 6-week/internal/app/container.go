package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	redisclient "github.com/DaniilKalts/rbk-school/6-week/internal/adapter/cache/redis"
	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/client"
	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/database/postgres"
	"github.com/DaniilKalts/rbk-school/6-week/internal/cache"
	"github.com/DaniilKalts/rbk-school/6-week/internal/config"
	"github.com/DaniilKalts/rbk-school/6-week/internal/repository"
	"github.com/DaniilKalts/rbk-school/6-week/internal/service"
	"github.com/DaniilKalts/rbk-school/6-week/pkg/jwt"
)

type Container struct {
	Config *config.Config
	Logger *zap.Logger

	DB    *pgxpool.Pool
	Redis *redis.Client

	Repositories *repository.Repositories
	Caches       *cache.Caches

	Clients      *client.Clients
	TokenManager *jwt.Manager
	Services     *service.Services
}

func NewContainer(cfg *config.Config, logger *zap.Logger) (_ *Container, err error) {
	ctx := context.Background()

	db, err := postgres.NewClient(ctx, &cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	defer func() {
		if err != nil {
			db.Close()
		}
	}()

	redisClient, err := redisclient.NewClient(ctx, &cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}
	defer func() {
		if err != nil {
			_ = redisClient.Close()
		}
	}()

	repositories := repository.NewRepositories(db)
	caches := cache.NewCaches(redisClient, cfg.Redis.WeatherCacheTTL)

	clients := client.NewClients(&http.Client{Timeout: cfg.Server.HTTPTimeout})
	tokenManager := jwt.NewManager([]byte(cfg.JWT.Secret), cfg.JWT.AccessTokenTTL, caches.TokenBlacklist)
	services := service.NewServices(repositories, caches, clients, tokenManager)

	return &Container{
		Config:       cfg,
		Logger:       logger,
		DB:           db,
		Redis:        redisClient,
		Repositories: repositories,
		Caches:       caches,
		Clients:      clients,
		TokenManager: tokenManager,
		Services:     services,
	}, nil
}

func (c *Container) Close() {
	c.DB.Close()
	if err := c.Redis.Close(); err != nil {
		c.Logger.Warn("не удалось закрыть клиент redis", zap.Error(err))
	}
}
