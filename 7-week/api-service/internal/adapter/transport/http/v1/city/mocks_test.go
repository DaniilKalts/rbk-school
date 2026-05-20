package city_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	domaincity "github.com/DaniilKalts/rbk-school/7-week/internal/domain/city"
	servicecity "github.com/DaniilKalts/rbk-school/7-week/internal/service/city"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) Create(ctx context.Context, userID uuid.UUID, input servicecity.CreateInput) (*domaincity.City, error) {
	args := m.Called(ctx, userID, input)
	out, _ := args.Get(0).(*domaincity.City)
	return out, args.Error(1)
}

func (m *mockService) List(ctx context.Context, userID uuid.UUID) ([]domaincity.City, error) {
	args := m.Called(ctx, userID)
	out, _ := args.Get(0).([]domaincity.City)
	return out, args.Error(1)
}

func (m *mockService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	args := m.Called(ctx, userID, id)
	return args.Error(0)
}
