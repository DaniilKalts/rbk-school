package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/DaniilKalts/rbk-school/7-week/internal/adapter/transport/http/middleware"
	"github.com/DaniilKalts/rbk-school/7-week/internal/adapter/transport/http/swagger"
	"github.com/DaniilKalts/rbk-school/7-week/internal/adapter/transport/http/v1"
)

func NewRouter(logger *zap.Logger, deps v1.Dependencies, handlerTimeout time.Duration) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(logger))

	swagger.RegisterRoutes(r)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Group(func(r chi.Router) {
		r.Use(chimiddleware.Timeout(handlerTimeout))
		v1.RegisterRoutes(r, deps)
	})

	return r
}
