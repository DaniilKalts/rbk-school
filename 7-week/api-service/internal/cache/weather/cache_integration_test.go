//go:build integration

package weather_test

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

	weathercache "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/cache/weather"
	domainweather "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/weather"
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

func newCache(t *testing.T, ttl time.Duration) (context.Context, *weathercache.Cache) {
	t.Helper()
	reset(t)
	return context.Background(), weathercache.NewCache(client, ttl)
}

func sampleWeather() domainweather.Weather {
	return domainweather.Weather{
		City:        "Almaty",
		Temperature: 20.5,
		FeelsLike:   19.0,
		Description: "clear sky",
		RequestedAt: time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC),
	}
}

func TestWeatherCache_Get_Miss(t *testing.T) {
	ctx, cache := newCache(t, time.Minute)

	got, ok, err := cache.Get(ctx, "Almaty")

	require.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, domainweather.Weather{}, got)
}

func TestWeatherCache_Set_Then_Get_Roundtrip(t *testing.T) {
	ctx, cache := newCache(t, time.Minute)
	w := sampleWeather()

	require.NoError(t, cache.Set(ctx, "Almaty", w))

	got, ok, err := cache.Get(ctx, "Almaty")

	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, w.City, got.City)
	assert.InDelta(t, w.Temperature, got.Temperature, 0.001)
	assert.Equal(t, w.Description, got.Description)
}

func TestWeatherCache_Set_StripsRequestedAt(t *testing.T) {
	ctx, cache := newCache(t, time.Minute)
	w := sampleWeather()
	require.False(t, w.RequestedAt.IsZero())

	require.NoError(t, cache.Set(ctx, "Almaty", w))

	got, ok, err := cache.Get(ctx, "Almaty")
	require.NoError(t, err)
	require.True(t, ok)
	assert.True(t, got.RequestedAt.IsZero())
}

func TestWeatherCache_KeyIsCaseAndWhitespaceInsensitive(t *testing.T) {
	ctx, cache := newCache(t, time.Minute)
	require.NoError(t, cache.Set(ctx, "  Almaty  ", sampleWeather()))

	got, ok, err := cache.Get(ctx, "almaty")

	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "Almaty", got.City)
}

func TestWeatherCache_TTL_Expires(t *testing.T) {
	ctx, cache := newCache(t, 100*time.Millisecond)
	require.NoError(t, cache.Set(ctx, "Almaty", sampleWeather()))

	_, ok, err := cache.Get(ctx, "Almaty")
	require.NoError(t, err)
	require.True(t, ok)

	time.Sleep(200 * time.Millisecond)

	_, ok, err = cache.Get(ctx, "Almaty")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestWeatherCache_Get_MalformedPayload(t *testing.T) {
	ctx, cache := newCache(t, time.Minute)
	require.NoError(t, client.Set(ctx, "weather:almaty", "not-json", time.Minute).Err())

	_, _, err := cache.Get(ctx, "Almaty")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "декодирование кеша погоды")
}

func TestWeatherCache_Get_ContextCancelledWrapsError(t *testing.T) {
	ctx, cache := newCache(t, time.Minute)
	cancelled, cancel := context.WithCancel(ctx)
	cancel()

	_, _, err := cache.Get(cancelled, "Almaty")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "получение кеша погоды")
}
