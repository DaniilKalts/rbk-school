package weather

import (
	"context"

	"github.com/google/uuid"

	domaincity "github.com/DaniilKalts/rbk-school/4-week/internal/domain/city"
	domainhistory "github.com/DaniilKalts/rbk-school/4-week/internal/domain/history"
	domainuser "github.com/DaniilKalts/rbk-school/4-week/internal/domain/user"
	domainweather "github.com/DaniilKalts/rbk-school/4-week/internal/domain/weather"

	geocodingdto "github.com/DaniilKalts/rbk-school/4-week/internal/client/geocoding/dto"
	openmeteodto "github.com/DaniilKalts/rbk-school/4-week/internal/client/openmeteo/dto"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error)
}

type CityRepository interface {
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]domaincity.City, error)
}

type HistoryRepository interface {
	CreateHistory(ctx context.Context, history domainhistory.History) (*domainhistory.History, error)
	ListHistory(ctx context.Context, userID uuid.UUID, city string, limit int, offset int) ([]domainhistory.History, error)
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
