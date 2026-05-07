package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/helpers"

	domainuser "github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
	serviceuser "github.com/DaniilKalts/rbk-school/5-week/internal/service/user"
)

type Service interface {
	Create(ctx context.Context, input serviceuser.CreateInput) (*domainuser.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error)
	List(ctx context.Context) ([]domainuser.User, error)
	Update(ctx context.Context, id uuid.UUID, input serviceuser.UpdateInput) (*domainuser.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func WriteServiceError(w http.ResponseWriter, err error) {
	status, msg := http.StatusInternalServerError, "internal server error"
	switch {
	case errors.Is(err, domainuser.ErrNotFound):
		status, msg = http.StatusNotFound, err.Error()
	case errors.Is(err, domainuser.ErrEmailAlreadyExists):
		status, msg = http.StatusConflict, err.Error()
	case errors.Is(err, domainuser.ErrInvalidID), errors.Is(err, domainuser.ErrInvalidFirstName), errors.Is(err, domainuser.ErrInvalidLastName), errors.Is(err, domainuser.ErrInvalidEmail), errors.Is(err, domainuser.ErrInvalidPassword), errors.Is(err, domainuser.ErrInvalidRole):
		status, msg = http.StatusBadRequest, err.Error()
	}
	helpers.JSON(w, status, helpers.NewErrorResponse(status, msg))
}
