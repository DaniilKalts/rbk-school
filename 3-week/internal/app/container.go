package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	redisclient "github.com/redis/go-redis/v9"

	cacheredis "github.com/DaniilKalts/rbk-school/3-week/internal/adapters/cache/redis"
	"github.com/DaniilKalts/rbk-school/3-week/internal/adapters/client/geocoding"
	"github.com/DaniilKalts/rbk-school/3-week/internal/adapters/client/openmeteo"
	"github.com/DaniilKalts/rbk-school/3-week/internal/adapters/database/postgres"
	docshttp "github.com/DaniilKalts/rbk-school/3-week/internal/adapters/transport/http/docs"
	cityhttp "github.com/DaniilKalts/rbk-school/3-week/internal/adapters/transport/http/v1/city"
	userhttp "github.com/DaniilKalts/rbk-school/3-week/internal/adapters/transport/http/v1/user"
	weatherhttp "github.com/DaniilKalts/rbk-school/3-week/internal/adapters/transport/http/v1/weather"
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

	userRepository := userrepo.New(db)
	userService := userservice.New(userRepository)

	cityRepository := cityrepo.New(db)
	cityService := cityservice.New(cityRepository, userRepository)

	weatherRepository := weatherrepo.New(db)
	geocodingClient := geocoding.NewClient(httpClient)
	openMeteoClient := openmeteo.NewClient(httpClient)
	weatherCache := cacheredis.NewWeatherCache(redisClient, cfg.Redis.WeatherCacheTTL)
	weatherService := weatherservice.New(userRepository, cityRepository, weatherRepository, geocodingClient, openMeteoClient, weatherCache)

	router := newRouter(userService, cityService, weatherService)

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
