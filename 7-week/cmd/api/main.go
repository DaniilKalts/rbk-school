package main

import (
	"flag"
	"log"

	"go.uber.org/zap"

	"github.com/DaniilKalts/rbk-school/7-week/internal/app"
	"github.com/DaniilKalts/rbk-school/7-week/internal/config"
	"github.com/DaniilKalts/rbk-school/7-week/pkg/logger"
)

func main() {
	configPath := flag.String("config-path", ".env", "путь к файлу конфигурации")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("не удалось загрузить конфигурацию: %v", err)
	}

	zapLogger, err := logger.New(cfg.Logger.Level, cfg.Logger.Format)
	if err != nil {
		log.Fatalf("не удалось создать логгер: %v", err)
	}
	defer func() {
		_ = zapLogger.Sync()
	}()

	a, err := app.NewApp(&cfg, zapLogger)
	if err != nil {
		zapLogger.Fatal("не удалось собрать приложение", zap.Error(err))
	}

	if err := a.Run(); err != nil {
		zapLogger.Fatal("не удалось запустить приложение", zap.Error(err))
	}
}
