package user

import (
	"github.com/go-chi/chi/v5"
)

func RegisterCurrentUserRoutes(r chi.Router, service Service, tokenRevoker TokenRevoker) {
	h := NewHandler(service, tokenRevoker)
	r.Get("/users/me", h.GetMe)
	r.Patch("/users/me", h.UpdateMe)
	r.Delete("/users/me", h.DeleteMe)
}

func RegisterAdminRoutes(r chi.Router, service Service, tokenRevoker TokenRevoker) {
	h := NewHandler(service, tokenRevoker)
	r.Post("/users", h.Create)
	r.Get("/users", h.List)
	r.Get("/users/{id}", h.GetByID)
	r.Patch("/users/{id}", h.Update)
	r.Delete("/users/{id}", h.Delete)
}
