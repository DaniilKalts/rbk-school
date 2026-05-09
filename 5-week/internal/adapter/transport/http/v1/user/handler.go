package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapter/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
	"github.com/DaniilKalts/rbk-school/5-week/internal/utils"

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

func WriteServiceError(w http.ResponseWriter, err error) {
	status, msg := http.StatusInternalServerError, "internal server error"
	switch {
	case errors.Is(err, user.ErrNotFound):
		status, msg = http.StatusNotFound, err.Error()
	case errors.Is(err, user.ErrEmailAlreadyExists):
		status, msg = http.StatusConflict, err.Error()
	case errors.Is(err, user.ErrInvalidID), errors.Is(err, user.ErrInvalidFirstName), errors.Is(err, user.ErrInvalidLastName), errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrInvalidPassword), errors.Is(err, user.ErrInvalidRole):
		status, msg = http.StatusBadRequest, err.Error()
	case errors.Is(err, utils.ErrInvalidToken):
		status, msg = http.StatusUnauthorized, err.Error()
	}
	helpers.JSON(w, status, helpers.NewErrorResponse(status, msg))
}
