package main

import (
	"flag"
	"log"

	"github.com/DaniilKalts/rbk-school/2-week/internal/app"
	"github.com/DaniilKalts/rbk-school/2-week/internal/config"
)

func main() {
	configPath := flag.String("config-path", ".env", "path to config file")

	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	container, err := app.NewContainer(cfg)
	if err != nil {
		log.Fatalf("failed to build app container: %v", err)
	}

	application, err := app.New(container)
	if err != nil {
		log.Fatalf("failed to build app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
