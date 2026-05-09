package city

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapter/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/city"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"

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

func CurrentUserID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	claims, ok := helpers.ClaimsFromContext(r.Context())
	if !ok {
		helpers.JSON(w, http.StatusUnauthorized, helpers.NewErrorResponse(http.StatusUnauthorized, "отсутствуют claims аутентификации"))
		return uuid.Nil, false
	}
	return claims.UserID, true
}

func WriteServiceError(w http.ResponseWriter, err error) {
	status, msg := http.StatusInternalServerError, "internal server error"
	switch {
	case errors.Is(err, user.ErrNotFound), errors.Is(err, city.ErrNotFound):
		status, msg = http.StatusNotFound, err.Error()
	case errors.Is(err, city.ErrAlreadyExists):
		status, msg = http.StatusConflict, err.Error()
	case errors.Is(err, city.ErrInvalidID), errors.Is(err, city.ErrInvalidUserID), errors.Is(err, city.ErrInvalidName), errors.Is(err, user.ErrInvalidID):
		status, msg = http.StatusBadRequest, err.Error()
	}
	helpers.JSON(w, status, helpers.NewErrorResponse(status, msg))
}
