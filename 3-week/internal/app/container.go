package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	redisclient "github.com/redis/go-redis/v9"

	cacheredis "github.com/DaniilKalts/rbk-school/3-week/internal/adapters/cache/redis"
	"github.com/DaniilKalts/rbk-school/3-week/internal/adapters/database/postgres"
	docshttp "github.com/DaniilKalts/rbk-school/3-week/internal/adapters/transport/http/docs"
	cityhttp "github.com/DaniilKalts/rbk-school/3-week/internal/adapters/transport/http/v1/city"
	userhttp "github.com/DaniilKalts/rbk-school/3-week/internal/adapters/transport/http/v1/user"
	weatherhttp "github.com/DaniilKalts/rbk-school/3-week/internal/adapters/transport/http/v1/weather"
	"github.com/DaniilKalts/rbk-school/3-week/internal/clients/geocoding"
	"github.com/DaniilKalts/rbk-school/3-week/internal/clients/openmeteo"
	"github.com/DaniilKalts/rbk-school/3-week/internal/config"
	cityrepo "github.com/DaniilKalts/rbk-school/3-week/internal/repository/city"
	userrepo "github.com/DaniilKalts/rbk-school/3-week/internal/repository/user"
	weatherrepo "github.com/DaniilKalts/rbk-school/3-week/internal/repository/weather"
	cityservice "github.com/DaniilKalts/rbk-school/3-week/internal/service/city"
	userservice "github.com/DaniilKalts/rbk-school/3-week/internal/service/user"
	weatherservice "github.com/DaniilKalts/rbk-school/3-week/internal/service/weather"
)

type Container struct {
	Config *config.Config

	DB         *pgxpool.Pool
	Redis      *redisclient.Client
	HTTPClient *http.Client
	Router     *http.ServeMux
}

func NewContainer(cfg *config.Config) (*Container, error) {
	if cfg == nil {
		return nil, fmt.Errorf("container: config is nil")
	}

	db, err := postgres.New(context.Background(), &cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("container: create postgres client: %w", err)
	}

	redisClient, err := cacheredis.New(context.Background(), &cfg.Redis)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("container: create redis client: %w", err)
	}

	httpClient := &http.Client{Timeout: cfg.Server.HTTPTimeout}

	repos := initRepositories(db)
	clients := initClients(httpClient, redisClient, cfg)
	services := initServices(repos, clients)

	router := newRouter(services.user, services.city, services.weather)

	return &Container{
		Config:     cfg,
		DB:         db,
		Redis:      redisClient,
		HTTPClient: httpClient,
		Router:     router,
	}, nil
}

func (c *Container) Close() {
	if c == nil {
		return
	}

	if c.DB != nil {
		c.DB.Close()
	}
	if c.Redis != nil {
		_ = c.Redis.Close()
	}
}

type repositories struct {
	user    *userrepo.Repository
	city    *cityrepo.Repository
	weather *weatherrepo.Repository
}

func initRepositories(db *pgxpool.Pool) *repositories {
	return &repositories{
		user:    userrepo.New(db),
		city:    cityrepo.New(db),
		weather: weatherrepo.New(db),
	}
}

type clients struct {
	geocoding    *geocoding.Client
	openMeteo    *openmeteo.Client
	weatherCache *weatherrepo.WeatherCache
}

func initClients(httpClient *http.Client, redisClient *redisclient.Client, cfg *config.Config) *clients {
	return &clients{
		geocoding:    geocoding.NewClient(httpClient),
		openMeteo:    openmeteo.NewClient(httpClient),
		weatherCache: weatherrepo.NewWeatherCache(redisClient, cfg.Redis.WeatherCacheTTL),
	}
}

type services struct {
	user    *userservice.Service
	city    *cityservice.Service
	weather *weatherservice.Service
}

func initServices(repos *repositories, clients *clients) *services {
	return &services{
		user:    userservice.New(repos.user),
		city:    cityservice.New(repos.city, repos.user),
		weather: weatherservice.New(repos.user, repos.city, repos.weather, clients.geocoding, clients.openMeteo, clients.weatherCache),
	}
}

func newRouter(userService userhttp.Service, cityService cityhttp.Service, weatherService weatherhttp.Service) *http.ServeMux {
	mux := http.NewServeMux()
	docshttp.RegisterRoutes(mux)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	userhttp.RegisterRoutes(mux, userService)
	cityhttp.RegisterRoutes(mux, cityService)
	weatherhttp.RegisterRoutes(mux, weatherService)

	return mux
}
