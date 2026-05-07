package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/DaniilKalts/rbk-school/5-week/internal/config"
)

type App struct {
	container *Container
	server    *http.Server
}

func NewApp(cfg *config.Config) *App {
	container, err := NewContainer(cfg)
	if err != nil {
		log.Fatalf("не удалось собрать DI-контейнер приложения: %v", err)
	}

	a := &App{container: container}
	a.initDeps()

	return a
}

func (a *App) initDeps() {
	inits := []func(){
		a.initHTTPServer,
	}

	for _, fn := range inits {
		fn()
	}
}

func (a *App) initHTTPServer() {
	cfg := a.container.Config()
	a.server = &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      a.container.Router(),
		ReadTimeout:  cfg.Server.HTTPTimeout,
		WriteTimeout: cfg.Server.HTTPTimeout,
		IdleTimeout:  cfg.Server.HTTPTimeout,
	}
}

func (a *App) Run() error {
	if a == nil || a.server == nil {
		return fmt.Errorf("app: сервер не инициализирован")
	}
	if a.container != nil {
		defer a.container.Close()
	}

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	log.Printf("server is running on %s", a.server.Addr)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("app: ошибка запуска сервера: %w", err)
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-sigCtx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.server.ReadTimeout)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("app: ошибка остановки сервера: %w", err)
	}

	if err := <-errCh; err != nil {
		return err
	}

	log.Print("server is stopped")

	return nil
}
