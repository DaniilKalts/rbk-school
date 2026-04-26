package user

import "net/http"

func RegisterRoutes(mux *http.ServeMux, service Service) {
	h := NewHandler(service)

	mux.HandleFunc("POST /api/v1/users", h.Create)
	mux.HandleFunc("GET /api/v1/users", h.List)
	mux.HandleFunc("GET /api/v1/users/{id}", h.GetByID)
	mux.HandleFunc("PUT /api/v1/users/{id}", h.Update)
	mux.HandleFunc("DELETE /api/v1/users/{id}", h.Delete)
}
