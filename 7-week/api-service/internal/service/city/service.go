package city

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/city"
	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/user"
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

	if _, err := s.userRepository.GetByID(ctx, userID); err != nil {
		return nil, err
	}

	return s.repository.Create(ctx, *c)
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]city.City, error) {
	if userID == uuid.Nil {
		return nil, city.ErrInvalidUserID
	}

	if _, err := s.userRepository.GetByID(ctx, userID); err != nil {
		return nil, err
	}

	return s.repository.ListByUserID(ctx, userID)
}

func (s *Service) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	if id == uuid.Nil {
		return city.ErrInvalidID
	}

	if userID == uuid.Nil {
		return city.ErrInvalidUserID
	}

	if _, err := s.userRepository.GetByID(ctx, userID); err != nil {
		return err
	}

	return s.repository.Delete(ctx, userID, id)
}
