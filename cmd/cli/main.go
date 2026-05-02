package main

import (
	"log/slog"
	"os"

	"github.com/Romasmi/s-shop-microservices/cmd/api/config"
	"github.com/Romasmi/s-shop-microservices/internal/app"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.LoadConfig(".")
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	appCfg := app.Config{
		DBUser:       cfg.Db.User,
		DBPassword:   cfg.Db.Password,
		DBHost:       cfg.Db.Host,
		DBPort:       cfg.Db.Port,
		DBName:       cfg.Db.Name,
		KafkaBrokers: cfg.Kafka.Brokers,
		KafkaTopic:   cfg.Kafka.Topic,
		KafkaGroupID: cfg.Kafka.GroupID,
	}

	application, err := app.NewApp(appCfg)
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
