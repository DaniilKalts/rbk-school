package weather

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	weatherDTO "github.com/DaniilKalts/rbk-school/2-week/internal/client/openmeteo/dto"
)

type Service interface {
	GetWeatherByCity(city string) (weatherDTO.WeatherResponse, error)
	GetWeatherByCountry(countryCode string) ([]weatherDTO.WeatherResponse, error)
	GetTopWarmestCities(countryCode string, limit int) ([]weatherDTO.WeatherResponse, error)
}

type Handler struct {
	service Service
}

func NewWeatherHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetWeatherByCity(w http.ResponseWriter, r *http.Request) {
	city := strings.TrimSpace(chi.URLParam(r, "city"))
	if city == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "city is required"})
		return
	}

	weather, err := h.service.GetWeatherByCity(city)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "failed to get weather"})
		return
	}

	writeJSON(w, http.StatusOK, weather)
}

func (h *Handler) GetWeatherByCountry(w http.ResponseWriter, r *http.Request) {
	countryCode := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "country")))
	if countryCode == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "country code is required"})
		return
	}

	weathers, err := h.service.GetWeatherByCountry(countryCode)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "failed to get country weather"})
		return
	}

	writeJSON(w, http.StatusOK, weathers)
}

func (h *Handler) GetTopWarmestCities(w http.ResponseWriter, r *http.Request) {
	countryCode := strings.ToUpper(strings.TrimSpace(chi.URLParam(r, "country")))
	if countryCode == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "country code is required"})
		return
	}

	limit := 3

	limitParam := strings.TrimSpace(r.URL.Query().Get("limit"))
	if limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err != nil || parsedLimit <= 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"message": "limit must be a positive integer"})
			return
		}

		limit = parsedLimit
	}

	weathers, err := h.service.GetTopWarmestCities(countryCode, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "failed to get warmest cities"})
		return
	}

	writeJSON(w, http.StatusOK, weathers)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
