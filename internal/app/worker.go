package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	kafkaint "github.com/Romasmi/s-shop-microservices/internal/interface/kafka"
)

type Worker struct {
	*App
}

func NewWorker(app *App) *Worker {
	return &Worker{App: app}
}

func (w *Worker) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	consumer := kafkaint.NewUserConsumer(w.Cfg.Kafka.Brokers, w.Cfg.Kafka.Topic, w.Cfg.Kafka.GroupID)
	defer consumer.Close()

	go consumer.Start(ctx)

	<-ctx.Done()
	slog.Info("Shutting down worker...")
	return nil
}
