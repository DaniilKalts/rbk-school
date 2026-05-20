package weather_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	domainhistory "github.com/DaniilKalts/rbk-school/7-week/internal/domain/history"
	domainweather "github.com/DaniilKalts/rbk-school/7-week/internal/domain/weather"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domainweather.Weather, error) {
	args := m.Called(ctx, userID)
	out, _ := args.Get(0).([]domainweather.Weather)
	return out, args.Error(1)
}

func (m *mockService) GetHistory(ctx context.Context, userID uuid.UUID, city string, limit int, offset int) ([]domainhistory.History, error) {
	args := m.Called(ctx, userID, city, limit, offset)
	out, _ := args.Get(0).([]domainhistory.History)
	return out, args.Error(1)
}
