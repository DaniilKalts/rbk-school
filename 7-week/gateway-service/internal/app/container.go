package app

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/adapter/client/geocoding"
	"github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/adapter/client/openmeteo"
	v1 "github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/adapter/transport/http/v1"
	"github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/config"
	weathersvc "github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/service/weather"
)

type Container struct {
	Config *config.Config
	Logger *zap.Logger

	Dependencies v1.Dependencies
}

func NewContainer(cfg *config.Config, logger *zap.Logger) (*Container, error) {
	httpClient := &http.Client{Timeout: cfg.External.Timeout}

	geocodingClient := geocoding.NewClient(httpClient, cfg.External.GeocodingURL)
	openMeteoClient := openmeteo.NewClient(httpClient, cfg.External.OpenMeteoURL)

	weatherService := weathersvc.NewService(geocodingClient, openMeteoClient)

	return &Container{
		Config: cfg,
		Logger: logger,
		Dependencies: v1.Dependencies{
			WeatherService: weatherService,
		},
	}, nil
}

func (c *Container) Close() {}
