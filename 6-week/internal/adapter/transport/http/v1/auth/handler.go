package auth

import (
	"context"

	"github.com/DaniilKalts/rbk-school/5-week/internal/service/auth"
)

type Service interface {
	Register(ctx context.Context, input auth.RegisterInput) (*auth.Token, error)
	Login(ctx context.Context, input auth.LoginInput) (*auth.Token, error)
	Logout(ctx context.Context, accessToken string) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}
