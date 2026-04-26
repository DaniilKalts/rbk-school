package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/DaniilKalts/rbk-school/3-week/internal/adapters/database/postgres"
	docshttp "github.com/DaniilKalts/rbk-school/3-week/internal/adapters/transport/http/docs"
	userhttp "github.com/DaniilKalts/rbk-school/3-week/internal/adapters/transport/http/v1/user"
	"github.com/DaniilKalts/rbk-school/3-week/internal/config"
	userrepo "github.com/DaniilKalts/rbk-school/3-week/internal/repository/user"
	userservice "github.com/DaniilKalts/rbk-school/3-week/internal/service/user"
)

type Container struct {
	Config *config.Config

	DB     *pgxpool.Pool
	Router *http.ServeMux
}

func NewContainer(cfg *config.Config) (*Container, error) {
	if cfg == nil {
		return nil, fmt.Errorf("container: config is nil")
	}

	db, err := postgres.New(context.Background(), &cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("container: create postgres client: %w", err)
	}

	userRepository := userrepo.New(db)
	userService := userservice.New(userRepository)

	router := newRouter(userService)

	return &Container{
		Config: cfg,
		DB:     db,
		Router: router,
	}, nil
}

func (c *Container) Close() {
	if c == nil || c.DB == nil {
		return
	}

	c.DB.Close()
}

func newRouter(userService userhttp.Service) *http.ServeMux {
	mux := http.NewServeMux()
	docshttp.RegisterRoutes(mux)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	userhttp.RegisterRoutes(mux, userService)

	return mux
}
