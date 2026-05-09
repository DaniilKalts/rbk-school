package middleware

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapter/transport/http/httpx"
)

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := httpx.ClaimsFromContext(r.Context())
			if !ok {
				httpx.WriteError(w, http.StatusUnauthorized, "отсутствуют claims аутентификации")
				return
			}

			if claims.Role != role {
				httpx.WriteError(w, http.StatusForbidden, "доступ запрещен")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
