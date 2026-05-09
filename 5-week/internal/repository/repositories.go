package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/5-week/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/city"
	"github.com/DaniilKalts/rbk-school/5-week/internal/domain/history"
	"github.com/DaniilKalts/rbk-school/5-week/internal/repository/user"
	"github.com/DaniilKalts/rbk-school/5-week/internal/repository/weather"

	cityrepo "github.com/DaniilKalts/rbk-school/5-week/internal/repository/city"
)

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
	User    *user.Repository
	City    CityRepository
	Weather WeatherRepository
}

func NewRepositories(db sqlc.DBTX) *Repositories {
	return &Repositories{
		User:    user.NewRepository(db),
		City:    cityrepo.NewRepository(db),
		Weather: weather.NewRepository(db),
	}
}
