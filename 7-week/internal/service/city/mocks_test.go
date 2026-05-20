package city_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	domaincity "github.com/DaniilKalts/rbk-school/7-week/internal/domain/city"
	domainuser "github.com/DaniilKalts/rbk-school/7-week/internal/domain/user"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Create(ctx context.Context, c domaincity.City) (*domaincity.City, error) {
	args := m.Called(ctx, c)
	out, _ := args.Get(0).(*domaincity.City)
	return out, args.Error(1)
}

func (m *mockRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domaincity.City, error) {
	args := m.Called(ctx, userID)
	out, _ := args.Get(0).([]domaincity.City)
	return out, args.Error(1)
}

func (m *mockRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	args := m.Called(ctx, userID, id)
	return args.Error(0)
}

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error) {
	args := m.Called(ctx, id)
	out, _ := args.Get(0).(*domainuser.User)
	return out, args.Error(1)
}
