package app

import (
	"fmt"
	"net/http"
)

type App struct {
	container *Container
	server    *http.Server
}

func New(container *Container) (*App, error) {
	if container == nil {
		return nil, fmt.Errorf("app: container is nil")
	}

	server := &http.Server{
		Addr:         container.Config.Server.Addr,
		Handler:      container.Router,
		ReadTimeout:  container.Config.Server.HTTPTimeout,
		WriteTimeout: container.Config.Server.HTTPTimeout,
		IdleTimeout:  container.Config.Server.HTTPTimeout,
	}

	return &App{
		container: container,
		server:    server,
	}, nil
}

func (a *App) Run() error {
	if a == nil || a.server == nil {
		return fmt.Errorf("app: server is not initialized")
	}

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("app: run server: %w", err)
	}

	return nil
}
