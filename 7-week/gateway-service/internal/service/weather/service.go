package weather

import (
	"context"
	"fmt"

	geocodingdto "github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/adapter/client/geocoding/dto"
	openmeteodto "github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/adapter/client/openmeteo/dto"
)

type GeocodingClient interface {
	GetCoordsByCity(ctx context.Context, city string) (geocodingdto.CoordsResponse, error)
}

type OpenMeteoClient interface {
	GetWeatherByCoords(ctx context.Context, latitude, longitude float64) (openmeteodto.WeatherResponse, error)
}

type Snapshot struct {
	City                string
	Latitude            float64
	Longitude           float64
	Temperature         float64
	ApparentTemperature float64
	WeatherCode         int
}

type Service struct {
	geocodingClient GeocodingClient
	openMeteoClient OpenMeteoClient
}

func NewService(geocodingClient GeocodingClient, openMeteoClient OpenMeteoClient) *Service {
	return &Service{
		geocodingClient: geocodingClient,
		openMeteoClient: openMeteoClient,
	}
}

func (s *Service) GetByCity(ctx context.Context, city string) (Snapshot, error) {
	coords, err := s.geocodingClient.GetCoordsByCity(ctx, city)
	if err != nil {
		return Snapshot{}, fmt.Errorf("получение координат для города %q: %w", city, err)
	}

	weather, err := s.openMeteoClient.GetWeatherByCoords(ctx, coords.Latitude, coords.Longitude)
	if err != nil {
		return Snapshot{}, fmt.Errorf("получение погоды для города %q: %w", city, err)
	}

	return Snapshot{
		City:                city,
		Latitude:            coords.Latitude,
		Longitude:           coords.Longitude,
		Temperature:         weather.Current.Temperature2M,
		ApparentTemperature: weather.Current.ApparentTemperature,
		WeatherCode:         weather.Current.WeatherCode,
	}, nil
}
