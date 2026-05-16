package transporthttp

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/swagger"
	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/v1"
)

func NewRouter(deps v1.Dependencies, handlerTimeout time.Duration) http.Handler {
	r := chi.NewRouter()
	swagger.RegisterRoutes(r)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.Timeout(handlerTimeout))
		v1.RegisterRoutes(r, deps)
	})

	return r
}
