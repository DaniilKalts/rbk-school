package weather_test

import (
	"context"

	"github.com/stretchr/testify/mock"

	geocodingdto "github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/adapter/client/geocoding/dto"
	openmeteodto "github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/adapter/client/openmeteo/dto"
)

type mockGeocodingClient struct{ mock.Mock }

func (m *mockGeocodingClient) GetCoordsByCity(ctx context.Context, city string) (geocodingdto.CoordsResponse, error) {
	args := m.Called(ctx, city)
	out, _ := args.Get(0).(geocodingdto.CoordsResponse)
	return out, args.Error(1)
}

type mockOpenMeteoClient struct{ mock.Mock }

func (m *mockOpenMeteoClient) GetWeatherByCoords(ctx context.Context, latitude, longitude float64) (openmeteodto.WeatherResponse, error) {
	args := m.Called(ctx, latitude, longitude)
	out, _ := args.Get(0).(openmeteodto.WeatherResponse)
	return out, args.Error(1)
}
