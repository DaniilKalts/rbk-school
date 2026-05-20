package city_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domaincity "github.com/DaniilKalts/rbk-school/7-week/internal/domain/city"
	domainuser "github.com/DaniilKalts/rbk-school/7-week/internal/domain/user"
	servicecity "github.com/DaniilKalts/rbk-school/7-week/internal/service/city"
)

func newService(t *testing.T) (*mockRepository, *mockUserRepository, *servicecity.Service) {
	t.Helper()

	repo := new(mockRepository)
	userRepo := new(mockUserRepository)
	t.Cleanup(func() {
		repo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	return repo, userRepo, servicecity.NewService(repo, userRepo)
}

func TestService_Create(t *testing.T) {
	userID := uuid.New()
	created := &domaincity.City{UserID: userID, Name: "Almaty"}

	tests := []struct {
		name      string
		userID    uuid.UUID
		input     servicecity.CreateInput
		setupMock func(*mockRepository, *mockUserRepository)
		want      *domaincity.City
		wantErr   error
	}{
		{
			name:   "success normalizes city name",
			userID: userID,
			input:  servicecity.CreateInput{Name: " almaty "},
			setupMock: func(r *mockRepository, ur *mockUserRepository) {
				ur.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				r.On("Create", mock.Anything,
					mock.MatchedBy(func(c domaincity.City) bool {
						return c.UserID == userID && c.Name == "Almaty"
					}),
				).Return(created, nil).Once()
			},
			want: created,
		},
		{
			name:    "nil user id",
			userID:  uuid.Nil,
			input:   servicecity.CreateInput{Name: "Almaty"},
			wantErr: domaincity.ErrInvalidUserID,
		},
		{
			name:    "blank name",
			userID:  userID,
			input:   servicecity.CreateInput{Name: "   "},
			wantErr: domaincity.ErrInvalidName,
		},
		{
			name:   "user not found",
			userID: userID,
			input:  servicecity.CreateInput{Name: "Almaty"},
			setupMock: func(_ *mockRepository, ur *mockUserRepository) {
				ur.On("GetByID", mock.Anything, userID).Return(nil, domainuser.ErrNotFound).Once()
			},
			wantErr: domainuser.ErrNotFound,
		},
		{
			name:   "repo error propagates",
			userID: userID,
			input:  servicecity.CreateInput{Name: "Almaty"},
			setupMock: func(r *mockRepository, ur *mockUserRepository) {
				ur.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				r.On("Create", mock.Anything, mock.Anything).
					Return(nil, domaincity.ErrAlreadyExists).Once()
			},
			wantErr: domaincity.ErrAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, userRepo, svc := newService(t)
			if tt.setupMock != nil {
				tt.setupMock(repo, userRepo)
			}

			got, err := svc.Create(context.Background(), tt.userID, tt.input)

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
	userID := uuid.New()
	expected := []domaincity.City{{Name: "Almaty"}, {Name: "Astana"}}
	repoErr := errors.New("db down")

	tests := []struct {
		name      string
		userID    uuid.UUID
		setupMock func(*mockRepository, *mockUserRepository)
		want      []domaincity.City
		wantErr   error
	}{
		{
			name:   "success",
			userID: userID,
			setupMock: func(r *mockRepository, ur *mockUserRepository) {
				ur.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				r.On("ListByUserID", mock.Anything, userID).Return(expected, nil).Once()
			},
			want: expected,
		},
		{
			name:    "nil user id",
			userID:  uuid.Nil,
			wantErr: domaincity.ErrInvalidUserID,
		},
		{
			name:   "user not found",
			userID: userID,
			setupMock: func(_ *mockRepository, ur *mockUserRepository) {
				ur.On("GetByID", mock.Anything, userID).Return(nil, domainuser.ErrNotFound).Once()
			},
			wantErr: domainuser.ErrNotFound,
		},
		{
			name:   "repo error propagates",
			userID: userID,
			setupMock: func(r *mockRepository, ur *mockUserRepository) {
				ur.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				r.On("ListByUserID", mock.Anything, userID).Return(nil, repoErr).Once()
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, userRepo, svc := newService(t)
			if tt.setupMock != nil {
				tt.setupMock(repo, userRepo)
			}

			got, err := svc.List(context.Background(), tt.userID)

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

func TestService_Delete(t *testing.T) {
	userID := uuid.New()
	cityID := uuid.New()
	repoErr := errors.New("db down")

	tests := []struct {
		name      string
		userID    uuid.UUID
		cityID    uuid.UUID
		setupMock func(*mockRepository, *mockUserRepository)
		wantErr   error
	}{
		{
			name:   "success",
			userID: userID,
			cityID: cityID,
			setupMock: func(r *mockRepository, ur *mockUserRepository) {
				ur.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				r.On("Delete", mock.Anything, userID, cityID).Return(nil).Once()
			},
		},
		{
			name:    "nil city id",
			userID:  userID,
			cityID:  uuid.Nil,
			wantErr: domaincity.ErrInvalidID,
		},
		{
			name:    "nil user id",
			userID:  uuid.Nil,
			cityID:  cityID,
			wantErr: domaincity.ErrInvalidUserID,
		},
		{
			name:   "user not found",
			userID: userID,
			cityID: cityID,
			setupMock: func(_ *mockRepository, ur *mockUserRepository) {
				ur.On("GetByID", mock.Anything, userID).Return(nil, domainuser.ErrNotFound).Once()
			},
			wantErr: domainuser.ErrNotFound,
		},
		{
			name:   "repo error propagates",
			userID: userID,
			cityID: cityID,
			setupMock: func(r *mockRepository, ur *mockUserRepository) {
				ur.On("GetByID", mock.Anything, userID).
					Return(&domainuser.User{ID: userID}, nil).Once()
				r.On("Delete", mock.Anything, userID, cityID).Return(repoErr).Once()
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, userRepo, svc := newService(t)
			if tt.setupMock != nil {
				tt.setupMock(repo, userRepo)
			}

			err := svc.Delete(context.Background(), tt.userID, tt.cityID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
