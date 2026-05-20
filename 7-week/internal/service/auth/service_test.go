package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainuser "github.com/DaniilKalts/rbk-school/7-week/internal/domain/user"
	"github.com/DaniilKalts/rbk-school/7-week/internal/service/auth"
)

func newService(t *testing.T) (*mockRepository, *mockTokenManager, *auth.Service) {
	t.Helper()

	repo := new(mockRepository)
	tm := new(mockTokenManager)
	t.Cleanup(func() {
		repo.AssertExpectations(t)
		tm.AssertExpectations(t)
	})

	return repo, tm, auth.NewService(repo, tm)
}

func TestService_Register(t *testing.T) {
	id := uuid.New()
	expiresAt := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)
	created := &domainuser.User{
		ID:        id,
		FirstName: "Daniil",
		LastName:  "Kalts",
		Email:     "daniil@example.com",
		Role:      domainuser.RoleUser,
	}
	validInput := auth.RegisterInput{
		FirstName: "Daniil",
		LastName:  "Kalts",
		Email:     "daniil@example.com",
		Password:  "secret",
	}
	tokenErr := errors.New("signing failed")

	tests := []struct {
		name      string
		input     auth.RegisterInput
		setupMock func(*mockRepository, *mockTokenManager)
		wantToken *auth.Token
		wantErr   error
	}{
		{
			name:  "success",
			input: validInput,
			setupMock: func(r *mockRepository, tm *mockTokenManager) {
				r.On("Create", mock.Anything, mock.Anything, mock.Anything).
					Return(created, nil).Once()
				tm.On("Generate", id, "daniil@example.com", string(domainuser.RoleUser)).
					Return("token-abc", expiresAt, nil).Once()
			},
			wantToken: &auth.Token{AccessToken: "token-abc", ExpiresAt: expiresAt},
		},
		{
			name: "empty password",
			input: auth.RegisterInput{
				FirstName: "Daniil",
				LastName:  "Kalts",
				Email:     "daniil@example.com",
				Password:  "",
			},
			wantErr: domainuser.ErrInvalidPassword,
		},
		{
			name: "invalid email",
			input: auth.RegisterInput{
				FirstName: "Daniil",
				LastName:  "Kalts",
				Email:     "not-an-email",
				Password:  "secret",
			},
			wantErr: domainuser.ErrInvalidEmail,
		},
		{
			name:  "email already exists",
			input: validInput,
			setupMock: func(r *mockRepository, _ *mockTokenManager) {
				r.On("Create", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, domainuser.ErrEmailAlreadyExists).Once()
			},
			wantErr: domainuser.ErrEmailAlreadyExists,
		},
		{
			name:  "token generation error",
			input: validInput,
			setupMock: func(r *mockRepository, tm *mockTokenManager) {
				r.On("Create", mock.Anything, mock.Anything, mock.Anything).
					Return(created, nil).Once()
				tm.On("Generate", id, "daniil@example.com", string(domainuser.RoleUser)).
					Return("", time.Time{}, tokenErr).Once()
			},
			wantErr: tokenErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, tm, svc := newService(t)
			if tt.setupMock != nil {
				tt.setupMock(repo, tm)
			}

			token, err := svc.Register(context.Background(), tt.input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, token)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantToken, token)
		})
	}
}

func TestService_Login(t *testing.T) {
	id := uuid.New()
	expiresAt := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)
	password, err := domainuser.NewPassword("secret")
	require.NoError(t, err)
	user := &domainuser.User{
		ID:    id,
		Email: "daniil@example.com",
		Role:  domainuser.RoleUser,
	}
	repoErr := errors.New("db down")
	tokenErr := errors.New("signing failed")

	tests := []struct {
		name      string
		input     auth.LoginInput
		setupMock func(*mockRepository, *mockTokenManager)
		wantToken *auth.Token
		wantErr   error
	}{
		{
			name:  "success normalizes email",
			input: auth.LoginInput{Email: "DANIIL@example.com", Password: "secret"},
			setupMock: func(r *mockRepository, tm *mockTokenManager) {
				r.On("GetCredentialsByEmail", mock.Anything, "daniil@example.com").
					Return(user, password, nil).Once()
				tm.On("Generate", id, "daniil@example.com", string(domainuser.RoleUser)).
					Return("token-abc", expiresAt, nil).Once()
			},
			wantToken: &auth.Token{AccessToken: "token-abc", ExpiresAt: expiresAt},
		},
		{
			name:    "blank email",
			input:   auth.LoginInput{Email: "   ", Password: "secret"},
			wantErr: auth.ErrInvalidCredentials,
		},
		{
			name:    "empty password",
			input:   auth.LoginInput{Email: "daniil@example.com", Password: ""},
			wantErr: auth.ErrInvalidCredentials,
		},
		{
			name:  "user not found maps to invalid credentials",
			input: auth.LoginInput{Email: "daniil@example.com", Password: "secret"},
			setupMock: func(r *mockRepository, _ *mockTokenManager) {
				r.On("GetCredentialsByEmail", mock.Anything, "daniil@example.com").
					Return(nil, domainuser.Password{}, domainuser.ErrNotFound).Once()
			},
			wantErr: auth.ErrInvalidCredentials,
		},
		{
			name:  "repo error propagates",
			input: auth.LoginInput{Email: "daniil@example.com", Password: "secret"},
			setupMock: func(r *mockRepository, _ *mockTokenManager) {
				r.On("GetCredentialsByEmail", mock.Anything, "daniil@example.com").
					Return(nil, domainuser.Password{}, repoErr).Once()
			},
			wantErr: repoErr,
		},
		{
			name:  "wrong password",
			input: auth.LoginInput{Email: "daniil@example.com", Password: "wrong"},
			setupMock: func(r *mockRepository, _ *mockTokenManager) {
				r.On("GetCredentialsByEmail", mock.Anything, "daniil@example.com").
					Return(user, password, nil).Once()
			},
			wantErr: auth.ErrInvalidCredentials,
		},
		{
			name:  "token generation error after successful auth",
			input: auth.LoginInput{Email: "daniil@example.com", Password: "secret"},
			setupMock: func(r *mockRepository, tm *mockTokenManager) {
				r.On("GetCredentialsByEmail", mock.Anything, "daniil@example.com").
					Return(user, password, nil).Once()
				tm.On("Generate", id, "daniil@example.com", string(domainuser.RoleUser)).
					Return("", time.Time{}, tokenErr).Once()
			},
			wantErr: tokenErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, tm, svc := newService(t)
			if tt.setupMock != nil {
				tt.setupMock(repo, tm)
			}

			token, err := svc.Login(context.Background(), tt.input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, token)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantToken, token)
		})
	}
}

func TestService_Logout(t *testing.T) {
	revokeErr := errors.New("redis down")

	tests := []struct {
		name      string
		token     string
		setupMock func(*mockTokenManager)
		wantErr   error
	}{
		{
			name:  "success",
			token: "token-abc",
			setupMock: func(tm *mockTokenManager) {
				tm.On("Revoke", mock.Anything, "token-abc").Return(nil).Once()
			},
		},
		{
			name:  "revoke error propagates",
			token: "token-abc",
			setupMock: func(tm *mockTokenManager) {
				tm.On("Revoke", mock.Anything, "token-abc").Return(revokeErr).Once()
			},
			wantErr: revokeErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, tm, svc := newService(t)
			tt.setupMock(tm)

			err := svc.Logout(context.Background(), tt.token)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
