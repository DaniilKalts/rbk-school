package weather

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	countryStateCityDTO "github.com/DaniilKalts/rbk-school/2-week/internal/client/countrystatecity/dto"
	geocodingDTO "github.com/DaniilKalts/rbk-school/2-week/internal/client/geocoding/dto"
	openMeteoDTO "github.com/DaniilKalts/rbk-school/2-week/internal/client/openmeteo/dto"
	"github.com/DaniilKalts/rbk-school/2-week/internal/domain"
)

type CityListClient interface {
	GetStatesByCountry(ctx context.Context, countryCode string) ([]countryStateCityDTO.StateResponse, error)
}

type GeocodingClient interface {
	GetCoordsByState(ctx context.Context, countryCode string) (geocodingDTO.CoordsResponse, error)
}

type WeatherClient interface {
	GetWeatherByCoords(ctx context.Context, latitude, longitude float64) (openMeteoDTO.WeatherResponse, error)
}

type Service struct {
	cityListClient  CityListClient
	geocodingClient GeocodingClient
	weatherClient   WeatherClient
}

func NewService(
	cityListClient CityListClient,
	geocodingClient GeocodingClient,
	weatherClient WeatherClient,
) *Service {
	return &Service{
		cityListClient:  cityListClient,
		geocodingClient: geocodingClient,
		weatherClient:   weatherClient,
	}
}

func (s *Service) GetWeatherByCity(ctx context.Context, city string) (domain.Weather, error) {
	coords, err := s.geocodingClient.GetCoordsByState(ctx, city)
	if err != nil {
		return domain.Weather{}, fmt.Errorf("get coordinates for city %q: %w", city, err)
	}

	weatherResponse, err := s.weatherClient.GetWeatherByCoords(ctx, coords.Latitude, coords.Longitude)
	if err != nil {
		return domain.Weather{}, fmt.Errorf("get weather for city %q: %w", city, err)
	}

	conditions := domain.Conditions{
		Temperature: weatherResponse.Current.Temperature2M,
		FeelsLike:   weatherResponse.Current.ApparentTemperature,
	}

	weather, err := domain.NewWeather(city, coords.Latitude, coords.Longitude, conditions)
	if err != nil {
		return domain.Weather{}, fmt.Errorf("build weather for city %q: %w", city, err)
	}

	return weather, nil
}

func (s *Service) GetWeatherByCountry(ctx context.Context, countryCode string) ([]domain.Weather, error) {
	states, err := s.cityListClient.GetStatesByCountry(ctx, countryCode)
	if err != nil {
		return nil, fmt.Errorf("get states for country %q: %w", countryCode, err)
	}

	weathers := make([]domain.Weather, 0, len(states))

	for _, state := range states {
		latitude, err := strconv.ParseFloat(state.Latitude, 64)
		if err != nil {
			return nil, fmt.Errorf("parse latitude for city %q: %w", state.Name, err)
		}

		longitude, err := strconv.ParseFloat(state.Longitude, 64)
		if err != nil {
			return nil, fmt.Errorf("parse longitude for city %q: %w", state.Name, err)
		}

		weatherResponse, err := s.weatherClient.GetWeatherByCoords(ctx, latitude, longitude)
		if err != nil {
			continue
		}

		conditions := domain.Conditions{
			Temperature: weatherResponse.Current.Temperature2M,
			FeelsLike:   weatherResponse.Current.ApparentTemperature,
		}

		weather, err := domain.NewWeather(state.Name, latitude, longitude, conditions)
		if err != nil {
			return nil, fmt.Errorf("build weather for city %q: %w", state.Name, err)
		}

		weathers = append(weathers, weather)
	}

	return weathers, nil
}

func (s *Service) GetTopWarmestCities(ctx context.Context, countryCode string, limit int) ([]domain.Weather, error) {
	weathers, err := s.GetWeatherByCountry(ctx, countryCode)
	if err != nil {
		return nil, fmt.Errorf("get weather for country %q: %w", countryCode, err)
	}

	if limit < 0 {
		return nil, fmt.Errorf("limit must be non-negative")
	}

	sort.Slice(weathers, func(i, j int) bool {
		return weathers[i].Conditions.Temperature > weathers[j].Conditions.Temperature
	})

	if limit > len(weathers) {
		limit = len(weathers)
	}

	return weathers[:limit], nil
}
