package geocoding

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client := NewClient(server.Client())
	client.baseURL = server.URL
	return client
}

func TestNewClient_DefaultsHTTPClient(t *testing.T) {
	client := NewClient(nil)

	require.NotNil(t, client.httpClient)
	assert.Equal(t, defaultTimeout, client.httpClient.Timeout)
	assert.Equal(t, defaultBaseURL, client.baseURL)
}

func TestClient_GetCoordsByCity(t *testing.T) {
	tests := []struct {
		name       string
		city       string
		handler    http.HandlerFunc
		wantLat    float64
		wantLon    float64
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success returns first result",
			city: "Almaty",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				query := r.URL.Query()
				assert.Equal(t, "Almaty", query.Get("name"))
				assert.Equal(t, "1", query.Get("count"))
				assert.Equal(t, "en", query.Get("language"))
				assert.Equal(t, "json", query.Get("format"))

				_, _ = w.Write([]byte(`{"results":[{"latitude":43.25,"longitude":76.95},{"latitude":1,"longitude":2}]}`))
			},
			wantLat: 43.25,
			wantLon: 76.95,
		},
		{
			name: "empty results",
			city: "Nowhereville",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(`{"results":[]}`))
			},
			wantErr:    true,
			wantErrMsg: "не найдено результатов",
		},
		{
			name: "non-200 status",
			city: "Almaty",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadGateway)
			},
			wantErr:    true,
			wantErrMsg: "неожиданный статус-код 502",
		},
		{
			name: "invalid json body",
			city: "Almaty",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(`{not json`))
			},
			wantErr:    true,
			wantErrMsg: "ошибка декодирования",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := newTestClient(t, tc.handler)

			got, err := client.GetCoordsByCity(context.Background(), tc.city)
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.NoError(t, err)
			assert.InDelta(t, tc.wantLat, got.Latitude, 0.0001)
			assert.InDelta(t, tc.wantLon, got.Longitude, 0.0001)
		})
	}
}

func TestClient_GetCoordsByCity_ContextCanceled(t *testing.T) {
	client := newTestClient(t, func(_ http.ResponseWriter, _ *http.Request) {})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetCoordsByCity(ctx, "Almaty")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка выполнения запроса")
}

func TestClient_GetCoordsByCity_InvalidBaseURL(t *testing.T) {
	client := NewClient(nil)
	client.baseURL = "://broken"

	_, err := client.GetCoordsByCity(context.Background(), "Almaty")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка разбора базового URL")
}
