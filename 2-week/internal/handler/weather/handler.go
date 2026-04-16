package weather

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/DaniilKalts/rbk-school/2-week/internal/domain"
	handlerDTO "github.com/DaniilKalts/rbk-school/2-week/internal/handler/weather/dto"
	"github.com/go-chi/chi/v5"
)

type Service interface {
	GetWeatherByCity(ctx context.Context, city string) (domain.Weather, error)
	GetWeatherByCountry(ctx context.Context, countryCode string) ([]domain.Weather, error)
	GetTopWarmestCities(ctx context.Context, countryCode string, limit int) ([]domain.Weather, error)
}

type Handler struct {
	service Service
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewWeatherHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetWeatherByCity(w http.ResponseWriter, r *http.Request) {
	city := strings.TrimSpace(chi.URLParam(r, "city"))
	if city == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "City is required."})
		return
	}

	weather, err := h.service.GetWeatherByCity(r.Context(), city)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Could not get weather for the city."})
		return
	}

	writeJSON(w, http.StatusOK, handlerDTO.FromDomainWeather(weather))
}

func (h *Handler) GetWeatherByCountry(w http.ResponseWriter, r *http.Request) {
	countryCode := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "country")))
	if countryCode == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Country code is required."})
		return
	}

	weathers, err := h.service.GetWeatherByCountry(r.Context(), countryCode)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Could not get weather for the country."})
		return
	}

	writeJSON(w, http.StatusOK, handlerDTO.FromDomainWeathers(weathers))
}

func (h *Handler) GetTopWarmestCities(w http.ResponseWriter, r *http.Request) {
	countryCode := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "country")))
	if countryCode == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Country code is required."})
		return
	}

	limit := 3

	limitParam := strings.TrimSpace(r.URL.Query().Get("limit"))
	if limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err != nil || parsedLimit <= 0 {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Limit must be a positive number."})
			return
		}

		limit = parsedLimit
	}

	weathers, err := h.service.GetTopWarmestCities(r.Context(), countryCode, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Could not get the warmest cities for the country."})
		return
	}

	writeJSON(w, http.StatusOK, handlerDTO.FromDomainWeathers(weathers))
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
