package weather

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapter/client/geocoding/dto"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/city"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/history"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/weather"

	openmeteodto "github.com/DaniilKalts/rbk-school/5-week/internal/adapter/client/openmeteo/dto"
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
