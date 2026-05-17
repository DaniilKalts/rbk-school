package middleware

import (
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapped := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			defer func() {
				status := wrapped.Status()
				if status == 0 {
					status = http.StatusOK
				}

				fields := []zap.Field{
					zap.String("request_id", RequestIDFromContext(r.Context())),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Int("status", status),
					zap.Duration("duration", time.Since(start)),
					zap.String("remote_addr", r.RemoteAddr),
					zap.String("user_agent", r.UserAgent()),
				}

				switch {
				case status >= 500:
					logger.Error("HTTP-запрос", fields...)
				case status >= 400:
					logger.Warn("HTTP-запрос", fields...)
				default:
					logger.Info("HTTP-запрос", fields...)
				}
			}()

			next.ServeHTTP(wrapped, r)
		})
	}
}
