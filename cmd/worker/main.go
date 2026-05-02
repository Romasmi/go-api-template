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

	worker := app.NewWorker(application)
	if err := worker.Run(); err != nil {
		log.Fatalf("Worker error: %v", err)
	}
}
