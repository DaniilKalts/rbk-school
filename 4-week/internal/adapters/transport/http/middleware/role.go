package middleware

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/helpers"
)

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := helpers.ClaimsFromContext(r.Context())
			if !ok {
				response := helpers.NewErrorResponse(http.StatusUnauthorized, "missing authentication claims")
				helpers.JSON(w, http.StatusUnauthorized, response)
				return
			}

			if claims.Role != role {
				response := helpers.NewErrorResponse(http.StatusForbidden, "forbidden")
				helpers.JSON(w, http.StatusForbidden, response)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
