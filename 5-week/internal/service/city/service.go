package city

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/city"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
)

type Repository interface {
	Create(ctx context.Context, c city.City) (*city.City, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]city.City, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
}

type Service struct {
	repository     Repository
	userRepository UserRepository
}

type CreateInput struct {
	Name string
}

func NewService(repository Repository, userRepository UserRepository) *Service {
	return &Service{repository: repository, userRepository: userRepository}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, input CreateInput) (*city.City, error) {
	if userID == uuid.Nil {
		return nil, city.ErrInvalidUserID
	}

	c, err := city.NewCity(userID, input.Name)
	if err != nil {
		return nil, err
	}

	_, err = s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	created, err := s.repository.Create(ctx, *c)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]city.City, error) {
	if userID == uuid.Nil {
		return nil, city.ErrInvalidUserID
	}

	_, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	cities, err := s.repository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return cities, nil
}

func (s *Service) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	if id == uuid.Nil {
		return city.ErrInvalidID
	}

	if userID == uuid.Nil {
		return city.ErrInvalidUserID
	}

	_, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, userID, id)
}
