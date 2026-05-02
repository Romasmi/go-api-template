package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpcint "github.com/Romasmi/s-shop-microservices/internal/interface/grpc"
	httpint "github.com/Romasmi/s-shop-microservices/internal/interface/http"
)

type Api struct {
	*App
}

func NewApi(app *App) *Api {
	return &Api{App: app}
}

func (a *Api) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// gRPC Server
	grpcServer := grpcint.NewServer(a.App)

	grpcAddr := fmt.Sprintf(":%d", a.Cfg.Server.GRPCPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		log.Printf("Starting gRPC server on %s", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// gRPC Gateway
	gwServer, err := httpint.NewGatewayServer(grpcAddr, a.Cfg.Server.HTTPPort)
	if err != nil {
		return fmt.Errorf("failed to create gateway server: %w", err)
	}

	go func() {
		log.Printf("Starting HTTP gateway on %s", gwServer.Addr)
		if err := gwServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP gateway error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down API...")
	grpcServer.GracefulStop()
	if err := gwServer.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP gateway shutdown error: %v", err)
	}
	return nil
}
