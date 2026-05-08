package middleware

import (
	"context"
	"net/http"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapters/transport/http/helpers"
)

type JWTManager interface {
	Validate(tokenString string) (*helpers.Claims, error)
	IsRevoked(ctx context.Context, tokenString string) (bool, error)
}

func Auth(jwtManager JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := helpers.BearerTokenFromRequest(r)
			if !ok {
				response := helpers.NewErrorResponse(http.StatusUnauthorized, "отсутствует или некорректный заголовок Authorization")
				helpers.JSON(w, http.StatusUnauthorized, response)
				return
			}

			claims, err := jwtManager.Validate(token)
			if err != nil {
				response := helpers.NewErrorResponse(http.StatusUnauthorized, "некорректный или просроченный токен")
				helpers.JSON(w, http.StatusUnauthorized, response)
				return
			}

			revoked, err := jwtManager.IsRevoked(r.Context(), token)
			if err != nil {
				response := helpers.NewErrorResponse(http.StatusServiceUnavailable, "сервис временно недоступен")
				helpers.JSON(w, http.StatusServiceUnavailable, response)
				return
			}
			if revoked {
				response := helpers.NewErrorResponse(http.StatusUnauthorized, "некорректный или просроченный токен")
				helpers.JSON(w, http.StatusUnauthorized, response)
				return
			}

			ctx := helpers.WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
