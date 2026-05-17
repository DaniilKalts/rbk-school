package httpx

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/6-week/pkg/jwt"
)

type Claims = jwt.Claims

type contextKey string

const claimsKey contextKey = "claims"

func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*Claims)
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
