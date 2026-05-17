package city_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/httpx"
	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/v1/city"
	domaincity "github.com/DaniilKalts/rbk-school/6-week/internal/domain/city"
	domainuser "github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
	servicecity "github.com/DaniilKalts/rbk-school/6-week/internal/service/city"
	"github.com/DaniilKalts/rbk-school/6-week/pkg/jwt"
)

func newRouter(t *testing.T) (*mockService, http.Handler) {
	t.Helper()

	svc := new(mockService)
	t.Cleanup(func() { svc.AssertExpectations(t) })

	r := chi.NewRouter()
	city.RegisterRoutes(r, svc)
	return svc, r
}

func doAuthed(t *testing.T, r http.Handler, userID uuid.UUID, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req = req.WithContext(httpx.WithClaims(req.Context(), &jwt.Claims{UserID: userID}))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestHandler_Create(t *testing.T) {
	userID := uuid.New()
	cityID := uuid.New()
	createdAt := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)
	created := &domaincity.City{
		ID:        cityID,
		UserID:    userID,
		Name:      "Almaty",
		CreatedAt: createdAt,
	}
	validBody := `{"city":"Almaty"}`

	tests := []struct {
		name       string
		claimsID   uuid.UUID
		body       string
		setupMock  func(*mockService)
		wantStatus int
		wantBody   string
	}{
		{
			name:     "success",
			claimsID: userID,
			body:     validBody,
			setupMock: func(s *mockService) {
				s.On("Create", mock.Anything, userID, servicecity.CreateInput{Name: "Almaty"}).
					Return(created, nil).Once()
			},
			wantStatus: http.StatusCreated,
			wantBody: `{"id":"` + cityID.String() + `","user_id":"` + userID.String() +
				`","city":"Almaty","created_at":"2026-05-17T12:00:00Z"}`,
		},
		{
			name:       "missing claims",
			claimsID:   uuid.Nil,
			body:       validBody,
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"code":401,"message":"отсутствуют claims аутентификации"}`,
		},
		{
			name:       "malformed JSON",
			claimsID:   userID,
			body:       `{`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"некорректное тело запроса"}`,
		},
		{
			name:     "invalid name maps to 400",
			claimsID: userID,
			body:     validBody,
			setupMock: func(s *mockService) {
				s.On("Create", mock.Anything, userID, mock.Anything).
					Return(nil, domaincity.ErrInvalidName).Once()
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"` + domaincity.ErrInvalidName.Error() + `"}`,
		},
		{
			name:     "user not found maps to 404",
			claimsID: userID,
			body:     validBody,
			setupMock: func(s *mockService) {
				s.On("Create", mock.Anything, userID, mock.Anything).
					Return(nil, domainuser.ErrNotFound).Once()
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":404,"message":"` + domainuser.ErrNotFound.Error() + `"}`,
		},
		{
			name:     "already exists maps to 409",
			claimsID: userID,
			body:     validBody,
			setupMock: func(s *mockService) {
				s.On("Create", mock.Anything, userID, mock.Anything).
					Return(nil, domaincity.ErrAlreadyExists).Once()
			},
			wantStatus: http.StatusConflict,
			wantBody:   `{"code":409,"message":"` + domaincity.ErrAlreadyExists.Error() + `"}`,
		},
		{
			name:     "unexpected error maps to 500",
			claimsID: userID,
			body:     validBody,
			setupMock: func(s *mockService) {
				s.On("Create", mock.Anything, userID, mock.Anything).
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

			w := doAuthed(t, r, tt.claimsID, http.MethodPost, "/cities", tt.body)

			require.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestHandler_List(t *testing.T) {
	userID := uuid.New()
	cityID := uuid.New()
	createdAt := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)
	cities := []domaincity.City{
		{ID: cityID, UserID: userID, Name: "Almaty", CreatedAt: createdAt},
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
				s.On("List", mock.Anything, userID).Return(cities, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody: `[{"id":"` + cityID.String() + `","user_id":"` + userID.String() +
				`","city":"Almaty","created_at":"2026-05-17T12:00:00Z"}]`,
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
				s.On("List", mock.Anything, userID).Return(nil, domainuser.ErrNotFound).Once()
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":404,"message":"` + domainuser.ErrNotFound.Error() + `"}`,
		},
		{
			name:     "unexpected error maps to 500",
			claimsID: userID,
			setupMock: func(s *mockService) {
				s.On("List", mock.Anything, userID).Return(nil, errors.New("db down")).Once()
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

			w := doAuthed(t, r, tt.claimsID, http.MethodGet, "/cities", "")

			require.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	userID := uuid.New()
	cityID := uuid.New()

	tests := []struct {
		name       string
		claimsID   uuid.UUID
		path       string
		setupMock  func(*mockService)
		wantStatus int
		wantBody   string
	}{
		{
			name:     "success",
			claimsID: userID,
			path:     "/cities/" + cityID.String(),
			setupMock: func(s *mockService) {
				s.On("Delete", mock.Anything, userID, cityID).Return(nil).Once()
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "missing claims",
			claimsID:   uuid.Nil,
			path:       "/cities/" + cityID.String(),
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"code":401,"message":"отсутствуют claims аутентификации"}`,
		},
		{
			name:       "invalid city id maps to 400",
			claimsID:   userID,
			path:       "/cities/not-a-uuid",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"invalid city id"}`,
		},
		{
			name:     "city not found maps to 404",
			claimsID: userID,
			path:     "/cities/" + cityID.String(),
			setupMock: func(s *mockService) {
				s.On("Delete", mock.Anything, userID, cityID).Return(domaincity.ErrNotFound).Once()
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":404,"message":"` + domaincity.ErrNotFound.Error() + `"}`,
		},
		{
			name:     "user not found maps to 404",
			claimsID: userID,
			path:     "/cities/" + cityID.String(),
			setupMock: func(s *mockService) {
				s.On("Delete", mock.Anything, userID, cityID).Return(domainuser.ErrNotFound).Once()
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":404,"message":"` + domainuser.ErrNotFound.Error() + `"}`,
		},
		{
			name:     "unexpected error maps to 500",
			claimsID: userID,
			path:     "/cities/" + cityID.String(),
			setupMock: func(s *mockService) {
				s.On("Delete", mock.Anything, userID, cityID).Return(errors.New("db down")).Once()
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

			w := doAuthed(t, r, tt.claimsID, http.MethodDelete, tt.path, "")

			require.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != "" {
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
				assert.JSONEq(t, tt.wantBody, w.Body.String())
			} else {
				assert.Empty(t, w.Body.String())
			}
		})
	}
}
