package weather_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	geocodingdto "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/adapter/client/geocoding/dto"
	openmeteodto "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/adapter/client/openmeteo/dto"
	domaincity "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/city"
	domainhistory "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/history"
	domainuser "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/user"
	domainweather "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/weather"
	serviceweather "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/service/weather"
)

type weatherMocks struct {
	user      *mockUserRepository
	city      *mockCityRepository
	history   *mockHistoryRepository
	geocoding *mockGeocodingClient
	weather   *mockWeatherClient
	cache     *mockWeatherCache
}

func newService(t *testing.T) (*weatherMocks, *serviceweather.Service) {
	t.Helper()

	m := &weatherMocks{
		user:      new(mockUserRepository),
		city:      new(mockCityRepository),
		history:   new(mockHistoryRepository),
		geocoding: new(mockGeocodingClient),
		weather:   new(mockWeatherClient),
		cache:     new(mockWeatherCache),
	}
	t.Cleanup(func() {
		m.user.AssertExpectations(t)
		m.city.AssertExpectations(t)
		m.history.AssertExpectations(t)
		m.geocoding.AssertExpectations(t)
		m.weather.AssertExpectations(t)
		m.cache.AssertExpectations(t)
	})

	svc := serviceweather.NewService(m.user, m.city, m.history, m.geocoding, m.weather, m.cache)
	return m, svc
}

