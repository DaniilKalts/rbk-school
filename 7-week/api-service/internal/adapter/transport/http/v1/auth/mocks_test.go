package auth_test

import (
	"context"

	"github.com/stretchr/testify/mock"

	serviceauth "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/service/auth"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) Register(ctx context.Context, input serviceauth.RegisterInput) (*serviceauth.Token, error) {
	args := m.Called(ctx, input)
	out, _ := args.Get(0).(*serviceauth.Token)
	return out, args.Error(1)
}

func (m *mockService) Login(ctx context.Context, input serviceauth.LoginInput) (*serviceauth.Token, error) {
	args := m.Called(ctx, input)
	out, _ := args.Get(0).(*serviceauth.Token)
	return out, args.Error(1)
}

func (m *mockService) Logout(ctx context.Context, accessToken string) error {
	args := m.Called(ctx, accessToken)
	return args.Error(0)
}
