package transporthttp

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/docs"
	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/v1"
)

func NewRouter(deps v1.Dependencies) http.Handler {
	r := chi.NewRouter()
	docs.RegisterRoutes(r)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	v1.RegisterRoutes(r, deps)

	return r
}
