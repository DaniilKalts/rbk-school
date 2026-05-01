package docs

import (
	"net/http"
	"os"
	"path/filepath"
)

func RegisterRoutes(mux *http.ServeMux) {
	swaggerDir := firstExistingPath(filepath.Join("web", "swagger"), filepath.Join("4-week", "web", "swagger"))
	openAPIPath := firstExistingPath(filepath.Join("api", "v1", "openapi.yaml"), filepath.Join("4-week", "api", "v1", "openapi.yaml"))

	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	mux.HandleFunc("GET /swagger/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		http.ServeFile(w, r, filepath.Join(swaggerDir, "index.html"))
	})

	mux.Handle("GET /swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir(swaggerDir))))

	mux.HandleFunc("GET /api/v1/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
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
