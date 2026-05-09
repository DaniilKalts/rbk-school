package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	domainuser "github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
)

var ErrInvalidCredentials = errors.New("неверный email или пароль")

type Repository interface {
	Create(ctx context.Context, u domainuser.User, password domainuser.Password) (*domainuser.User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*Credentials, error)
}

type Credentials struct {
	ID       uuid.UUID
	Email    string
	Role     domainuser.Role
	Password domainuser.Password
}

func (c Credentials) Verify(plain string) bool {
	return c.Password.Matches(plain)
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

func (s *Service) Register(ctx context.Context, input RegisterInput) (*Token, error) {
	password, err := domainuser.NewPassword(input.Password)
	if err != nil {
		return nil, err
	}

	u, err := domainuser.NewUser(input.FirstName, input.LastName, input.Email, domainuser.RoleUser)
	if err != nil {
		return nil, err
	}

	created, err := s.repository.Create(ctx, *u, password)
	if err != nil {
		return nil, err
	}

	return s.generateToken(created.ID, created.Email, created.Role)
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*Token, error) {
	email := domainuser.NormalizeEmail(input.Email)
	if email == "" || input.Password == "" {
		return nil, ErrInvalidCredentials
	}

	credentials, err := s.repository.GetCredentialsByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domainuser.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if !credentials.Verify(input.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.generateToken(credentials.ID, credentials.Email, credentials.Role)
}

func (s *Service) Logout(ctx context.Context, accessToken string) error {
	return s.tokenManager.Revoke(ctx, accessToken)
}

func (s *Service) generateToken(userID uuid.UUID, email string, role domainuser.Role) (*Token, error) {
	accessToken, expiresAt, err := s.tokenManager.Generate(userID, email, string(role))
	if err != nil {
		return nil, err
	}

	return &Token{AccessToken: accessToken, ExpiresAt: expiresAt}, nil
}
