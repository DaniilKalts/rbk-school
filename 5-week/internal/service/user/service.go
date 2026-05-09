package user

import (
	"context"

	"github.com/google/uuid"

	domainuser "github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
)

type Repository interface {
	Create(ctx context.Context, u domainuser.User, password domainuser.Password) (*domainuser.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error)
	GetByEmail(ctx context.Context, email string) (*domainuser.User, error)
	List(ctx context.Context) ([]domainuser.User, error)
	Update(ctx context.Context, u domainuser.User) (*domainuser.User, error)
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

func (s *Service) Create(ctx context.Context, input CreateInput) (*domainuser.User, error) {
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

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error) {
	if id == uuid.Nil {
		return nil, domainuser.ErrInvalidID
	}

	u, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	email = domainuser.NormalizeEmail(email)
	if email == "" {
		return nil, domainuser.ErrInvalidEmail
	}

	u, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) List(ctx context.Context) ([]domainuser.User, error) {
	users, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*domainuser.User, error) {
	if id == uuid.Nil {
		return nil, domainuser.ErrInvalidID
	}

	existing, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := existing.UpdateProfile(input.FirstName, input.LastName, input.Email); err != nil {
		return nil, err
	}

	updated, err := s.repository.Update(ctx, *existing)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return domainuser.ErrInvalidID
	}

	if err := s.repository.SoftDelete(ctx, id); err != nil {
		return err
	}

	return nil
}