func TestService_GetByUserID(t *testing.T) {
	userID := uuid.New()
	historyAt := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)
	cached := domainweather.Weather{
		City: "Almaty", Temperature: 15, FeelsLike: 14, Description: "rain",
	}
	geoErr := errors.New("geocoding 500")
	historyErr := errors.New("history insert failed")
	cityErr := errors.New("db down")
	weatherErr := errors.New("openmeteo 500")
	cacheErr := errors.New("redis down")

	tests := []struct {
		name        string
		userID      uuid.UUID
		setupMock   func(*weatherMocks)
		assertExtra func(t *testing.T, m *weatherMocks, got []domainweather.Weather)
		wantErr     error
	}{
		{
			name:   "success cache miss fetches and stores",
			userID: userID,
			setupMock: func(m *weatherMocks) {
				m.user.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				m.city.On("ListByUserID", mock.Anything, userID).
					Return([]domaincity.City{{UserID: userID, Name: "Almaty"}}, nil).Once()
				m.cache.On("Get", mock.Anything, "Almaty").
					Return(domainweather.Weather{}, false, nil).Once()
				m.geocoding.On("GetCoordsByCity", mock.Anything, "Almaty").
					Return(geocodingdto.CoordsResponse{Latitude: 43.25, Longitude: 76.95}, nil).Once()
				m.weather.On("GetWeatherByCoords", mock.Anything, 43.25, 76.95).
					Return(openmeteodto.WeatherResponse{
						Current: openmeteodto.CurrentWeather{
							Temperature2M: 20, ApparentTemperature: 19, WeatherCode: 0,
						},
					}, nil).Once()
				m.cache.On("Set", mock.Anything, "Almaty", mock.AnythingOfType("weather.Weather")).
					Return(nil).Once()
				m.history.On("CreateHistory", mock.Anything, mock.AnythingOfType("history.History")).
					Return(&domainhistory.History{RequestedAt: historyAt}, nil).Once()
			},
			assertExtra: func(t *testing.T, _ *weatherMocks, got []domainweather.Weather) {
				require.Len(t, got, 1)
				assert.Equal(t, "Almaty", got[0].City)
				assert.Equal(t, 20.0, got[0].Temperature)
				assert.Equal(t, 19.0, got[0].FeelsLike)
				assert.Equal(t, "clear sky", got[0].Description)
				assert.Equal(t, historyAt, got[0].RequestedAt)
			},
		},
		{
			name:   "success cache hit skips external calls",
			userID: userID,
			setupMock: func(m *weatherMocks) {
				m.user.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				m.city.On("ListByUserID", mock.Anything, userID).
					Return([]domaincity.City{{UserID: userID, Name: "Almaty"}}, nil).Once()
				m.cache.On("Get", mock.Anything, "Almaty").Return(cached, true, nil).Once()
				m.history.On("CreateHistory", mock.Anything, mock.AnythingOfType("history.History")).
					Return(&domainhistory.History{RequestedAt: historyAt}, nil).Once()
			},
			assertExtra: func(t *testing.T, m *weatherMocks, got []domainweather.Weather) {
				require.Len(t, got, 1)
				assert.Equal(t, "Almaty", got[0].City)
				assert.Equal(t, 15.0, got[0].Temperature)
				assert.Equal(t, historyAt, got[0].RequestedAt)
				m.geocoding.AssertNotCalled(t, "GetCoordsByCity")
				m.weather.AssertNotCalled(t, "GetWeatherByCoords")
				m.cache.AssertNotCalled(t, "Set")
			},
		},
		{
			name:    "nil user id",
			userID:  uuid.Nil,
			wantErr: domainuser.ErrInvalidID,
		},
		{
			name:   "user not found",
			userID: userID,
			setupMock: func(m *weatherMocks) {
				m.user.On("GetByID", mock.Anything, userID).Return(nil, domainuser.ErrNotFound).Once()
			},
			wantErr: domainuser.ErrNotFound,
		},
		{
			name:   "city list error",
			userID: userID,
			setupMock: func(m *weatherMocks) {
				m.user.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				m.city.On("ListByUserID", mock.Anything, userID).Return(nil, cityErr).Once()
			},
			wantErr: cityErr,
		},
		{
			name:   "user has no cities returns empty result",
			userID: userID,
			setupMock: func(m *weatherMocks) {
				m.user.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				m.city.On("ListByUserID", mock.Anything, userID).
					Return([]domaincity.City{}, nil).Once()
			},
			assertExtra: func(t *testing.T, m *weatherMocks, got []domainweather.Weather) {
				assert.Empty(t, got)
				m.cache.AssertNotCalled(t, "Get")
				m.geocoding.AssertNotCalled(t, "GetCoordsByCity")
				m.weather.AssertNotCalled(t, "GetWeatherByCoords")
				m.history.AssertNotCalled(t, "CreateHistory")
			},
		},
		{
			name:   "geocoding error",
			userID: userID,
			setupMock: func(m *weatherMocks) {
				m.user.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				m.city.On("ListByUserID", mock.Anything, userID).
					Return([]domaincity.City{{UserID: userID, Name: "Almaty"}}, nil).Once()
				m.cache.On("Get", mock.Anything, "Almaty").
					Return(domainweather.Weather{}, false, nil).Once()
				m.geocoding.On("GetCoordsByCity", mock.Anything, "Almaty").
					Return(geocodingdto.CoordsResponse{}, geoErr).Once()
			},
			wantErr: geoErr,
		},
		{
			name:   "history error",
			userID: userID,
			setupMock: func(m *weatherMocks) {
				m.user.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				m.city.On("ListByUserID", mock.Anything, userID).
					Return([]domaincity.City{{UserID: userID, Name: "Almaty"}}, nil).Once()
				m.cache.On("Get", mock.Anything, "Almaty").Return(cached, true, nil).Once()
				m.history.On("CreateHistory", mock.Anything, mock.AnythingOfType("history.History")).
					Return(nil, historyErr).Once()
			},
			wantErr: historyErr,
		},
		{
			name:   "weather client error",
			userID: userID,
			setupMock: func(m *weatherMocks) {
				m.user.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				m.city.On("ListByUserID", mock.Anything, userID).
					Return([]domaincity.City{{UserID: userID, Name: "Almaty"}}, nil).Once()
				m.cache.On("Get", mock.Anything, "Almaty").
					Return(domainweather.Weather{}, false, nil).Once()
				m.geocoding.On("GetCoordsByCity", mock.Anything, "Almaty").
					Return(geocodingdto.CoordsResponse{Latitude: 43.25, Longitude: 76.95}, nil).Once()
				m.weather.On("GetWeatherByCoords", mock.Anything, 43.25, 76.95).
					Return(openmeteodto.WeatherResponse{}, weatherErr).Once()
			},
			wantErr: weatherErr,
		},
		{
			name:   "invalid temperature from weather API",
			userID: userID,
			setupMock: func(m *weatherMocks) {
				m.user.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				m.city.On("ListByUserID", mock.Anything, userID).
					Return([]domaincity.City{{UserID: userID, Name: "Almaty"}}, nil).Once()
				m.cache.On("Get", mock.Anything, "Almaty").
					Return(domainweather.Weather{}, false, nil).Once()
				m.geocoding.On("GetCoordsByCity", mock.Anything, "Almaty").
					Return(geocodingdto.CoordsResponse{Latitude: 43.25, Longitude: 76.95}, nil).Once()
				m.weather.On("GetWeatherByCoords", mock.Anything, 43.25, 76.95).
					Return(openmeteodto.WeatherResponse{
						Current: openmeteodto.CurrentWeather{
							Temperature2M: 999, ApparentTemperature: 19, WeatherCode: 0,
						},
					}, nil).Once()
			},
			wantErr: domainweather.ErrInvalidTemperature,
		},
		{
			name:   "cache get error falls through to API",
			userID: userID,
			setupMock: func(m *weatherMocks) {
				m.user.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				m.city.On("ListByUserID", mock.Anything, userID).
					Return([]domaincity.City{{UserID: userID, Name: "Almaty"}}, nil).Once()
				m.cache.On("Get", mock.Anything, "Almaty").
					Return(domainweather.Weather{}, false, cacheErr).Once()
				m.geocoding.On("GetCoordsByCity", mock.Anything, "Almaty").
					Return(geocodingdto.CoordsResponse{Latitude: 43.25, Longitude: 76.95}, nil).Once()
				m.weather.On("GetWeatherByCoords", mock.Anything, 43.25, 76.95).
					Return(openmeteodto.WeatherResponse{
						Current: openmeteodto.CurrentWeather{Temperature2M: 20, ApparentTemperature: 19, WeatherCode: 0},
					}, nil).Once()
				m.cache.On("Set", mock.Anything, "Almaty", mock.AnythingOfType("weather.Weather")).
					Return(nil).Once()
				m.history.On("CreateHistory", mock.Anything, mock.AnythingOfType("history.History")).
					Return(&domainhistory.History{RequestedAt: historyAt}, nil).Once()
			},
			assertExtra: func(t *testing.T, _ *weatherMocks, got []domainweather.Weather) {
				require.Len(t, got, 1)
				assert.Equal(t, "Almaty", got[0].City)
			},
		},
		{
			name:   "cache set error does not fail request",
			userID: userID,
			setupMock: func(m *weatherMocks) {
				m.user.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				m.city.On("ListByUserID", mock.Anything, userID).
					Return([]domaincity.City{{UserID: userID, Name: "Almaty"}}, nil).Once()
				m.cache.On("Get", mock.Anything, "Almaty").
					Return(domainweather.Weather{}, false, nil).Once()
				m.geocoding.On("GetCoordsByCity", mock.Anything, "Almaty").
					Return(geocodingdto.CoordsResponse{Latitude: 43.25, Longitude: 76.95}, nil).Once()
				m.weather.On("GetWeatherByCoords", mock.Anything, 43.25, 76.95).
					Return(openmeteodto.WeatherResponse{
						Current: openmeteodto.CurrentWeather{Temperature2M: 20, ApparentTemperature: 19, WeatherCode: 0},
					}, nil).Once()
				m.cache.On("Set", mock.Anything, "Almaty", mock.AnythingOfType("weather.Weather")).
					Return(cacheErr).Once()
				m.history.On("CreateHistory", mock.Anything, mock.AnythingOfType("history.History")).
					Return(&domainhistory.History{RequestedAt: historyAt}, nil).Once()
			},
			assertExtra: func(t *testing.T, _ *weatherMocks, got []domainweather.Weather) {
				require.Len(t, got, 1)
			},
		},
		{
			name:   "multiple cities in parallel",
			userID: userID,
			setupMock: func(m *weatherMocks) {
				m.user.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				m.city.On("ListByUserID", mock.Anything, userID).
					Return([]domaincity.City{
						{UserID: userID, Name: "Almaty"},
						{UserID: userID, Name: "Astana"},
					}, nil).Once()
				m.cache.On("Get", mock.Anything, "Almaty").
					Return(domainweather.Weather{City: "Almaty", Temperature: 15, FeelsLike: 14, Description: "rain"}, true, nil).Once()
				m.cache.On("Get", mock.Anything, "Astana").
					Return(domainweather.Weather{City: "Astana", Temperature: 5, FeelsLike: 3, Description: "snow"}, true, nil).Once()
				m.history.On("CreateHistory", mock.Anything, mock.AnythingOfType("history.History")).
					Return(&domainhistory.History{RequestedAt: historyAt}, nil).Twice()
			},
			assertExtra: func(t *testing.T, _ *weatherMocks, got []domainweather.Weather) {
				require.Len(t, got, 2)
				assert.Equal(t, "Almaty", got[0].City)
				assert.Equal(t, "Astana", got[1].City)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, svc := newService(t)
			if tt.setupMock != nil {
				tt.setupMock(m)
			}

			got, err := svc.GetByUserID(context.Background(), tt.userID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}
			require.NoError(t, err)
			if tt.assertExtra != nil {
				tt.assertExtra(t, m, got)
			}
		})
	}
}

