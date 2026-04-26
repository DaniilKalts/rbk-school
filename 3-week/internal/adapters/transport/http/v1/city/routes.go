package city

import "net/http"

func RegisterRoutes(mux *http.ServeMux, service Service) {
	h := NewHandler(service)

	mux.HandleFunc("POST /api/v1/users/{id}/cities", h.Create)
	mux.HandleFunc("GET /api/v1/users/{id}/cities", h.List)
	mux.HandleFunc("DELETE /api/v1/users/{id}/cities/{city_id}", h.Delete)
}
