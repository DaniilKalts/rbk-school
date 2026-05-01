package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	domainuser "github.com/DaniilKalts/rbk-school/4-week/internal/domain/user"
	"github.com/DaniilKalts/rbk-school/4-week/internal/utils"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type Repository interface {
	Create(ctx context.Context, u domainuser.User, passwordHash string, salt string) (*domainuser.User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*domainuser.Credentials, error)
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
	if err := domainuser.ValidatePassword(input.Password); err != nil {
		return nil, err
	}

	u, err := domainuser.New(uuid.New(), input.FirstName, input.LastName, input.Email, domainuser.RoleUser)
	if err != nil {
		return nil, err
	}

	salt, err := utils.GenerateSalt()
	if err != nil {
		return nil, err
	}

	passwordHash, err := utils.HashPassword(input.Password, salt)
	if err != nil {
		return nil, err
	}

	created, err := s.repository.Create(ctx, *u, passwordHash, salt)
	if err != nil {
		return nil, err
	}

	return s.generateToken(created.ID, created.Email, created.Role)
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*Token, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if email == "" {
		return nil, domainuser.ErrInvalidEmail
	}

	if err := domainuser.ValidatePassword(input.Password); err != nil {
		return nil, err
	}

	credentials, err := s.repository.GetCredentialsByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domainuser.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if !utils.VerifyPassword(input.Password, credentials.Salt, credentials.PasswordHash) {
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
