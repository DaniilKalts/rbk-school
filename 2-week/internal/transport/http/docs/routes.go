package docs

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	r.Handle("/swagger/*", http.StripPrefix("/swagger/", http.FileServer(http.Dir(filepath.Join("web", "swagger")))))

	r.Get("/api/v1/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		http.ServeFile(w, r, filepath.Join("api", "v1", "openapi.yaml"))
	})
}
