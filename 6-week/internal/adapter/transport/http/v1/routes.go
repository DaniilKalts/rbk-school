package v1

import (
	"github.com/go-chi/chi/v5"

	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/middleware"
	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/v1/auth"
	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/v1/city"
	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/v1/user"
	"github.com/DaniilKalts/rbk-school/6-week/internal/adapter/transport/http/v1/weather"
	"github.com/DaniilKalts/rbk-school/6-week/pkg/jwt"
)

type Dependencies struct {
	AuthService    auth.Service
	CityService    city.Service
	WeatherService weather.Service
	UserService    user.Service
	TokenManager   *jwt.Manager
}

func RegisterRoutes(r chi.Router, deps Dependencies) {
	r.Route("/api/v1", func(r chi.Router) {
		auth.RegisterRoutes(r, deps.AuthService)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(deps.TokenManager))

			city.RegisterRoutes(r, deps.CityService)
			weather.RegisterRoutes(r, deps.WeatherService)
			user.RegisterCurrentUserRoutes(r, deps.UserService, deps.TokenManager)

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("admin"))
				user.RegisterAdminRoutes(r, deps.UserService, deps.TokenManager)
			})
		})
	})
}
