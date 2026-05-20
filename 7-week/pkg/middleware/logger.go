package middleware

import (
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/DaniilKalts/rbk-school/7-week/pkg/httpx"
	"github.com/DaniilKalts/rbk-school/7-week/pkg/logger"
)

func Logger(base *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapped := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			scoped := base.With(zap.String("request_id", httpx.RequestIDFromContext(r.Context())))
			r = r.WithContext(logger.WithContext(r.Context(), scoped))

			defer func() {
				status := wrapped.Status()
				if status == 0 {
					status = http.StatusOK
				}

				fields := []zap.Field{
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Int("status", status),
					zap.Duration("duration", time.Since(start)),
					zap.String("remote_addr", r.RemoteAddr),
					zap.String("user_agent", r.UserAgent()),
				}

				switch {
				case status >= 500:
					scoped.Error("HTTP-запрос", fields...)
				case status >= 400:
					scoped.Warn("HTTP-запрос", fields...)
				default:
					scoped.Info("HTTP-запрос", fields...)
				}
			}()

			next.ServeHTTP(wrapped, r)
		})
	}
}
