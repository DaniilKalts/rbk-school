package app

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/DaniilKalts/rbk-school/2-week/internal/client/countrystatecity"
	"github.com/DaniilKalts/rbk-school/2-week/internal/client/geocoding"
	"github.com/DaniilKalts/rbk-school/2-week/internal/client/openmeteo"
	"github.com/DaniilKalts/rbk-school/2-week/internal/config"
	weatherService "github.com/DaniilKalts/rbk-school/2-week/internal/service/weather"
	weatherHandler "github.com/DaniilKalts/rbk-school/2-week/internal/transport/http/v1/weather"
)

type Container struct {
	Config *config.Config

	HTTPClient *http.Client
	Router     *chi.Mux
}

func NewContainer(cfg *config.Config) (*Container, error) {
	if cfg == nil {
		return nil, fmt.Errorf("container: config is nil")
	}

	httpClient := &http.Client{Timeout: cfg.Server.HTTPTimeout}

	countryStateCityClient := countrystatecity.NewClient(httpClient, cfg.CountryStateCity.APIKey)
	geocodingClient := geocoding.NewClient(httpClient)
	openMeteoClient := openmeteo.NewClient(httpClient)

	weatherSvc := weatherService.NewService(countryStateCityClient, geocodingClient, openMeteoClient)
	weatherHTTPHandler := weatherHandler.NewWeatherHandler(weatherSvc)

	router := newRouter(cfg, weatherHTTPHandler)

	return &Container{
		Config:     cfg,
		HTTPClient: httpClient,
		Router:     router,
	}, nil
}

func newRouter(cfg *config.Config, weatherHandler *weatherHandler.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(cfg.Server.HTTPTimeout))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Route("/api/v1/weather", weatherHandler.RegisterRoutes)

	return r
}
