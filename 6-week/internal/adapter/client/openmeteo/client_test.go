package openmeteo

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

func TestClient_GetWeatherByCoords(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantTemp   float64
		wantApp    float64
		wantCode   int
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				query := r.URL.Query()
				assert.Equal(t, "43.25", query.Get("latitude"))
				assert.Equal(t, "76.95", query.Get("longitude"))
				assert.Equal(t, "temperature_2m,apparent_temperature,weather_code", query.Get("current"))

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"current":{"temperature_2m":21.5,"apparent_temperature":20.1,"weather_code":3}}`))
			},
			wantTemp: 21.5,
			wantApp:  20.1,
			wantCode: 3,
		},
		{
			name: "non-200 status",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr:    true,
			wantErrMsg: "неожиданный статус-код 500",
		},
		{
			name: "invalid json body",
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

			got, err := client.GetWeatherByCoords(context.Background(), 43.25, 76.95)
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrMsg)
				return
			}

			require.NoError(t, err)
			assert.InDelta(t, tc.wantTemp, got.Current.Temperature2M, 0.0001)
			assert.InDelta(t, tc.wantApp, got.Current.ApparentTemperature, 0.0001)
			assert.Equal(t, tc.wantCode, got.Current.WeatherCode)
		})
	}
}

func TestClient_GetWeatherByCoords_ContextCanceled(t *testing.T) {
	client := newTestClient(t, func(_ http.ResponseWriter, _ *http.Request) {})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetWeatherByCoords(ctx, 0, 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка выполнения запроса")
}

func TestClient_GetWeatherByCoords_InvalidBaseURL(t *testing.T) {
	client := NewClient(nil)
	client.baseURL = "://broken"

	_, err := client.GetWeatherByCoords(context.Background(), 0, 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка разбора базового URL")
}
