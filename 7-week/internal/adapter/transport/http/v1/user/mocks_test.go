package user_test

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	domainuser "github.com/DaniilKalts/rbk-school/7-week/internal/domain/user"
	serviceuser "github.com/DaniilKalts/rbk-school/7-week/internal/service/user"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) Create(ctx context.Context, input serviceuser.CreateInput) (*domainuser.User, error) {
	args := m.Called(ctx, input)
	out, _ := args.Get(0).(*domainuser.User)
	return out, args.Error(1)
}

func (m *mockService) GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error) {
	args := m.Called(ctx, id)
	out, _ := args.Get(0).(*domainuser.User)
	return out, args.Error(1)
}

func (m *mockService) List(ctx context.Context) ([]domainuser.User, error) {
	args := m.Called(ctx)
	out, _ := args.Get(0).([]domainuser.User)
	return out, args.Error(1)
}

func (m *mockService) Update(ctx context.Context, id uuid.UUID, input serviceuser.UpdateInput) (*domainuser.User, error) {
	args := m.Called(ctx, id, input)
	out, _ := args.Get(0).(*domainuser.User)
	return out, args.Error(1)
}

func (m *mockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockTokenRevoker struct {
	mock.Mock
}

func (m *mockTokenRevoker) Revoke(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
