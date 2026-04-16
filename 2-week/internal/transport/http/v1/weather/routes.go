package weather

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/{city}", h.GetWeatherByCity)
	r.Get("/country/{country}", h.GetWeatherByCountry)
	r.Get("/country/{country}/top", h.GetTopWarmestCities)
}
