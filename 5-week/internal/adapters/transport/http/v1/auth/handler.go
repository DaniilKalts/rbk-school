package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
	serviceauth "github.com/DaniilKalts/rbk-school/5-week/internal/service/auth"
	"github.com/DaniilKalts/rbk-school/5-week/internal/utils"
)

type Service interface {
	Register(ctx context.Context, input serviceauth.RegisterInput) (*serviceauth.Token, error)
	Login(ctx context.Context, input serviceauth.LoginInput) (*serviceauth.Token, error)
	Logout(ctx context.Context, accessToken string) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func WriteServiceError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	msg := "internal server error"
	switch {
	case errors.Is(err, user.ErrNotFound):
		status, msg = http.StatusNotFound, err.Error()
	case errors.Is(err, user.ErrEmailAlreadyExists):
		status, msg = http.StatusConflict, err.Error()
	case errors.Is(err, user.ErrInvalidID),
		errors.Is(err, user.ErrInvalidFirstName),
		errors.Is(err, user.ErrInvalidLastName),
		errors.Is(err, user.ErrInvalidEmail),
		errors.Is(err, user.ErrInvalidPassword),
		errors.Is(err, user.ErrInvalidRole):
		status, msg = http.StatusBadRequest, err.Error()
	case errors.Is(err, serviceauth.ErrInvalidCredentials), errors.Is(err, utils.ErrInvalidToken):
		status, msg = http.StatusUnauthorized, err.Error()
	}
	helpers.JSON(w, status, helpers.NewErrorResponse(status, msg))
}
