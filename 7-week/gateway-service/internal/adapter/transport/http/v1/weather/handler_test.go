package weather_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/adapter/transport/http/v1/weather"
	weathersvc "github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/service/weather"
)

func newRouter(t *testing.T) (*mockService, http.Handler) {
	t.Helper()

	svc := new(mockService)
	t.Cleanup(func() { svc.AssertExpectations(t) })

	r := chi.NewRouter()
	weather.RegisterRoutes(r, svc)
	return svc, r
}

func TestHandler_GetByCity(t *testing.T) {
	upstreamErr := errors.New("upstream is down")

	tests := []struct {
		name       string
		query      string
		setupMock  func(*mockService)
		wantStatus int
		wantBody   string
	}{
		{
			name:  "success",
			query: "?city=Almaty",
			setupMock: func(s *mockService) {
				s.On("GetByCity", mock.Anything, "Almaty").
					Return(weathersvc.Snapshot{
						City: "Almaty", Latitude: 43.25, Longitude: 76.95,
						Temperature: 20, ApparentTemperature: 19, WeatherCode: 1,
					}, nil).Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing city query parameter",
			query:      "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "blank city query parameter",
			query:      "?city=%20%20",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "service error",
			query: "?city=Almaty",
			setupMock: func(s *mockService) {
				s.On("GetByCity", mock.Anything, "Almaty").
					Return(weathersvc.Snapshot{}, upstreamErr).Once()
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:  "city trimmed before reaching service",
			query: "?city=%20Astana%20",
			setupMock: func(s *mockService) {
				s.On("GetByCity", mock.Anything, "Astana").
					Return(weathersvc.Snapshot{
						City: "Astana", Latitude: 51.18, Longitude: 71.44,
						Temperature: 5, ApparentTemperature: 3, WeatherCode: 2,
					}, nil).Once()
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, r := newRouter(t)
			if tt.setupMock != nil {
				tt.setupMock(svc)
			}

			req := httptest.NewRequest(http.MethodGet, "/weather"+tt.query, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var got weather.Response
				require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
				assert.NotEmpty(t, got.City)
				assert.NotZero(t, got.Temperature)
			}
		})
	}
}
