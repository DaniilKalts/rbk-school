package weather_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/rbk-school/7-week/pkg/httpx"
	"github.com/DaniilKalts/rbk-school/7-week/api-service/internal/adapter/transport/http/v1/weather"
	domainhistory "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/history"
	domainuser "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/user"
	domainweather "github.com/DaniilKalts/rbk-school/7-week/api-service/internal/domain/weather"
	"github.com/DaniilKalts/rbk-school/7-week/pkg/jwt"
)

func newRouter(t *testing.T) (*mockService, http.Handler) {
	t.Helper()

	svc := new(mockService)
	t.Cleanup(func() { svc.AssertExpectations(t) })

	r := chi.NewRouter()
	weather.RegisterRoutes(r, svc)
	return svc, r
}

func doAuthed(t *testing.T, r http.Handler, userID uuid.UUID, method, path string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, nil)
	req = req.WithContext(httpx.WithClaims(req.Context(), &jwt.Claims{UserID: userID}))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestHandler_GetByUser(t *testing.T) {
	userID := uuid.New()
	requestedAt := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)
	weathers := []domainweather.Weather{
		{City: "Almaty", Temperature: 20, FeelsLike: 19, Description: "clear sky", RequestedAt: requestedAt},
	}

	tests := []struct {
		name       string
		claimsID   uuid.UUID
		setupMock  func(*mockService)
		wantStatus int
		wantBody   string
	}{
		{
			name:     "success",
			claimsID: userID,
			setupMock: func(s *mockService) {
				s.On("GetByUserID", mock.Anything, userID).Return(weathers, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody: `{"user_id":"` + userID.String() + `","weather":[` +
				`{"city":"Almaty","temperature":20,"feels_like":19,"description":"clear sky","requested_at":"2026-05-17T12:00:00Z"}` +
				`]}`,
		},
		{
			name:       "missing claims",
			claimsID:   uuid.Nil,
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"code":401,"message":"отсутствуют claims аутентификации"}`,
		},
		{
			name:     "user not found maps to 404",
			claimsID: userID,
			setupMock: func(s *mockService) {
				s.On("GetByUserID", mock.Anything, userID).Return(nil, domainuser.ErrNotFound).Once()
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":404,"message":"` + domainuser.ErrNotFound.Error() + `"}`,
		},
		{
			name:     "unexpected error maps to 500",
			claimsID: userID,
			setupMock: func(s *mockService) {
				s.On("GetByUserID", mock.Anything, userID).Return(nil, errors.New("upstream down")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"code":500,"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, r := newRouter(t)
			if tt.setupMock != nil {
				tt.setupMock(svc)
			}

			w := doAuthed(t, r, tt.claimsID, http.MethodGet, "/weather")

			require.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestHandler_History(t *testing.T) {
	userID := uuid.New()
	requestedAt := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)
	items := []domainhistory.History{
		{City: "Almaty", Temperature: 20, Description: "clear sky", RequestedAt: requestedAt},
	}

	tests := []struct {
		name       string
		claimsID   uuid.UUID
		query      string
		setupMock  func(*mockService)
		wantStatus int
		wantBody   string
	}{
		{
			name:     "success without filters",
			claimsID: userID,
			query:    "",
			setupMock: func(s *mockService) {
				s.On("GetHistory", mock.Anything, userID, "", 0, 0).Return(items, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody: `{"user_id":"` + userID.String() + `","history":[` +
				`{"city":"Almaty","temperature":20,"description":"clear sky","requested_at":"2026-05-17T12:00:00Z"}` +
				`]}`,
		},
		{
			name:     "success with filters normalizes city",
			claimsID: userID,
			query:    "?city=%20almaty%20&limit=5&offset=2",
			setupMock: func(s *mockService) {
				s.On("GetHistory", mock.Anything, userID, "almaty", 5, 2).Return(items, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody: `{"user_id":"` + userID.String() + `","city":"Almaty","history":[` +
				`{"city":"Almaty","temperature":20,"description":"clear sky","requested_at":"2026-05-17T12:00:00Z"}` +
				`]}`,
		},
		{
			name:       "missing claims",
			claimsID:   uuid.Nil,
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"code":401,"message":"отсутствуют claims аутентификации"}`,
		},
		{
			name:       "invalid limit maps to 400",
			claimsID:   userID,
			query:      "?limit=abc",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"limit должен быть положительным числом"}`,
		},
		{
			name:       "invalid offset maps to 400",
			claimsID:   userID,
			query:      "?offset=-1",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"offset должен быть неотрицательным числом"}`,
		},
		{
			name:     "service validation error maps to 400",
			claimsID: userID,
			setupMock: func(s *mockService) {
				s.On("GetHistory", mock.Anything, userID, "", 0, 0).
					Return(nil, domainweather.ErrInvalidLimit).Once()
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"` + domainweather.ErrInvalidLimit.Error() + `"}`,
		},
		{
			name:     "user not found maps to 404",
			claimsID: userID,
			setupMock: func(s *mockService) {
				s.On("GetHistory", mock.Anything, userID, "", 0, 0).
					Return(nil, domainuser.ErrNotFound).Once()
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":404,"message":"` + domainuser.ErrNotFound.Error() + `"}`,
		},
		{
			name:     "unexpected error maps to 500",
			claimsID: userID,
			setupMock: func(s *mockService) {
				s.On("GetHistory", mock.Anything, userID, "", 0, 0).
					Return(nil, errors.New("db down")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"code":500,"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, r := newRouter(t)
			if tt.setupMock != nil {
				tt.setupMock(svc)
			}

			w := doAuthed(t, r, tt.claimsID, http.MethodGet, "/weather/history"+tt.query)

			require.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}
