package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
)

type Repository interface {
	Create(ctx context.Context, u user.User, password user.Password) (*user.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	List(ctx context.Context) ([]user.User, error)
	Update(ctx context.Context, u user.User) (*user.User, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repository Repository
}

type CreateInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type UpdateInput struct {
	FirstName string
	LastName  string
	Email     string
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*user.User, error) {
	password, err := user.NewPassword(input.Password)
	if err != nil {
		return nil, err
	}

	u, err := user.NewUser(input.FirstName, input.LastName, input.Email, user.RoleUser)
	if err != nil {
		return nil, err
	}

	return s.repository.Create(ctx, *u, password)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	if id == uuid.Nil {
		return nil, user.ErrInvalidID
	}

	return s.repository.GetByID(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	email = user.NormalizeEmail(email)
	if email == "" {
		return nil, user.ErrInvalidEmail
	}

	return s.repository.GetByEmail(ctx, email)
}

func (s *Service) List(ctx context.Context) ([]user.User, error) {
	return s.repository.List(ctx)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*user.User, error) {
	if id == uuid.Nil {
		return nil, user.ErrInvalidID
	}

	existing, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := existing.UpdateProfile(input.FirstName, input.LastName, input.Email); err != nil {
		return nil, err
	}

	return s.repository.Update(ctx, *existing)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return user.ErrInvalidID
	}

	return s.repository.SoftDelete(ctx, id)
}
