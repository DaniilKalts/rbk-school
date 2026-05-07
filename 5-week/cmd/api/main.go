package main

import (
	"flag"
	"log"

	"github.com/DaniilKalts/rbk-school/5-week/internal/app"
	"github.com/DaniilKalts/rbk-school/5-week/internal/config"
)

func main() {
	configPath := flag.String("config-path", ".env", "path to config file")

	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("не удалось загрузить конфигурацию: %v", err)
	}

	a := app.NewApp(cfg)

	if err := a.Run(); err != nil {
		log.Fatalf("не удалось запустить приложение: %v", err)
	}
}
