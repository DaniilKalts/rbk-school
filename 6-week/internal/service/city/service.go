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
