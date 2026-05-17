package middleware

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/httpx"
	"github.com/DaniilKalts/rbk-school/6-week/pkg/logger"
)

type JWTManager interface {
	Validate(tokenString string) (*httpx.Claims, error)
	IsRevoked(ctx context.Context, tokenString string) (bool, error)
}

func Auth(jwtManager JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := httpx.BearerTokenFromRequest(r)
			if !ok {
				httpx.WriteError(w, http.StatusUnauthorized, "отсутствует или некорректный заголовок Authorization")
				return
			}

			claims, err := jwtManager.Validate(token)
			if err != nil {
				httpx.WriteError(w, http.StatusUnauthorized, "некорректный или просроченный токен")
				return
			}

			revoked, err := jwtManager.IsRevoked(r.Context(), token)
			if err != nil {
				logger.FromContext(r.Context()).Warn("проверка отозванности токена", zap.Error(err))
				httpx.WriteError(w, http.StatusServiceUnavailable, "сервис временно недоступен")
				return
			}
			if revoked {
				httpx.WriteError(w, http.StatusUnauthorized, "некорректный или просроченный токен")
				return
			}

			ctx := httpx.WithClaims(r.Context(), claims)
			ctx = logger.WithContext(ctx, logger.FromContext(ctx).With(zap.String("user_id", claims.UserID.String())))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
