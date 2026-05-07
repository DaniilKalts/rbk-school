package auth

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Post("/auth/register", h.Post())
	r.Post("/auth/login", h.Login())
	r.Post("/auth/logout", h.Delete())
}
