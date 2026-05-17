package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/client"
	"github.com/DaniilKalts/rbk-school/6-week/internal/cache"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/city"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/history"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/user"
	"github.com/DaniilKalts/rbk-school/6-week/internal/domain/weather"
	"github.com/DaniilKalts/rbk-school/6-week/internal/repository"
	"github.com/DaniilKalts/rbk-school/6-week/internal/service/auth"

	servicecity "github.com/DaniilKalts/rbk-school/6-week/internal/service/city"
	serviceuser "github.com/DaniilKalts/rbk-school/6-week/internal/service/user"
	serviceweather "github.com/DaniilKalts/rbk-school/6-week/internal/service/weather"
)

type AuthService interface {
	Register(ctx context.Context, input auth.RegisterInput) (*auth.Token, error)
	Login(ctx context.Context, input auth.LoginInput) (*auth.Token, error)
	Logout(ctx context.Context, accessToken string) error
}

type UserService interface {
	Create(ctx context.Context, input serviceuser.CreateInput) (*user.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*user.User, error)
	List(ctx context.Context) ([]user.User, error)
	Update(ctx context.Context, id uuid.UUID, input serviceuser.UpdateInput) (*user.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CityService interface {
	Create(ctx context.Context, userID uuid.UUID, input servicecity.CreateInput) (*city.City, error)
	List(ctx context.Context, userID uuid.UUID) ([]city.City, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

type WeatherService interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]weather.Weather, error)
	GetHistory(ctx context.Context, userID uuid.UUID, city string, limit int, offset int) ([]history.History, error)
}

type Services struct {
	Auth    AuthService
	User    UserService
	City    CityService
	Weather WeatherService
}

func NewServices(repositories *repository.Repositories, caches *cache.Caches, clients *client.Clients, tokenManager auth.TokenManager) *Services {
	return &Services{
		Auth: auth.NewService(repositories.User, tokenManager),
		User: serviceuser.NewService(repositories.User),
		City: servicecity.NewService(repositories.City, repositories.User),
		Weather: serviceweather.NewService(
			repositories.User,
			repositories.City,
			repositories.Weather,
			clients.Geocoding,
			clients.OpenMeteo,
			caches.Weather,
		),
	}
}
