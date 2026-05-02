package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
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

	for _, consumer := range w.App.Consumers {
		go consumer.Start(ctx)
	}

	<-ctx.Done()
	slog.Info("Shutting down worker...")
	return nil
}
