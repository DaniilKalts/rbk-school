package docs

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	swaggerDir := firstExistingPath(filepath.Join("web", "swagger"), filepath.Join("5-week", "web", "swagger"))
	openAPIPath := firstExistingPath(filepath.Join("api", "v1", "openapi.yaml"), filepath.Join("5-week", "api", "v1", "openapi.yaml"))

	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	r.Get("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		http.ServeFile(w, r, filepath.Join(swaggerDir, "index.html"))
	})

	r.Handle("/swagger/*", http.StripPrefix("/swagger/", http.FileServer(http.Dir(swaggerDir))))

	r.Get("/api/v1/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.Header().Set("Cache-Control", "no-store")
		http.ServeFile(w, r, openAPIPath)
	})
}

func firstExistingPath(paths ...string) string {
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return paths[0]
}
