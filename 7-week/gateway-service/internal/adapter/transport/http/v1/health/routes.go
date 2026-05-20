package health

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, client APIServiceClient) {
	h := NewHandler(client)

	r.Get("/ready", h.Ready)
}
