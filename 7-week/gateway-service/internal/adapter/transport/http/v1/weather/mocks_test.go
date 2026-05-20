package weather_test

import (
	"context"

	"github.com/stretchr/testify/mock"

	weathersvc "github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/service/weather"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) GetByCity(ctx context.Context, city string) (weathersvc.Snapshot, error) {
	args := m.Called(ctx, city)
	out, _ := args.Get(0).(weathersvc.Snapshot)
	return out, args.Error(1)
}
