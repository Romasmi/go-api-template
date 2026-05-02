package main

import (
	"log"

	"github.com/Romasmi/s-shop-microservices/internal/app"
	"github.com/Romasmi/s-shop-microservices/internal/config"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer application.Close()

	cliApp := app.NewCli(application)
	if err := cliApp.Run(); err != nil {
		log.Fatalf("CLI error: %v", err)
	}
}
