package weather

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	geocodingdto "github.com/DaniilKalts/rbk-school/3-week/internal/adapters/client/geocoding/dto"
	openmeteodto "github.com/DaniilKalts/rbk-school/3-week/internal/adapters/client/openmeteo/dto"
	domaincity "github.com/DaniilKalts/rbk-school/3-week/internal/domain/city"
	domainuser "github.com/DaniilKalts/rbk-school/3-week/internal/domain/user"
	domainweather "github.com/DaniilKalts/rbk-school/3-week/internal/domain/weather"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error)
}

type CityRepository interface {
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]domaincity.City, error)
}

type HistoryRepository interface {
	CreateHistory(ctx context.Context, history domainweather.History) (*domainweather.History, error)
}

type GeocodingClient interface {
	GetCoordsByCity(ctx context.Context, city string) (geocodingdto.CoordsResponse, error)
}

type WeatherClient interface {
	GetWeatherByCoords(ctx context.Context, latitude, longitude float64) (openmeteodto.WeatherResponse, error)
}

type Service struct {
	userRepository    UserRepository
	cityRepository    CityRepository
	historyRepository HistoryRepository
	geocodingClient   GeocodingClient
	weatherClient     WeatherClient
}

func New(
	userRepository UserRepository,
	cityRepository CityRepository,
	historyRepository HistoryRepository,
	geocodingClient GeocodingClient,
	weatherClient WeatherClient,
) *Service {
	return &Service{
		userRepository:    userRepository,
		cityRepository:    cityRepository,
		historyRepository: historyRepository,
		geocodingClient:   geocodingClient,
		weatherClient:     weatherClient,
	}
}

func (s *Service) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domainweather.Weather, error) {
	if userID == uuid.Nil {
		return nil, domainuser.ErrInvalidID
	}

	_, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	cities, err := s.cityRepository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list user cities: %w", err)
	}

	weathers := make([]domainweather.Weather, 0, len(cities))
	for _, city := range cities {
		weather, err := s.getWeatherByCity(ctx, city.Name)
		if err != nil {
			return nil, err
		}

		history, err := s.historyRepository.CreateHistory(ctx, domainweather.History{
			ID:          uuid.New(),
			UserID:      userID,
			City:        weather.City,
			Temperature: weather.Temperature,
			Description: weather.Description,
		})
		if err != nil {
			return nil, err
		}

		weather.RequestedAt = history.RequestedAt
		weathers = append(weathers, weather)
	}

	return weathers, nil
}

func (s *Service) getWeatherByCity(ctx context.Context, city string) (domainweather.Weather, error) {
	coords, err := s.geocodingClient.GetCoordsByCity(ctx, city)
	if err != nil {
		return domainweather.Weather{}, fmt.Errorf("get coordinates for city %q: %w", city, err)
	}

	weatherResponse, err := s.weatherClient.GetWeatherByCoords(ctx, coords.Latitude, coords.Longitude)
	if err != nil {
		return domainweather.Weather{}, fmt.Errorf("get weather for city %q: %w", city, err)
	}

	weather, err := domainweather.New(
		city,
		coords.Latitude,
		coords.Longitude,
		weatherResponse.Current.Temperature2M,
		weatherResponse.Current.ApparentTemperature,
		weatherResponse.Current.WeatherCode,
	)
	if err != nil {
		return domainweather.Weather{}, fmt.Errorf("build weather for city %q: %w", city, err)
	}

	return weather, nil
}
