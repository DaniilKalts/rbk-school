package user

import (
	"github.com/go-chi/chi/v5"
)

func RegisterCurrentUserRoutes(r chi.Router, service Service, tokenRevoker TokenRevoker) {
	h := NewHandler(service, tokenRevoker)
	r.Get("/users/me", h.Me())
	r.Patch("/users/me", h.MePatch())
	r.Delete("/users/me", h.MeDelete())
}

func RegisterAdminRoutes(r chi.Router, service Service, tokenRevoker TokenRevoker) {
	h := NewHandler(service, tokenRevoker)
	r.Post("/users", h.Post())
	r.Get("/users", h.Get())
	r.Get("/users/{id}", h.GetByID())
	r.Patch("/users/{id}", h.Patch())
	r.Delete("/users/{id}", h.Delete())
}
