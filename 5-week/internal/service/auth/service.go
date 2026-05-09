package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
)

var ErrInvalidCredentials = errors.New("неверный email или пароль")

type Repository interface {
	Create(ctx context.Context, u user.User, password user.Password) (*user.User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*user.User, user.Password, error)
}

type TokenManager interface {
	Generate(userID uuid.UUID, email string, role string) (string, time.Time, error)
	Revoke(ctx context.Context, token string) error
}

type Service struct {
	repository   Repository
	tokenManager TokenManager
}

type RegisterInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type LoginInput struct {
	Email    string
	Password string
}

type Token struct {
	AccessToken string
	ExpiresAt   time.Time
}

func NewService(repository Repository, tokenManager TokenManager) *Service {
	return &Service{repository: repository, tokenManager: tokenManager}
}

func (s *Service) generateToken(userID uuid.UUID, email string, role user.Role) (*Token, error) {
	accessToken, expiresAt, err := s.tokenManager.Generate(userID, email, string(role))
	if err != nil {
		return nil, err
	}

	return &Token{AccessToken: accessToken, ExpiresAt: expiresAt}, nil
}
