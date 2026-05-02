package main

import (
	"log/slog"
	"os"

	"github.com/Romasmi/s-shop-microservices/internal/app"
	"github.com/Romasmi/s-shop-microservices/internal/config"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig(".")
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		slog.Error("Failed to initialize app", "error", err)
		os.Exit(1)
	}
	defer application.Close()

	cliApp := app.NewCli(application)
	if err := cliApp.Run(); err != nil {
		slog.Error("CLI error", "error", err)
		os.Exit(1)
	}
}
