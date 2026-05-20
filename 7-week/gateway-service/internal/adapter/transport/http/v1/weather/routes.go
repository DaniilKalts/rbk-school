package weather

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Get("/weather", h.GetByCity)
}
