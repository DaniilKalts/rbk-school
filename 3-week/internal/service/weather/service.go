package weather

import (
	"context"
	"fmt"
	"sync"

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
	ListHistory(ctx context.Context, userID uuid.UUID, city string, limit int, offset int) ([]domainweather.History, error)
}

type GeocodingClient interface {
	GetCoordsByCity(ctx context.Context, city string) (geocodingdto.CoordsResponse, error)
}

type WeatherClient interface {
	GetWeatherByCoords(ctx context.Context, latitude, longitude float64) (openmeteodto.WeatherResponse, error)
}

type WeatherCache interface {
	Get(ctx context.Context, city string) (domainweather.Weather, bool, error)
	Set(ctx context.Context, city string, weather domainweather.Weather) error
}

type Service struct {
	userRepository    UserRepository
	cityRepository    CityRepository
	historyRepository HistoryRepository
	geocodingClient   GeocodingClient
	weatherClient     WeatherClient
	weatherCache      WeatherCache
}

func New(
	userRepository UserRepository,
	cityRepository CityRepository,
	historyRepository HistoryRepository,
	geocodingClient GeocodingClient,
	weatherClient WeatherClient,
	weatherCache WeatherCache,
) *Service {
	return &Service{
		userRepository:    userRepository,
		cityRepository:    cityRepository,
		historyRepository: historyRepository,
		geocodingClient:   geocodingClient,
		weatherClient:     weatherClient,
		weatherCache:      weatherCache,
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

	weathers := make([]domainweather.Weather, len(cities))
	errCh := make(chan error, len(cities))
	var wg sync.WaitGroup

	for i, city := range cities {
		wg.Add(1)
		go func(i int, city domaincity.City) {
			defer wg.Done()

			weather, err := s.getWeatherByCity(ctx, city.Name)
			if err != nil {
				errCh <- err
				return
			}

			history, err := s.historyRepository.CreateHistory(ctx, domainweather.History{
				ID:          uuid.New(),
				UserID:      userID,
				City:        weather.City,
				Temperature: weather.Temperature,
				Description: weather.Description,
			})
			if err != nil {
				errCh <- err
				return
			}

			weather.RequestedAt = history.RequestedAt
			weathers[i] = weather
		}(i, city)
	}
	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		return nil, err
	}

	return weathers, nil
}

func (s *Service) GetHistory(ctx context.Context, userID uuid.UUID, city string, limit int, offset int) ([]domainweather.History, error) {
	if userID == uuid.Nil {
		return nil, domainuser.ErrInvalidID
	}

	city = domainweather.NormalizeCityName(city)

	if limit < 0 {
		return nil, domainweather.ErrInvalidLimit
	}
	if offset < 0 {
		return nil, domainweather.ErrInvalidOffset
	}

	_, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	history, err := s.historyRepository.ListHistory(ctx, userID, city, limit, offset)
	if err != nil {
		return nil, err
	}

	return history, nil
}

func (s *Service) getWeatherByCity(ctx context.Context, city string) (domainweather.Weather, error) {
	cacheKey := domainweather.NormalizeCityName(city)
	if s.weatherCache != nil {
		weather, ok, err := s.weatherCache.Get(ctx, cacheKey)
		if err == nil && ok {
			return weather, nil
		}
	}

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

	if s.weatherCache != nil {
		_ = s.weatherCache.Set(ctx, weather.City, weather)
	}

	return weather, nil
}
