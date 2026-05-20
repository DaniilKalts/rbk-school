package httpx_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/adapter/transport/http/httpx"
)

func TestClaimsRoundtrip(t *testing.T) {
	claims := &httpx.Claims{UserID: uuid.New(), Email: "u@example.com", Role: "user"}

	ctx := httpx.WithClaims(context.Background(), claims)
	got, ok := httpx.ClaimsFromContext(ctx)

	require.True(t, ok)
	assert.Equal(t, claims, got)
}

func TestClaimsFromContext_Missing(t *testing.T) {
	got, ok := httpx.ClaimsFromContext(context.Background())

	assert.False(t, ok)
	assert.Nil(t, got)
}

func TestRequestIDRoundtrip(t *testing.T) {
	ctx := httpx.WithRequestID(context.Background(), "req-123")

	assert.Equal(t, "req-123", httpx.RequestIDFromContext(ctx))
}

func TestRequestIDFromContext_Missing(t *testing.T) {
	assert.Equal(t, "", httpx.RequestIDFromContext(context.Background()))
}

func TestCurrentUserID_Success(t *testing.T) {
	id := uuid.New()
	claims := &httpx.Claims{UserID: id}
	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(httpx.WithClaims(context.Background(), claims))
	w := httptest.NewRecorder()

	got, ok := httpx.CurrentUserID(w, req)

	require.True(t, ok)
	assert.Equal(t, id, got)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCurrentUserID_NoClaims(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	_, ok := httpx.CurrentUserID(w, req)

	require.False(t, ok)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body httpx.ErrorResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
	assert.Equal(t, http.StatusUnauthorized, body.Code)
}

func TestCurrentUserID_NilUUID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(
		httpx.WithClaims(context.Background(), &httpx.Claims{UserID: uuid.Nil}),
	)
	w := httptest.NewRecorder()

	_, ok := httpx.CurrentUserID(w, req)

	require.False(t, ok)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
