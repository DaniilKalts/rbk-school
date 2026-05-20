package auth_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	domainuser "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/user"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Create(ctx context.Context, u domainuser.User, password domainuser.Password) (*domainuser.User, error) {
	args := m.Called(ctx, u, password)
	out, _ := args.Get(0).(*domainuser.User)
	return out, args.Error(1)
}

func (m *mockRepository) GetCredentialsByEmail(ctx context.Context, email string) (*domainuser.User, domainuser.Password, error) {
	args := m.Called(ctx, email)
	u, _ := args.Get(0).(*domainuser.User)
	pwd, _ := args.Get(1).(domainuser.Password)
	return u, pwd, args.Error(2)
}

type mockTokenManager struct {
	mock.Mock
}

func (m *mockTokenManager) Generate(userID uuid.UUID, email string, role string) (string, time.Time, error) {
	args := m.Called(userID, email, role)
	token, _ := args.Get(0).(string)
	exp, _ := args.Get(1).(time.Time)
	return token, exp, args.Error(2)
}

func (m *mockTokenManager) Revoke(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
