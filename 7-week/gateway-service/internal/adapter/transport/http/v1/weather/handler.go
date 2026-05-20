package weather

import (
	"context"
	"net/http"
	"strings"

	"github.com/DaniilKalts/rbk-school/7-week/pkg/httpx"

	weathersvc "github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/service/weather"
)

type Service interface {
	GetByCity(ctx context.Context, city string) (weathersvc.Snapshot, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetByCity(w http.ResponseWriter, r *http.Request) {
	city := strings.TrimSpace(r.URL.Query().Get("city"))
	if city == "" {
		httpx.WriteError(w, http.StatusBadRequest, "параметр city обязателен")
		return
	}

	snapshot, err := h.service.GetByCity(r.Context(), city)
	if err != nil {
		httpx.WriteInternalError(w, r, err)
		return
	}

	httpx.JSON(w, http.StatusOK, Response{
		City:                snapshot.City,
		Latitude:            snapshot.Latitude,
		Longitude:           snapshot.Longitude,
		Temperature:         snapshot.Temperature,
		ApparentTemperature: snapshot.ApparentTemperature,
		WeatherCode:         snapshot.WeatherCode,
	})
}
