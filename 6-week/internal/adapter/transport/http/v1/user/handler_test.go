package user_test

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
	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/v1/user"
	domainuser "github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
	serviceuser "github.com/DaniilKalts/rbk-school/6-week/internal/service/user"
	"github.com/DaniilKalts/rbk-school/6-week/pkg/jwt"
)

func newRouter(t *testing.T) (*mockService, *mockTokenRevoker, http.Handler) {
	t.Helper()

	svc := new(mockService)
	tr := new(mockTokenRevoker)
	t.Cleanup(func() {
		svc.AssertExpectations(t)
		tr.AssertExpectations(t)
	})

	r := chi.NewRouter()
	user.RegisterCurrentUserRoutes(r, svc, tr)
	user.RegisterAdminRoutes(r, svc, tr)
	return svc, tr, r
}

func fixedUser(id uuid.UUID) *domainuser.User {
	createdAt := time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC)
	return &domainuser.User{
		ID:        id,
		FirstName: "Daniil",
		LastName:  "Kalts",
		Email:     "daniil@example.com",
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

func userResponseJSON(u *domainuser.User) string {
	return `{"id":"` + u.ID.String() +
		`","first_name":"` + u.FirstName +
		`","last_name":"` + u.LastName +
		`","email":"` + u.Email +
		`","created_at":"2026-05-17T12:00:00Z","updated_at":"2026-05-17T12:00:00Z"}`
}

func authedRequest(method, path, body string, claimsID uuid.UUID, bearer string) *http.Request {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	if claimsID != uuid.Nil {
		req = req.WithContext(httpx.WithClaims(req.Context(), &jwt.Claims{UserID: claimsID}))
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	return req
}

func TestHandler_GetMe(t *testing.T) {
	userID := uuid.New()
	u := fixedUser(userID)

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
				s.On("GetByID", mock.Anything, userID).Return(u, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody:   userResponseJSON(u),
		},
		{
			name:       "missing claims",
			claimsID:   uuid.Nil,
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"code":401,"message":"отсутствуют claims аутентификации"}`,
		},
		{
			name:     "not found maps to 404",
			claimsID: userID,
			setupMock: func(s *mockService) {
				s.On("GetByID", mock.Anything, userID).Return(nil, domainuser.ErrNotFound).Once()
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":404,"message":"` + domainuser.ErrNotFound.Error() + `"}`,
		},
		{
			name:     "unexpected error maps to 500",
			claimsID: userID,
			setupMock: func(s *mockService) {
				s.On("GetByID", mock.Anything, userID).Return(nil, errors.New("db down")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"code":500,"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _, r := newRouter(t)
			if tt.setupMock != nil {
				tt.setupMock(svc)
			}

			req := authedRequest(http.MethodGet, "/users/me", "", tt.claimsID, "")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestHandler_UpdateMe(t *testing.T) {
	userID := uuid.New()
	u := fixedUser(userID)
	validBody := `{"first_name":"Daniil","last_name":"Kalts","email":"daniil@example.com"}`

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
				s.On("Update", mock.Anything, userID, serviceuser.UpdateInput{
					FirstName: "Daniil",
					LastName:  "Kalts",
					Email:     "daniil@example.com",
				}).Return(u, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody:   userResponseJSON(u),
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
			name:     "invalid email maps to 400",
			claimsID: userID,
			body:     validBody,
			setupMock: func(s *mockService) {
				s.On("Update", mock.Anything, userID, mock.Anything).
					Return(nil, domainuser.ErrInvalidEmail).Once()
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"` + domainuser.ErrInvalidEmail.Error() + `"}`,
		},
		{
			name:     "not found maps to 404",
			claimsID: userID,
			body:     validBody,
			setupMock: func(s *mockService) {
				s.On("Update", mock.Anything, userID, mock.Anything).
					Return(nil, domainuser.ErrNotFound).Once()
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":404,"message":"` + domainuser.ErrNotFound.Error() + `"}`,
		},
		{
			name:     "email taken maps to 409",
			claimsID: userID,
			body:     validBody,
			setupMock: func(s *mockService) {
				s.On("Update", mock.Anything, userID, mock.Anything).
					Return(nil, domainuser.ErrEmailAlreadyExists).Once()
			},
			wantStatus: http.StatusConflict,
			wantBody:   `{"code":409,"message":"` + domainuser.ErrEmailAlreadyExists.Error() + `"}`,
		},
		{
			name:     "unexpected error maps to 500",
			claimsID: userID,
			body:     validBody,
			setupMock: func(s *mockService) {
				s.On("Update", mock.Anything, userID, mock.Anything).
					Return(nil, errors.New("db down")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"code":500,"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _, r := newRouter(t)
			if tt.setupMock != nil {
				tt.setupMock(svc)
			}

			req := authedRequest(http.MethodPatch, "/users/me", tt.body, tt.claimsID, "")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestHandler_DeleteMe(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name       string
		claimsID   uuid.UUID
		bearer     string
		setupMock  func(*mockService, *mockTokenRevoker)
		wantStatus int
		wantBody   string
	}{
		{
			name:     "success",
			claimsID: userID,
			bearer:   "tok-abc",
			setupMock: func(s *mockService, tr *mockTokenRevoker) {
				tr.On("Revoke", mock.Anything, "tok-abc").Return(nil).Once()
				s.On("Delete", mock.Anything, userID).Return(nil).Once()
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "missing claims",
			claimsID:   uuid.Nil,
			bearer:     "tok-abc",
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"code":401,"message":"отсутствуют claims аутентификации"}`,
		},
		{
			name:       "missing Authorization header",
			claimsID:   userID,
			bearer:     "",
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"code":401,"message":"отсутствует или некорректный заголовок Authorization"}`,
		},
		{
			name:     "invalid token maps to 401",
			claimsID: userID,
			bearer:   "tok-abc",
			setupMock: func(_ *mockService, tr *mockTokenRevoker) {
				tr.On("Revoke", mock.Anything, "tok-abc").Return(jwt.ErrInvalidToken).Once()
			},
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"code":401,"message":"` + jwt.ErrInvalidToken.Error() + `"}`,
		},
		{
			name:     "revoke unexpected error maps to 500",
			claimsID: userID,
			bearer:   "tok-abc",
			setupMock: func(_ *mockService, tr *mockTokenRevoker) {
				tr.On("Revoke", mock.Anything, "tok-abc").Return(errors.New("redis down")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"code":500,"message":"internal server error"}`,
		},
		{
			name:     "delete not found maps to 404",
			claimsID: userID,
			bearer:   "tok-abc",
			setupMock: func(s *mockService, tr *mockTokenRevoker) {
				tr.On("Revoke", mock.Anything, "tok-abc").Return(nil).Once()
				s.On("Delete", mock.Anything, userID).Return(domainuser.ErrNotFound).Once()
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":404,"message":"` + domainuser.ErrNotFound.Error() + `"}`,
		},
		{
			name:     "delete unexpected error maps to 500",
			claimsID: userID,
			bearer:   "tok-abc",
			setupMock: func(s *mockService, tr *mockTokenRevoker) {
				tr.On("Revoke", mock.Anything, "tok-abc").Return(nil).Once()
				s.On("Delete", mock.Anything, userID).Return(errors.New("db down")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"code":500,"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, tr, r := newRouter(t)
			if tt.setupMock != nil {
				tt.setupMock(svc, tr)
			}

			req := authedRequest(http.MethodDelete, "/users/me", "", tt.claimsID, tt.bearer)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, w.Body.String())
			} else {
				assert.Empty(t, w.Body.String())
			}
		})
	}
}

func TestHandler_Create(t *testing.T) {
	userID := uuid.New()
	u := fixedUser(userID)
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
				s.On("Create", mock.Anything, serviceuser.CreateInput{
					FirstName: "Daniil",
					LastName:  "Kalts",
					Email:     "daniil@example.com",
					Password:  "secret",
				}).Return(u, nil).Once()
			},
			wantStatus: http.StatusCreated,
			wantBody:   userResponseJSON(u),
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
				s.On("Create", mock.Anything, mock.Anything).
					Return(nil, domainuser.ErrInvalidEmail).Once()
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"` + domainuser.ErrInvalidEmail.Error() + `"}`,
		},
		{
			name: "email taken maps to 409",
			body: validBody,
			setupMock: func(s *mockService) {
				s.On("Create", mock.Anything, mock.Anything).
					Return(nil, domainuser.ErrEmailAlreadyExists).Once()
			},
			wantStatus: http.StatusConflict,
			wantBody:   `{"code":409,"message":"` + domainuser.ErrEmailAlreadyExists.Error() + `"}`,
		},
		{
			name: "unexpected error maps to 500",
			body: validBody,
			setupMock: func(s *mockService) {
				s.On("Create", mock.Anything, mock.Anything).
					Return(nil, errors.New("db down")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"code":500,"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _, r := newRouter(t)
			if tt.setupMock != nil {
				tt.setupMock(svc)
			}

			req := authedRequest(http.MethodPost, "/users", tt.body, uuid.Nil, "")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestHandler_List(t *testing.T) {
	id := uuid.New()
	users := []domainuser.User{*fixedUser(id)}

	tests := []struct {
		name       string
		setupMock  func(*mockService)
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			setupMock: func(s *mockService) {
				s.On("List", mock.Anything).Return(users, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody:   `[` + userResponseJSON(&users[0]) + `]`,
		},
		{
			name: "unexpected error maps to 500",
			setupMock: func(s *mockService) {
				s.On("List", mock.Anything).Return(nil, errors.New("db down")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"code":500,"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _, r := newRouter(t)
			tt.setupMock(svc)

			req := authedRequest(http.MethodGet, "/users", "", uuid.Nil, "")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestHandler_GetByID(t *testing.T) {
	id := uuid.New()
	u := fixedUser(id)

	tests := []struct {
		name       string
		path       string
		setupMock  func(*mockService)
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			path: "/users/" + id.String(),
			setupMock: func(s *mockService) {
				s.On("GetByID", mock.Anything, id).Return(u, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody:   userResponseJSON(u),
		},
		{
			name:       "invalid id maps to 400",
			path:       "/users/not-a-uuid",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"некорректный id пользователя"}`,
		},
		{
			name: "not found maps to 404",
			path: "/users/" + id.String(),
			setupMock: func(s *mockService) {
				s.On("GetByID", mock.Anything, id).Return(nil, domainuser.ErrNotFound).Once()
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":404,"message":"` + domainuser.ErrNotFound.Error() + `"}`,
		},
		{
			name: "unexpected error maps to 500",
			path: "/users/" + id.String(),
			setupMock: func(s *mockService) {
				s.On("GetByID", mock.Anything, id).Return(nil, errors.New("db down")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"code":500,"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _, r := newRouter(t)
			if tt.setupMock != nil {
				tt.setupMock(svc)
			}

			req := authedRequest(http.MethodGet, tt.path, "", uuid.Nil, "")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestHandler_Update(t *testing.T) {
	id := uuid.New()
	u := fixedUser(id)
	validBody := `{"first_name":"Daniil","last_name":"Kalts","email":"daniil@example.com"}`

	tests := []struct {
		name       string
		path       string
		body       string
		setupMock  func(*mockService)
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			path: "/users/" + id.String(),
			body: validBody,
			setupMock: func(s *mockService) {
				s.On("Update", mock.Anything, id, serviceuser.UpdateInput{
					FirstName: "Daniil",
					LastName:  "Kalts",
					Email:     "daniil@example.com",
				}).Return(u, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody:   userResponseJSON(u),
		},
		{
			name:       "invalid id maps to 400",
			path:       "/users/not-a-uuid",
			body:       validBody,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"некорректный id пользователя"}`,
		},
		{
			name:       "malformed JSON",
			path:       "/users/" + id.String(),
			body:       `{`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"некорректное тело запроса"}`,
		},
		{
			name: "invalid email maps to 400",
			path: "/users/" + id.String(),
			body: validBody,
			setupMock: func(s *mockService) {
				s.On("Update", mock.Anything, id, mock.Anything).
					Return(nil, domainuser.ErrInvalidEmail).Once()
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"` + domainuser.ErrInvalidEmail.Error() + `"}`,
		},
		{
			name: "not found maps to 404",
			path: "/users/" + id.String(),
			body: validBody,
			setupMock: func(s *mockService) {
				s.On("Update", mock.Anything, id, mock.Anything).
					Return(nil, domainuser.ErrNotFound).Once()
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":404,"message":"` + domainuser.ErrNotFound.Error() + `"}`,
		},
		{
			name: "email taken maps to 409",
			path: "/users/" + id.String(),
			body: validBody,
			setupMock: func(s *mockService) {
				s.On("Update", mock.Anything, id, mock.Anything).
					Return(nil, domainuser.ErrEmailAlreadyExists).Once()
			},
			wantStatus: http.StatusConflict,
			wantBody:   `{"code":409,"message":"` + domainuser.ErrEmailAlreadyExists.Error() + `"}`,
		},
		{
			name: "unexpected error maps to 500",
			path: "/users/" + id.String(),
			body: validBody,
			setupMock: func(s *mockService) {
				s.On("Update", mock.Anything, id, mock.Anything).
					Return(nil, errors.New("db down")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"code":500,"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _, r := newRouter(t)
			if tt.setupMock != nil {
				tt.setupMock(svc)
			}

			req := authedRequest(http.MethodPatch, tt.path, tt.body, uuid.Nil, "")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name       string
		path       string
		setupMock  func(*mockService)
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			path: "/users/" + id.String(),
			setupMock: func(s *mockService) {
				s.On("Delete", mock.Anything, id).Return(nil).Once()
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "invalid id maps to 400",
			path:       "/users/not-a-uuid",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":400,"message":"некорректный id пользователя"}`,
		},
		{
			name: "not found maps to 404",
			path: "/users/" + id.String(),
			setupMock: func(s *mockService) {
				s.On("Delete", mock.Anything, id).Return(domainuser.ErrNotFound).Once()
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":404,"message":"` + domainuser.ErrNotFound.Error() + `"}`,
		},
		{
			name: "unexpected error maps to 500",
			path: "/users/" + id.String(),
			setupMock: func(s *mockService) {
				s.On("Delete", mock.Anything, id).Return(errors.New("db down")).Once()
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"code":500,"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _, r := newRouter(t)
			if tt.setupMock != nil {
				tt.setupMock(svc)
			}

			req := authedRequest(http.MethodDelete, tt.path, "", uuid.Nil, "")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, w.Body.String())
			} else {
				assert.Empty(t, w.Body.String())
			}
		})
	}
}
