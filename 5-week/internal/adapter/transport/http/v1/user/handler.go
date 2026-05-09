package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"

	serviceuser "github.com/DaniilKalts/rbk-school/5-week/internal/service/user"
)

type Service interface {
	Create(ctx context.Context, input serviceuser.CreateInput) (*user.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
	List(ctx context.Context) ([]user.User, error)
	Update(ctx context.Context, id uuid.UUID, input serviceuser.UpdateInput) (*user.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type TokenRevoker interface {
	Revoke(ctx context.Context, token string) error
}

type Handler struct {
	service      Service
	tokenRevoker TokenRevoker
}

func NewHandler(service Service, tokenRevoker TokenRevoker) *Handler {
	return &Handler{service: service, tokenRevoker: tokenRevoker}
}
