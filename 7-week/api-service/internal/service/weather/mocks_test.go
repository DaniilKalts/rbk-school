package weather_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	gatewaydto "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/adapter/client/gateway/dto"
	domaincity "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/city"
	domainhistory "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/history"
	domainuser "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/user"
	domainweather "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/weather"
)

type mockUserRepository struct{ mock.Mock }

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error) {
	args := m.Called(ctx, id)
	out, _ := args.Get(0).(*domainuser.User)
	return out, args.Error(1)
}

type mockCityRepository struct{ mock.Mock }

func (m *mockCityRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domaincity.City, error) {
	args := m.Called(ctx, userID)
	out, _ := args.Get(0).([]domaincity.City)
	return out, args.Error(1)
}

type mockHistoryRepository struct{ mock.Mock }

func (m *mockHistoryRepository) CreateHistory(ctx context.Context, h domainhistory.History) (*domainhistory.History, error) {
	args := m.Called(ctx, h)
	out, _ := args.Get(0).(*domainhistory.History)
	return out, args.Error(1)
}

func (m *mockHistoryRepository) ListHistory(ctx context.Context, userID uuid.UUID, city string, limit int, offset int) ([]domainhistory.History, error) {
	args := m.Called(ctx, userID, city, limit, offset)
	out, _ := args.Get(0).([]domainhistory.History)
	return out, args.Error(1)
}

type mockGatewayClient struct{ mock.Mock }

func (m *mockGatewayClient) GetWeatherByCity(ctx context.Context, city string) (gatewaydto.WeatherResponse, error) {
	args := m.Called(ctx, city)
	out, _ := args.Get(0).(gatewaydto.WeatherResponse)
	return out, args.Error(1)
}

type mockWeatherCache struct{ mock.Mock }

func (m *mockWeatherCache) Get(ctx context.Context, city string) (domainweather.Weather, bool, error) {
	args := m.Called(ctx, city)
	w, _ := args.Get(0).(domainweather.Weather)
	ok, _ := args.Get(1).(bool)
	return w, ok, args.Error(2)
}

func (m *mockWeatherCache) Set(ctx context.Context, city string, w domainweather.Weather) error {
	args := m.Called(ctx, city, w)
	return args.Error(0)
}
