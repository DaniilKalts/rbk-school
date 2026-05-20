package middleware

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/adapter/transport/http/httpx"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewString()
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r.WithContext(httpx.WithRequestID(r.Context(), id)))
	})
}
