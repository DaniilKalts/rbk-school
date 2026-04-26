package weather

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/3-week/internal/adapters/transport/http/v1/weather/dto"
	domainuser "github.com/DaniilKalts/rbk-school/3-week/internal/domain/user"
	domainweather "github.com/DaniilKalts/rbk-school/3-week/internal/domain/weather"
	"github.com/DaniilKalts/rbk-school/3-week/internal/utils"
)

type Service interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domainweather.Weather, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userID, ok := parseUUID(w, r.PathValue("id"), "invalid user id")
	if !ok {
		return
	}

	weathers, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	utils.JSON(w, http.StatusOK, dto.ToUserWeatherResponse(userID, weathers))
}

func parseUUID(w http.ResponseWriter, value string, message string) (uuid.UUID, bool) {
	id, err := uuid.Parse(value)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, message)
		return uuid.Nil, false
	}

	return id, true
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domainuser.ErrNotFound):
		utils.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domainuser.ErrInvalidID),
		errors.Is(err, domainweather.ErrInvalidCity):
		utils.Error(w, http.StatusBadRequest, err.Error())
	default:
		utils.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
