package city

import (
	"context"

	"github.com/google/uuid"

	domaincity "github.com/DaniilKalts/rbk-school/5-week/internal/domain/city"
	domainuser "github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
)

type Repository interface {
	Create(ctx context.Context, c domaincity.City) (*domaincity.City, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]domaincity.City, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error)
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

func (s *Service) Create(ctx context.Context, userID uuid.UUID, input CreateInput) (*domaincity.City, error) {
	if userID == uuid.Nil {
		return nil, domaincity.ErrInvalidUserID
	}

	c, err := domaincity.NewCity(userID, input.Name)
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

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]domaincity.City, error) {
	if userID == uuid.Nil {
		return nil, domaincity.ErrInvalidUserID
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
		return domaincity.ErrInvalidID
	}

	if userID == uuid.Nil {
		return domaincity.ErrInvalidUserID
	}

	_, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, userID, id)
}
