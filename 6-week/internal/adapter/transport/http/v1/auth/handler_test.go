package auth_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/v1/auth"
	domainuser "github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
	serviceauth "github.com/DaniilKalts/rbk-school/6-week/internal/service/auth"
	"github.com/DaniilKalts/rbk-school/6-week/pkg/jwt"
)

func newRouter(t *testing.T) (*mockService, http.Handler) {
	t.Helper()

	svc := new(mockService)
	t.Cleanup(func() { svc.AssertExpectations(t) })

	r := chi.NewRouter()
	auth.RegisterRoutes(r, svc)
	return svc, r
}

func do(t *testing.T, r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	var rdr *strings.Reader
	if body != "" {
		rdr = strings.NewReader(body)
	}
	var req *http.Request
	if rdr != nil {
		req = httptest.NewRequest(method, path, rdr)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestHandler_Register(t *testing.T) {
	expiresAt := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)
	token := &serviceauth.Token{AccessToken: "tok-abc", ExpiresAt: expiresAt}
	validBody := `{"first_name":"Daniil","last_name":"Kalts","email":"daniil@example.com","password":"secret"}`

	tests := []struct {
		name       string
		body       string
		setupMock  func(*mockService)
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			body: validBody,
			setupMock: func(s *mockService) {
				s.On("Register", mock.Anything, serviceauth.RegisterInput{
					FirstName: "Daniil",
					LastName:  "Kalts",
					Email:     "daniil@example.com",
					Password:  "secret",
				}).Return(token, nil).Once()
			},
			wantStatus: http.StatusCreated,
			wantBody:   `{"access_token":"tok-abc","expires_at":"2026-05-17T12:00:00Z"}`,
		},
		{
			name:       "malformed JSON",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"некорректное тело запроса"}`,
		},
		{
			name: "invalid email maps to 400",
			body: validBody,
			setupMock: func(s *mockService) {
				s.On("Register", mock.Anything, mock.Anything).
					Return(nil, domainuser.ErrInvalidEmail).Once()
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"` + domainuser.ErrInvalidEmail.Error() + `"}`,
		},
		{
			name: "email already exists maps to 409",
			body: validBody,
			setupMock: func(s *mockService) {
				s.On("Register", mock.Anything, mock.Anything).
					Return(nil, domainuser.ErrEmailAlreadyExists).Once()
			},
			wantStatus: http.StatusConflict,
			wantBody:   `{"code":409,"message":"` + domainuser.ErrEmailAlreadyExists.Error() + `"}`,
		},
		{
			name: "unexpected error maps to 500",
			body: validBody,
			setupMock: func(s *mockService) {
				s.On("Register", mock.Anything, mock.Anything).
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

			w := do(t, r, http.MethodPost, "/auth/register", tt.body)

			require.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestHandler_Login(t *testing.T) {
	expiresAt := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)
	token := &serviceauth.Token{AccessToken: "tok-abc", ExpiresAt: expiresAt}
	validBody := `{"email":"daniil@example.com","password":"secret"}`

	tests := []struct {
		name       string
		body       string
		setupMock  func(*mockService)
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			body: validBody,
			setupMock: func(s *mockService) {
				s.On("Login", mock.Anything, serviceauth.LoginInput{
					Email:    "daniil@example.com",
					Password: "secret",
				}).Return(token, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"access_token":"tok-abc","expires_at":"2026-05-17T12:00:00Z"}`,
		},
		{
			name:       "malformed JSON",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"некорректное тело запроса"}`,
		},
		{
			name: "invalid email maps to 400",
			body: validBody,
			setupMock: func(s *mockService) {
				s.On("Login", mock.Anything, mock.Anything).
					Return(nil, domainuser.ErrInvalidEmail).Once()
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"` + domainuser.ErrInvalidEmail.Error() + `"}`,
		},
		{
			name: "invalid credentials maps to 401",
			body: validBody,
			setupMock: func(s *mockService) {
				s.On("Login", mock.Anything, mock.Anything).
					Return(nil, serviceauth.ErrInvalidCredentials).Once()
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"code":401,"message":"` + serviceauth.ErrInvalidCredentials.Error() + `"}`,
		},
		{
			name: "unexpected error maps to 500",
			body: validBody,
			setupMock: func(s *mockService) {
				s.On("Login", mock.Anything, mock.Anything).
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

			w := do(t, r, http.MethodPost, "/auth/login", tt.body)

			require.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestHandler_Logout(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		setupMock  func(*mockService)
		wantStatus int
		wantBody   string
	}{
		{
			name:       "success",
			authHeader: "Bearer tok-abc",
			setupMock: func(s *mockService) {
				s.On("Logout", mock.Anything, "tok-abc").Return(nil).Once()
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "missing Authorization header",
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"code":401,"message":"отсутствует или некорректный заголовок Authorization"}`,
		},
		{
			name:       "malformed Authorization header",
			authHeader: "Basic xyz",
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"code":401,"message":"отсутствует или некорректный заголовок Authorization"}`,
		},
		{
			name:       "invalid token maps to 401",
			authHeader: "Bearer tok-abc",
			setupMock: func(s *mockService) {
				s.On("Logout", mock.Anything, "tok-abc").Return(jwt.ErrInvalidToken).Once()
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"code":401,"message":"` + jwt.ErrInvalidToken.Error() + `"}`,
		},
		{
			name:       "unexpected error maps to 500",
			authHeader: "Bearer tok-abc",
			setupMock: func(s *mockService) {
				s.On("Logout", mock.Anything, "tok-abc").Return(errors.New("redis down")).Once()
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

			req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

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
