package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
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
