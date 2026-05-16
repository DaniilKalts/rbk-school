package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/city"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/history"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"

	cityrepo "github.com/DaniilKalts/rbk-school/6-week/internal/repository/city"
	userrepo "github.com/DaniilKalts/rbk-school/6-week/internal/repository/user"
	weatherrepo "github.com/DaniilKalts/rbk-school/6-week/internal/repository/weather"
)

type UserRepository interface {
	Create(ctx context.Context, u user.User, password user.Password) (*user.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	GetCredentialsByEmail(ctx context.Context, email string) (*user.User, user.Password, error)
	List(ctx context.Context) ([]user.User, error)
	Update(ctx context.Context, u user.User) (*user.User, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type CityRepository interface {
	Create(ctx context.Context, c city.City) (*city.City, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]city.City, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

type WeatherRepository interface {
	CreateHistory(ctx context.Context, history history.History) (*history.History, error)
	ListHistory(ctx context.Context, userID uuid.UUID, city string, limit int, offset int) ([]history.History, error)
}

type Repositories struct {
	User    UserRepository
	City    CityRepository
	Weather WeatherRepository
}

func NewRepositories(db sqlc.DBTX) *Repositories {
	return &Repositories{
		User:    userrepo.NewRepository(db),
		City:    cityrepo.NewRepository(db),
		Weather: weatherrepo.NewRepository(db),
	}
}
