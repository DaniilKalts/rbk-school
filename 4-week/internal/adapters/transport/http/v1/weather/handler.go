package weather

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/helpers"
	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/v1/weather/dto"
	domainhistory "github.com/DaniilKalts/rbk-school/4-week/internal/domain/history"
	domainweather "github.com/DaniilKalts/rbk-school/4-week/internal/domain/weather"
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

	helpers.JSON(w, http.StatusOK, dto.ToUserWeatherResponse(userID, weathers))
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := parseUUID(w, r.PathValue("id"), "invalid user id")
	if !ok {
		return
	}

	city := strings.TrimSpace(r.URL.Query().Get("city"))
	limit, ok := parseLimit(w, r.URL.Query().Get("limit"))
	if !ok {
		return
	}

	offset, ok := parseOffset(w, r.URL.Query().Get("offset"))
	if !ok {
		return
	}

	history, err := h.service.GetHistory(r.Context(), userID, city, limit, offset)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	helpers.JSON(w, http.StatusOK, dto.ToUserWeatherHistoryResponse(userID, domainweather.NormalizeCityName(city), history))
}
