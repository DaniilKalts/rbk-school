package swagger

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

const (
	openAPIURLPrefix = "/api/v1/"
	openAPIRootFile  = "openapi.yaml"
)

func RegisterRoutes(r chi.Router) {
	swaggerDir := firstExistingPath(filepath.Join("web", "swagger"), filepath.Join("6-week", "web", "swagger"))
	openAPIDir := firstExistingPath(filepath.Join("api", "v1"), filepath.Join("6-week", "api", "v1"))

	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	r.Get("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		http.ServeFile(w, r, filepath.Join(swaggerDir, "index.html"))
	})

	r.Handle("/swagger/*", http.StripPrefix("/swagger/", http.FileServer(http.Dir(swaggerDir))))

	openAPIFS := http.StripPrefix(openAPIURLPrefix, yamlOnly(http.FileServer(http.Dir(openAPIDir))))
	r.Handle(openAPIURLPrefix+openAPIRootFile, openAPIFS)
	for _, subdir := range openAPISubdirs(openAPIDir) {
		r.Handle(openAPIURLPrefix+subdir+"/*", openAPIFS)
	}
}

func openAPISubdirs(root string) []string {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}

	subdirs := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			subdirs = append(subdirs, e.Name())
		}
	}

	return subdirs
}

func yamlOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ext := strings.ToLower(filepath.Ext(r.URL.Path)); ext != ".yaml" && ext != ".yml" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/yaml")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
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
