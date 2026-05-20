package city

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Post("/cities", h.Create)
	r.Get("/cities", h.List)
	r.Delete("/cities/{city_id}", h.Delete)
}