func TestService_GetHistory(t *testing.T) {
	userID := uuid.New()
	expected := []domainhistory.History{{City: "Almaty", Temperature: 10}}
	repoErr := errors.New("db down")

	tests := []struct {
		name      string
		userID    uuid.UUID
		cityName  string
		limit     int
		offset    int
		setupMock func(*weatherMocks)
		want      []domainhistory.History
		wantErr   error
	}{
		{
			name:     "success normalizes city",
			userID:   userID,
			cityName: " almaty ",
			limit:    10,
			offset:   0,
			setupMock: func(m *weatherMocks) {
				m.user.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				m.history.On("ListHistory", mock.Anything, userID, "Almaty", 10, 0).
					Return(expected, nil).Once()
			},
			want: expected,
		},
		{
			name:    "nil user id",
			userID:  uuid.Nil,
			limit:   10,
			wantErr: domainuser.ErrInvalidID,
		},
		{
			name:    "negative limit",
			userID:  userID,
			limit:   -1,
			wantErr: domainweather.ErrInvalidLimit,
		},
		{
			name:    "negative offset",
			userID:  userID,
			limit:   10,
			offset:  -1,
			wantErr: domainweather.ErrInvalidOffset,
		},
		{
			name:     "user not found",
			userID:   userID,
			cityName: "Almaty",
			limit:    10,
			setupMock: func(m *weatherMocks) {
				m.user.On("GetByID", mock.Anything, userID).Return(nil, domainuser.ErrNotFound).Once()
			},
			wantErr: domainuser.ErrNotFound,
		},
		{
			name:     "repo error propagates",
			userID:   userID,
			cityName: "Almaty",
			limit:    10,
			setupMock: func(m *weatherMocks) {
				m.user.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				m.history.On("ListHistory", mock.Anything, userID, "Almaty", 10, 0).
					Return(nil, repoErr).Once()
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, svc := newService(t)
			if tt.setupMock != nil {
				tt.setupMock(m)
			}

			got, err := svc.GetHistory(context.Background(), tt.userID, tt.cityName, tt.limit, tt.offset)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
