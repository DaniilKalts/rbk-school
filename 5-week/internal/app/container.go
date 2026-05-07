package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	redisclient "github.com/DaniilKalts/rbk-school/5-week/internal/adapters/cache/redis"
	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/database/postgres"
	transporthttp "github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http"
	v1 "github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/v1"
	"github.com/DaniilKalts/rbk-school/5-week/internal/client"
	"github.com/DaniilKalts/rbk-school/5-week/internal/config"
	"github.com/DaniilKalts/rbk-school/5-week/internal/repository"
	"github.com/DaniilKalts/rbk-school/5-week/internal/service"
	"github.com/DaniilKalts/rbk-school/5-week/internal/utils"
)

type Container struct {
	// Конфигурация приложения.
	config *config.Config

	// Инфраструктурные зависимости.
	db         *pgxpool.Pool
	redis      *redis.Client
	httpClient *http.Client

	// Внешние клиенты.
	clients *client.Clients

	// Репозитории.
	repositories *repository.Repositories

	// Кеши.
	caches *repository.Caches

	// Сервисы и безопасность.
	tokenManager *utils.JWTManager
	services     *service.Services

	// Транспортный слой.
	router http.Handler
}

func NewContainer(cfg *config.Config) (*Container, error) {
	if cfg == nil {
		return nil, fmt.Errorf("container: конфигурация не задана")
	}

	return &Container{config: cfg}, nil
}

func (c *Container) Config() *config.Config {
	return c.config
}

func (c *Container) DB() *pgxpool.Pool {
	if c.db == nil {
		db, err := postgres.NewClient(context.Background(), &c.config.Postgres)
		if err != nil {
			log.Fatalf("не удалось создать клиент postgres: %v", err)
		}

		c.db = db
	}

	return c.db
}

func (c *Container) Redis() *redis.Client {
	if c.redis == nil {
		redisClient, err := redisclient.NewClient(context.Background(), &c.config.Redis)
		if err != nil {
			log.Fatalf("не удалось создать клиент redis: %v", err)
		}

		c.redis = redisClient
	}

	return c.redis
}

func (c *Container) HTTPClient() *http.Client {
	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: c.config.Server.HTTPTimeout}
	}

	return c.httpClient
}

func (c *Container) Clients() *client.Clients {
	if c.clients == nil {
		c.clients = client.NewClients(c.HTTPClient(), c.Redis(), c.config.Redis.WeatherCacheTTL)
	}

	return c.clients
}

func (c *Container) TokenBlacklist() repository.TokenBlacklist {
	return c.Caches().TokenBlacklist
}

func (c *Container) UserRepository() repository.UserRepository {
	return c.Repositories().User
}

func (c *Container) CityRepository() repository.CityRepository {
	return c.Repositories().City
}

func (c *Container) WeatherRepository() repository.WeatherRepository {
	return c.Repositories().Weather
}

func (c *Container) Repositories() *repository.Repositories {
	if c.repositories == nil {
		c.repositories = repository.NewRepositories(c.DB())
	}

	return c.repositories
}

func (c *Container) Caches() *repository.Caches {
	if c.caches == nil {
		c.caches = repository.NewCaches(c.Redis())
	}

	return c.caches
}

func (c *Container) TokenManager() *utils.JWTManager {
	if c.tokenManager == nil {
		c.tokenManager = utils.NewJWTManager([]byte(c.config.JWT.Secret), c.config.JWT.AccessTokenTTL, c.TokenBlacklist())
	}

	return c.tokenManager
}

func (c *Container) AuthService() service.AuthService {
	return c.Services().Auth
}

func (c *Container) UserService() service.UserService {
	return c.Services().User
}

func (c *Container) CityService() service.CityService {
	return c.Services().City
}

func (c *Container) WeatherService() service.WeatherService {
	return c.Services().Weather
}

func (c *Container) Services() *service.Services {
	if c.services == nil {
		c.services = service.NewServicesFromDependencies(c.Repositories(), c.Clients(), c.TokenManager())
	}

	return c.services
}

func (c *Container) Router() http.Handler {
	if c.router == nil {
		c.router = transporthttp.NewRouter(v1.Dependencies{
			AuthService:    c.AuthService(),
			CityService:    c.CityService(),
			WeatherService: c.WeatherService(),
			UserService:    c.UserService(),
			JWTManager:     c.TokenManager(),
		})
	}

	return c.router
}

func (c *Container) Close() {
	if c == nil {
		return
	}

	if c.db != nil {
		c.db.Close()
	}
	if c.redis != nil {
		_ = c.redis.Close()
	}
}
