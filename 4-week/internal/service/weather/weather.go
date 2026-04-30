package weather

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	domaincity "github.com/DaniilKalts/rbk-school/3-week/internal/domain/city"
	domainhistory "github.com/DaniilKalts/rbk-school/3-week/internal/domain/history"
	domainuser "github.com/DaniilKalts/rbk-school/3-week/internal/domain/user"
	domainweather "github.com/DaniilKalts/rbk-school/3-week/internal/domain/weather"
)

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

			historyModel, err := domainhistory.New(uuid.New(), userID, weather.City, weather.Temperature, weather.Description)
			if err != nil {
				errCh <- err
				return
			}

			history, err := s.historyRepository.CreateHistory(ctx, *historyModel)
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
