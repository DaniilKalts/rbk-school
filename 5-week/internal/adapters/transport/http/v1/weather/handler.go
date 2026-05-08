package weather

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/helpers"
	domainhistory "github.com/DaniilKalts/rbk-school/5-week/internal/domain/history"
	domainuser "github.com/DaniilKalts/rbk-school/5-week/internal/domain/user"
	domainweather "github.com/DaniilKalts/rbk-school/5-week/internal/domain/weather"
)

type Service interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domainweather.Weather, error)
	GetHistory(ctx context.Context, userID uuid.UUID, city string, limit int, offset int) ([]domainhistory.History, error)
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
	case errors.Is(err, domainuser.ErrNotFound):
		status, msg = http.StatusNotFound, err.Error()
	case errors.Is(err, domainuser.ErrInvalidID),
		errors.Is(err, domainweather.ErrInvalidCity),
		errors.Is(err, domainweather.ErrInvalidLatitude),
		errors.Is(err, domainweather.ErrInvalidLongitude),
		errors.Is(err, domainweather.ErrInvalidLimit),
		errors.Is(err, domainweather.ErrInvalidOffset):
		status, msg = http.StatusBadRequest, err.Error()
	}
	helpers.JSON(w, status, helpers.NewErrorResponse(status, msg))
}
