package v1

import (
	"github.com/go-chi/chi/v5"

	"github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/adapter/transport/http/v1/health"
	"github.com/DaniilKalts/rbk-school/7-week/gateway-service/internal/adapter/transport/http/v1/weather"
)

type Dependencies struct {
	WeatherService   weather.Service
	APIServiceClient health.APIServiceClient
}

func RegisterRoutes(r chi.Router, deps Dependencies) {
	r.Route("/api/v1", func(r chi.Router) {
		weather.RegisterRoutes(r, deps.WeatherService)
		health.RegisterRoutes(r, deps.APIServiceClient)
	})
}
