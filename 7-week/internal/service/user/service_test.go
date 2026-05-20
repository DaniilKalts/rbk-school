package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainuser "github.com/DaniilKalts/rbk-school/7-week/internal/domain/user"
	serviceuser "github.com/DaniilKalts/rbk-school/7-week/internal/service/user"
)

func newService(t *testing.T) (*mockRepository, *serviceuser.Service) {
	t.Helper()

	repo := new(mockRepository)
	t.Cleanup(func() { repo.AssertExpectations(t) })

	return repo, serviceuser.NewService(repo)
}

func TestService_Create(t *testing.T) {
	created := &domainuser.User{
		FirstName: "Daniil",
		LastName:  "Kalts",
		Email:     "daniil@example.com",
		Role:      domainuser.RoleUser,
	}
	repoErr := errors.New("db down")

	tests := []struct {
		name      string
		input     serviceuser.CreateInput
		setupMock func(*mockRepository)
		want      *domainuser.User
		wantErr   error
	}{
		{
			name: "success normalizes and trims fields",
			input: serviceuser.CreateInput{
				FirstName: " Daniil ",
				LastName:  " Kalts ",
				Email:     " DANIIL@example.com ",
				Password:  "secret",
			},
			setupMock: func(r *mockRepository) {
				r.On("Create",
					mock.Anything,
					mock.MatchedBy(func(u domainuser.User) bool {
						return u.FirstName == "Daniil" &&
							u.LastName == "Kalts" &&
							u.Email == "daniil@example.com" &&
							u.Role == domainuser.RoleUser
					}),
					mock.Anything,
				).Return(created, nil).Once()
			},
			want: created,
		},
		{
			name: "empty password",
			input: serviceuser.CreateInput{
				FirstName: "Daniil",
				LastName:  "Kalts",
				Email:     "daniil@example.com",
				Password:  "",
			},
			wantErr: domainuser.ErrInvalidPassword,
		},
		{
			name: "invalid email",
			input: serviceuser.CreateInput{
				FirstName: "Daniil",
				LastName:  "Kalts",
				Email:     "not-an-email",
				Password:  "secret",
			},
			wantErr: domainuser.ErrInvalidEmail,
		},
		{
			name: "repo error propagates",
			input: serviceuser.CreateInput{
				FirstName: "Daniil",
				LastName:  "Kalts",
				Email:     "daniil@example.com",
				Password:  "secret",
			},
			setupMock: func(r *mockRepository) {
				r.On("Create", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, repoErr).Once()
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, svc := newService(t)
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			got, err := svc.Create(context.Background(), tt.input)

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

func TestService_GetByID(t *testing.T) {
	id := uuid.New()
	expected := &domainuser.User{ID: id, Email: "daniil@example.com"}

	tests := []struct {
		name      string
		id        uuid.UUID
		setupMock func(*mockRepository)
		want      *domainuser.User
		wantErr   error
	}{
		{
			name: "success",
			id:   id,
			setupMock: func(r *mockRepository) {
				r.On("GetByID", mock.Anything, id).Return(expected, nil).Once()
			},
			want: expected,
		},
		{
			name:    "nil id",
			id:      uuid.Nil,
			wantErr: domainuser.ErrInvalidID,
		},
		{
			name: "not found",
			id:   id,
			setupMock: func(r *mockRepository) {
				r.On("GetByID", mock.Anything, id).Return(nil, domainuser.ErrNotFound).Once()
			},
			wantErr: domainuser.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, svc := newService(t)
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			got, err := svc.GetByID(context.Background(), tt.id)

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

func TestService_GetByEmail(t *testing.T) {
	expected := &domainuser.User{Email: "daniil@example.com"}

	tests := []struct {
		name      string
		email     string
		setupMock func(*mockRepository)
		want      *domainuser.User
		wantErr   error
	}{
		{
			name:  "success normalizes email",
			email: " DANIIL@example.com ",
			setupMock: func(r *mockRepository) {
				r.On("GetByEmail", mock.Anything, "daniil@example.com").Return(expected, nil).Once()
			},
			want: expected,
		},
		{
			name:    "blank email",
			email:   "   ",
			wantErr: domainuser.ErrInvalidEmail,
		},
		{
			name:  "not found",
			email: "daniil@example.com",
			setupMock: func(r *mockRepository) {
				r.On("GetByEmail", mock.Anything, "daniil@example.com").
					Return(nil, domainuser.ErrNotFound).Once()
			},
			wantErr: domainuser.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, svc := newService(t)
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			got, err := svc.GetByEmail(context.Background(), tt.email)

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

func TestService_List(t *testing.T) {
	expected := []domainuser.User{
		{FirstName: "Daniil", LastName: "Kalts", Email: "daniil@example.com", Role: domainuser.RoleAdmin},
		{FirstName: "Martin", LastName: "Kalts", Email: "martin@example.com", Role: domainuser.RoleUser},
	}
	repoErr := errors.New("db down")

	tests := []struct {
		name      string
		setupMock func(*mockRepository)
		want      []domainuser.User
		wantErr   error
	}{
		{
			name: "success",
			setupMock: func(r *mockRepository) {
				r.On("List", mock.Anything).Return(expected, nil).Once()
			},
			want: expected,
		},
		{
			name: "repo error propagates",
			setupMock: func(r *mockRepository) {
				r.On("List", mock.Anything).Return(nil, repoErr).Once()
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, svc := newService(t)
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			got, err := svc.List(context.Background())

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

func TestService_Update(t *testing.T) {
	id := uuid.New()
	validInput := serviceuser.UpdateInput{
		FirstName: "New",
		LastName:  "Surname",
		Email:     "NEW@example.com",
	}
	repoErr := errors.New("db down")

	newExisting := func() *domainuser.User {
		return &domainuser.User{
			ID:        id,
			FirstName: "Old",
			LastName:  "Name",
			Email:     "old@example.com",
			Role:      domainuser.RoleUser,
		}
	}

	tests := []struct {
		name      string
		id        uuid.UUID
		input     serviceuser.UpdateInput
		setupMock func(*mockRepository, *domainuser.User)
		wantErr   error
	}{
		{
			name:  "success",
			id:    id,
			input: validInput,
			setupMock: func(r *mockRepository, existing *domainuser.User) {
				r.On("GetByID", mock.Anything, id).Return(existing, nil).Once()
				r.On("Update", mock.Anything,
					mock.MatchedBy(func(u domainuser.User) bool {
						return u.ID == id &&
							u.FirstName == "New" &&
							u.LastName == "Surname" &&
							u.Email == "new@example.com"
					}),
				).Return(existing, nil).Once()
			},
		},
		{
			name:    "nil id",
			id:      uuid.Nil,
			input:   validInput,
			wantErr: domainuser.ErrInvalidID,
		},
		{
			name:  "not found",
			id:    id,
			input: validInput,
			setupMock: func(r *mockRepository, _ *domainuser.User) {
				r.On("GetByID", mock.Anything, id).Return(nil, domainuser.ErrNotFound).Once()
			},
			wantErr: domainuser.ErrNotFound,
		},
		{
			name: "invalid email",
			id:   id,
			input: serviceuser.UpdateInput{
				FirstName: "New",
				LastName:  "Surname",
				Email:     "not-an-email",
			},
			setupMock: func(r *mockRepository, existing *domainuser.User) {
				r.On("GetByID", mock.Anything, id).Return(existing, nil).Once()
			},
			wantErr: domainuser.ErrInvalidEmail,
		},
		{
			name:  "repo update error propagates",
			id:    id,
			input: validInput,
			setupMock: func(r *mockRepository, existing *domainuser.User) {
				r.On("GetByID", mock.Anything, id).Return(existing, nil).Once()
				r.On("Update", mock.Anything, mock.Anything).Return(nil, repoErr).Once()
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, svc := newService(t)
			existing := newExisting()
			if tt.setupMock != nil {
				tt.setupMock(repo, existing)
			}

			got, err := svc.Update(context.Background(), tt.id, tt.input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, existing, got)
		})
	}
}

func TestService_Delete(t *testing.T) {
	id := uuid.New()
	repoErr := errors.New("db down")

	tests := []struct {
		name      string
		id        uuid.UUID
		setupMock func(*mockRepository)
		wantErr   error
	}{
		{
			name: "success",
			id:   id,
			setupMock: func(r *mockRepository) {
				r.On("SoftDelete", mock.Anything, id).Return(nil).Once()
			},
		},
		{
			name:    "nil id",
			id:      uuid.Nil,
			wantErr: domainuser.ErrInvalidID,
		},
		{
			name: "repo error propagates",
			id:   id,
			setupMock: func(r *mockRepository) {
				r.On("SoftDelete", mock.Anything, id).Return(repoErr).Once()
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, svc := newService(t)
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			err := svc.Delete(context.Background(), tt.id)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
