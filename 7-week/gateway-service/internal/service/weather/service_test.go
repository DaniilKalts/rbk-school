package weather_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	geocodingdto "github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/adapter/client/geocoding/dto"
	openmeteodto "github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/adapter/client/openmeteo/dto"
	serviceweather "github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/service/weather"
)

type weatherMocks struct {
	geocoding *mockGeocodingClient
	openMeteo *mockOpenMeteoClient
}

func newService(t *testing.T) (*weatherMocks, *serviceweather.Service) {
	t.Helper()

	m := &weatherMocks{
		geocoding: new(mockGeocodingClient),
		openMeteo: new(mockOpenMeteoClient),
	}
	t.Cleanup(func() {
		m.geocoding.AssertExpectations(t)
		m.openMeteo.AssertExpectations(t)
	})

	svc := serviceweather.NewService(m.geocoding, m.openMeteo)
	return m, svc
}

func TestService_GetByCity(t *testing.T) {
	geoErr := errors.New("geocoding 500")
	weatherErr := errors.New("openmeteo 500")

	tests := []struct {
		name      string
		city      string
		setupMock func(*weatherMocks)
		want      serviceweather.Snapshot
		wantErr   error
	}{
		{
			name: "success",
			city: "Almaty",
			setupMock: func(m *weatherMocks) {
				m.geocoding.On("GetCoordsByCity", mock.Anything, "Almaty").
					Return(geocodingdto.CoordsResponse{Latitude: 43.25, Longitude: 76.95}, nil).Once()
				m.openMeteo.On("GetWeatherByCoords", mock.Anything, 43.25, 76.95).
					Return(openmeteodto.WeatherResponse{
						Current: openmeteodto.CurrentWeather{
							Temperature2M: 20.5, ApparentTemperature: 19.8, WeatherCode: 1,
						},
					}, nil).Once()
			},
			want: serviceweather.Snapshot{
				City: "Almaty", Latitude: 43.25, Longitude: 76.95,
				Temperature: 20.5, ApparentTemperature: 19.8, WeatherCode: 1,
			},
		},
		{
			name: "geocoding error propagates",
			city: "Almaty",
			setupMock: func(m *weatherMocks) {
				m.geocoding.On("GetCoordsByCity", mock.Anything, "Almaty").
					Return(geocodingdto.CoordsResponse{}, geoErr).Once()
			},
			wantErr: geoErr,
		},
		{
			name: "openmeteo error propagates",
			city: "Almaty",
			setupMock: func(m *weatherMocks) {
				m.geocoding.On("GetCoordsByCity", mock.Anything, "Almaty").
					Return(geocodingdto.CoordsResponse{Latitude: 43.25, Longitude: 76.95}, nil).Once()
				m.openMeteo.On("GetWeatherByCoords", mock.Anything, 43.25, 76.95).
					Return(openmeteodto.WeatherResponse{}, weatherErr).Once()
			},
			wantErr: weatherErr,
		},
		{
			name: "geocoding called with exact city name",
			city: "Saint-Petersburg",
			setupMock: func(m *weatherMocks) {
				m.geocoding.On("GetCoordsByCity", mock.Anything, "Saint-Petersburg").
					Return(geocodingdto.CoordsResponse{Latitude: 59.93, Longitude: 30.31}, nil).Once()
				m.openMeteo.On("GetWeatherByCoords", mock.Anything, 59.93, 30.31).
					Return(openmeteodto.WeatherResponse{
						Current: openmeteodto.CurrentWeather{Temperature2M: 5, ApparentTemperature: 3, WeatherCode: 71},
					}, nil).Once()
			},
			want: serviceweather.Snapshot{
				City: "Saint-Petersburg", Latitude: 59.93, Longitude: 30.31,
				Temperature: 5, ApparentTemperature: 3, WeatherCode: 71,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, svc := newService(t)
			if tt.setupMock != nil {
				tt.setupMock(m)
			}

			got, err := svc.GetByCity(context.Background(), tt.city)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, serviceweather.Snapshot{}, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
