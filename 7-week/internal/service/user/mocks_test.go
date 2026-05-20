package user_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	domainuser "github.com/DaniilKalts/rbk-school/7-week/internal/domain/user"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Create(ctx context.Context, u domainuser.User, password domainuser.Password) (*domainuser.User, error) {
	args := m.Called(ctx, u, password)
	out, _ := args.Get(0).(*domainuser.User)
	return out, args.Error(1)
}

func (m *mockRepository) GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error) {
	args := m.Called(ctx, id)
	out, _ := args.Get(0).(*domainuser.User)
	return out, args.Error(1)
}

func (m *mockRepository) GetByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	args := m.Called(ctx, email)
	out, _ := args.Get(0).(*domainuser.User)
	return out, args.Error(1)
}

func (m *mockRepository) List(ctx context.Context) ([]domainuser.User, error) {
	args := m.Called(ctx)
	out, _ := args.Get(0).([]domainuser.User)
	return out, args.Error(1)
}

func (m *mockRepository) Update(ctx context.Context, u domainuser.User) (*domainuser.User, error) {
	args := m.Called(ctx, u)
	out, _ := args.Get(0).(*domainuser.User)
	return out, args.Error(1)
}

func (m *mockRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
