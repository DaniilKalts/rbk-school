package helpers

import (
	"context"

	"github.com/DaniilKalts/rbk-school/5-week/internal/utils"
)

type Claims = utils.Claims

type contextKey string

const claimsKey contextKey = "claims"

func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	return claims, ok
}
