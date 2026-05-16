package city

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/city"

	servicecity "github.com/DaniilKalts/rbk-school/5-week/internal/service/city"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, input servicecity.CreateInput) (*city.City, error)
	List(ctx context.Context, userID uuid.UUID) ([]city.City, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}
