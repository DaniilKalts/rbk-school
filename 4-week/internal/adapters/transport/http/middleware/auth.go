package middleware

import (
	"net/http"

	"github.com/DaniilKalts/rbk-school/4-week/internal/adapters/transport/http/helpers"
)

type JWTManager interface {
	Validate(tokenString string) (*helpers.Claims, error)
}

func Auth(jwtManager JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := helpers.BearerTokenFromRequest(r)
			if !ok {
				response := helpers.NewErrorResponse(http.StatusUnauthorized, "missing or malformed authorization header")
				helpers.JSON(w, http.StatusUnauthorized, response)
				return
			}

			claims, err := jwtManager.Validate(token)
			if err != nil {
				response := helpers.NewErrorResponse(http.StatusUnauthorized, "invalid or expired token")
				helpers.JSON(w, http.StatusUnauthorized, response)
				return
			}

			ctx := helpers.WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
