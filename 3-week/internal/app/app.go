package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
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
	if a.container != nil {
		defer a.container.Close()
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	log.Printf("server is running on %s", a.server.Addr)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("app: run server: %w", err)
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.server.ReadTimeout)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("app: shutdown server: %w", err)
	}

	if err := <-errCh; err != nil {
		return err
	}

	log.Print("server is stopped")

	return nil
}
