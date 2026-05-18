//go:build integration

package blacklist_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	blacklistcache "github.com/DaniilKalts/rbk-school/6-week/internal/cache/blacklist"
)

var client *redis.Client

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := tcredis.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "redis start: %v\n", err)
		os.Exit(1)
	}

	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		_ = container.Terminate(ctx)
		fmt.Fprintf(os.Stderr, "endpoint: %v\n", err)
		os.Exit(1)
	}

	client = redis.NewClient(&redis.Options{Addr: endpoint})
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		_ = container.Terminate(ctx)
		fmt.Fprintf(os.Stderr, "ping: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	_ = client.Close()
	_ = container.Terminate(ctx)
	os.Exit(code)
}

func reset(t *testing.T) {
	t.Helper()
	require.NoError(t, client.FlushDB(context.Background()).Err())
}

func newCache(t *testing.T) (context.Context, *blacklistcache.Cache) {
	t.Helper()
	reset(t)
	return context.Background(), blacklistcache.NewCache(client)
}

func TestBlacklistCache_IsRevoked_UnknownToken(t *testing.T) {
	ctx, cache := newCache(t)

	revoked, err := cache.IsRevoked(ctx, "never-seen-before")

	require.NoError(t, err)
	assert.False(t, revoked)
}

func TestBlacklistCache_Revoke_Then_IsRevoked(t *testing.T) {
	ctx, cache := newCache(t)

	require.NoError(t, cache.Revoke(ctx, "tok-abc", time.Now().Add(time.Minute)))

	revoked, err := cache.IsRevoked(ctx, "tok-abc")
	require.NoError(t, err)
	assert.True(t, revoked)
}

func TestBlacklistCache_Revoke_PastExpiry_IsNoop(t *testing.T) {
	ctx, cache := newCache(t)

	err := cache.Revoke(ctx, "tok-stale", time.Now().Add(-time.Minute))
	require.NoError(t, err)

	revoked, err := cache.IsRevoked(ctx, "tok-stale")
	require.NoError(t, err)
	assert.False(t, revoked)
}

func TestBlacklistCache_TTL_Expires(t *testing.T) {
	ctx, cache := newCache(t)
	require.NoError(t, cache.Revoke(ctx, "tok-short", time.Now().Add(100*time.Millisecond)))

	revoked, err := cache.IsRevoked(ctx, "tok-short")
	require.NoError(t, err)
	require.True(t, revoked)

	time.Sleep(200 * time.Millisecond)

	revoked, err = cache.IsRevoked(ctx, "tok-short")
	require.NoError(t, err)
	assert.False(t, revoked)
}

func TestBlacklistCache_TokensAreIsolated(t *testing.T) {
	ctx, cache := newCache(t)
	require.NoError(t, cache.Revoke(ctx, "tok-a", time.Now().Add(time.Minute)))

	a, err := cache.IsRevoked(ctx, "tok-a")
	require.NoError(t, err)
	assert.True(t, a)

	b, err := cache.IsRevoked(ctx, "tok-b")
	require.NoError(t, err)
	assert.False(t, b)
}

func TestBlacklistCache_Revoke_ContextCancelledWrapsError(t *testing.T) {
	ctx, cache := newCache(t)
	cancelled, cancel := context.WithCancel(ctx)
	cancel()

	err := cache.Revoke(cancelled, "tok-x", time.Now().Add(time.Minute))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "добавление отозванного токена")
}

func TestBlacklistCache_IsRevoked_ContextCancelledWrapsError(t *testing.T) {
	ctx, cache := newCache(t)
	cancelled, cancel := context.WithCancel(ctx)
	cancel()

	_, err := cache.IsRevoked(cancelled, "tok-x")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "проверка отозванного токена")
}
