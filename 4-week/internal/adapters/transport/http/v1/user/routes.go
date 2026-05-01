package user

import "github.com/go-chi/chi/v5"

func RegisterCurrentUserRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Get("/users/me", h.Me)
}

func RegisterAdminRoutes(r chi.Router, service Service) {
	h := NewHandler(service)

	r.Get("/users", h.List)
	r.Get("/users/{id}", h.GetByID)
	r.Delete("/users/{id}", h.Delete)
}
