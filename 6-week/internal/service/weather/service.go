package weather

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/client/geocoding/dto"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/city"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/history"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/weather"
	"github.com/DaniilKalts/rbk-school/6-week/pkg/logger"

	openmeteodto "github.com/DaniilKalts/rbk-school/6-week/internal/adapter/client/openmeteo/dto"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
}

type CityRepository interface {
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]city.City, error)
}

type HistoryRepository interface {
	CreateHistory(ctx context.Context, history history.History) (*history.History, error)
	ListHistory(ctx context.Context, userID uuid.UUID, city string, limit int, offset int) ([]history.History, error)
}

type GeocodingClient interface {
	GetCoordsByCity(ctx context.Context, city string) (dto.CoordsResponse, error)
}

type WeatherClient interface {
	GetWeatherByCoords(ctx context.Context, latitude, longitude float64) (openmeteodto.WeatherResponse, error)
}

type WeatherCache interface {
	Get(ctx context.Context, city string) (weather.Weather, bool, error)
	Set(ctx context.Context, city string, weather weather.Weather) error
}

type Service struct {
	userRepository    UserRepository
	cityRepository    CityRepository
	historyRepository HistoryRepository
	geocodingClient   GeocodingClient
	weatherClient     WeatherClient
	weatherCache      WeatherCache
}

func NewService(
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

func (s *Service) GetByUserID(ctx context.Context, userID uuid.UUID) ([]weather.Weather, error) {
	if userID == uuid.Nil {
		return nil, user.ErrInvalidID
	}

	if _, err := s.userRepository.GetByID(ctx, userID); err != nil {
		return nil, err
	}

	cities, err := s.cityRepository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("получение списка городов пользователя: %w", err)
	}

	weathers := make([]weather.Weather, len(cities))
	g, gCtx := errgroup.WithContext(ctx)

	for i, c := range cities {
		idx := i
		cityName := c.Name

		g.Go(func() error {
			w, err := s.getWeatherByCity(gCtx, cityName)
			if err != nil {
				return err
			}

			h, err := s.createHistory(gCtx, userID, w)
			if err != nil {
				return err
			}

			w.RequestedAt = h.RequestedAt
			weathers[idx] = w

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return weathers, nil
}

func (s *Service) GetHistory(ctx context.Context, userID uuid.UUID, cityName string, limit int, offset int) ([]history.History, error) {
	if userID == uuid.Nil {
		return nil, user.ErrInvalidID
	}

	cityName = city.NormalizeCityName(cityName)

	if limit < 0 {
		return nil, weather.ErrInvalidLimit
	}
	if offset < 0 {
		return nil, weather.ErrInvalidOffset
	}

	if _, err := s.userRepository.GetByID(ctx, userID); err != nil {
		return nil, err
	}

	return s.historyRepository.ListHistory(ctx, userID, cityName, limit, offset)
}

func (s *Service) getWeatherByCity(ctx context.Context, cityName string) (weather.Weather, error) {
	cacheKey := city.NormalizeCityName(cityName)
	if s.weatherCache != nil {
		cached, ok, err := s.weatherCache.Get(ctx, cacheKey)
		if err != nil {
			logger.FromContext(ctx).Warn("кеш погоды: чтение", zap.String("key", cacheKey), zap.Error(err))
		} else if ok {
			return cached, nil
		}
	}

	coords, err := s.geocodingClient.GetCoordsByCity(ctx, cityName)
	if err != nil {
		return weather.Weather{}, fmt.Errorf("получение координат для города %q: %w", cityName, err)
	}

	weatherResponse, err := s.weatherClient.GetWeatherByCoords(ctx, coords.Latitude, coords.Longitude)
	if err != nil {
		return weather.Weather{}, fmt.Errorf("получение погоды для города %q: %w", cityName, err)
	}

	w, err := weather.NewWeather(
		cityName,
		weatherResponse.Current.Temperature2M,
		weatherResponse.Current.ApparentTemperature,
		weatherResponse.Current.WeatherCode,
	)
	if err != nil {
		return weather.Weather{}, fmt.Errorf("сборка модели погоды для города %q: %w", cityName, err)
	}

	if s.weatherCache != nil {
		if err := s.weatherCache.Set(ctx, cacheKey, w); err != nil {
			logger.FromContext(ctx).Warn("кеш погоды: запись", zap.String("key", cacheKey), zap.Error(err))
		}
	}

	return w, nil
}

func (s *Service) createHistory(ctx context.Context, userID uuid.UUID, w weather.Weather) (history.History, error) {
	model, err := history.NewHistory(userID, w.City, w.Temperature, w.Description)
	if err != nil {
		return history.History{}, err
	}

	created, err := s.historyRepository.CreateHistory(ctx, *model)
	if err != nil {
		return history.History{}, err
	}

	return *created, nil
}
