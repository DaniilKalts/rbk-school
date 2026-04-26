package weather

import "net/http"

func RegisterRoutes(mux *http.ServeMux, service Service) {
	h := NewHandler(service)

	mux.HandleFunc("GET /api/v1/users/{id}/weather", h.GetByUserID)
}
