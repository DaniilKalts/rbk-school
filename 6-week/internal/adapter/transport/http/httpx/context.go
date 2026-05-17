package httpx

import (
	"context"
	"net/http"

	"github.com/DaniilKalts/rbk-school/6-week/pkg/jwt"
	"github.com/google/uuid"
)

type Claims = jwt.Claims

type (
	claimsKey    struct{}
	requestIDKey struct{}
)

func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey{}).(*Claims)
	return claims, ok
}

func CurrentUserID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	claims, ok := ClaimsFromContext(r.Context())
	if !ok || claims.UserID == uuid.Nil {
		WriteError(w, http.StatusUnauthorized, "отсутствуют claims аутентификации")
		return uuid.Nil, false
	}
	return claims.UserID, true
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}
