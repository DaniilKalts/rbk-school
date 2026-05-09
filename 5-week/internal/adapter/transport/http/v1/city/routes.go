package city

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Post("/cities", h.Post())
	r.Get("/cities", h.Get())
	r.Delete("/cities/{city_id}", h.Delete())
}
