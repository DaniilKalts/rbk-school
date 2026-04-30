package auth

import "net/http"

func RegisterRoutes(mux *http.ServeMux, service Service) {
	h := NewHandler(service)

	mux.HandleFunc("POST /api/v1/auth/register", h.Register)
	mux.HandleFunc("POST /api/v1/auth/login", h.Login)
}
