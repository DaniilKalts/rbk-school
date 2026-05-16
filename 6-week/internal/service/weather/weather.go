package weather

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"

	domaincity "github.com/DaniilKalts/rbk-school/6-week/internal/domain/city"
	domainhistory "github.com/DaniilKalts/rbk-school/6-week/internal/domain/history"
	domainweather "github.com/DaniilKalts/rbk-school/6-week/internal/domain/weather"
)

func (s *Service) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domainweather.Weather, error) {
	if userID == uuid.Nil {
		return nil, user.ErrInvalidID
	}

	_, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	cities, err := s.cityRepository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("получение списка городов пользователя: %w", err)
	}

	weathers := make([]domainweather.Weather, len(cities))
	g, gCtx := errgroup.WithContext(ctx)

	for i, city := range cities {
		idx := i
		cityName := city.Name

		g.Go(func() error {
			weather, err := s.getWeatherByCity(gCtx, cityName)
			if err != nil {
				return err
			}

			history, err := s.createHistory(gCtx, userID, weather)
			if err != nil {
				return err
			}

			weather.RequestedAt = history.RequestedAt
			weathers[idx] = weather

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return weathers, nil
}

func (s *Service) getWeatherByCity(ctx context.Context, city string) (domainweather.Weather, error) {
	cacheKey := domaincity.NormalizeCityName(city)
	if s.weatherCache != nil {
		weather, ok, err := s.weatherCache.Get(ctx, cacheKey)
		if err != nil {
			log.Printf("weather cache get %q: %v", cacheKey, err)
		} else if ok {
			return weather, nil
		}
	}

	coords, err := s.geocodingClient.GetCoordsByCity(ctx, city)
	if err != nil {
		return domainweather.Weather{}, fmt.Errorf("получение координат для города %q: %w", city, err)
	}

	weatherResponse, err := s.weatherClient.GetWeatherByCoords(ctx, coords.Latitude, coords.Longitude)
	if err != nil {
		return domainweather.Weather{}, fmt.Errorf("получение погоды для города %q: %w", city, err)
	}

	weather, err := domainweather.NewWeather(
		city,
		weatherResponse.Current.Temperature2M,
		weatherResponse.Current.ApparentTemperature,
		weatherResponse.Current.WeatherCode,
	)
	if err != nil {
		return domainweather.Weather{}, fmt.Errorf("сборка модели погоды для города %q: %w", city, err)
	}

	if s.weatherCache != nil {
		if err := s.weatherCache.Set(ctx, cacheKey, weather); err != nil {
			log.Printf("weather cache set %q: %v", cacheKey, err)
		}
	}

	return weather, nil
}

func (s *Service) createHistory(ctx context.Context, userID uuid.UUID, weather domainweather.Weather) (domainhistory.History, error) {
	historyModel, err := domainhistory.NewHistory(userID, weather.City, weather.Temperature, weather.Description)
	if err != nil {
		return domainhistory.History{}, err
	}

	history, err := s.historyRepository.CreateHistory(ctx, *historyModel)
	if err != nil {
		return domainhistory.History{}, err
	}

	return *history, nil
}
