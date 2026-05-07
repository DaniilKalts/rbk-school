package city

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/helpers"
	domaincity "github.com/DaniilKalts/rbk-school/5-week/internal/domain/city"
	domainuser "github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
	servicecity "github.com/DaniilKalts/rbk-school/5-week/internal/service/city"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, input servicecity.CreateInput) (*domaincity.City, error)
	List(ctx context.Context, userID uuid.UUID) ([]domaincity.City, error)
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
	case errors.Is(err, domainuser.ErrNotFound), errors.Is(err, domaincity.ErrNotFound):
		status, msg = http.StatusNotFound, err.Error()
	case errors.Is(err, domaincity.ErrAlreadyExists):
		status, msg = http.StatusConflict, err.Error()
	case errors.Is(err, domaincity.ErrInvalidID), errors.Is(err, domaincity.ErrInvalidUserID), errors.Is(err, domaincity.ErrInvalidName), errors.Is(err, domainuser.ErrInvalidID):
		status, msg = http.StatusBadRequest, err.Error()
	}
	helpers.JSON(w, status, helpers.NewErrorResponse(status, msg))
}
