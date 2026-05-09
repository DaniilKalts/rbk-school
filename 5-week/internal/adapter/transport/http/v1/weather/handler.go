package weather

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapter/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/history"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/weather"
)

type Service interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]weather.Weather, error)
	GetHistory(ctx context.Context, userID uuid.UUID, city string, limit int, offset int) ([]history.History, error)
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
	case errors.Is(err, user.ErrNotFound):
		status, msg = http.StatusNotFound, err.Error()
	case errors.Is(err, user.ErrInvalidID),
		errors.Is(err, weather.ErrInvalidCity),
		errors.Is(err, weather.ErrInvalidLimit),
		errors.Is(err, weather.ErrInvalidOffset):
		status, msg = http.StatusBadRequest, err.Error()
	}
	helpers.JSON(w, status, helpers.NewErrorResponse(status, msg))
}
