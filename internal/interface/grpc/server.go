package grpc

import (
	api "github.com/Romasmi/s-shop-microservices/internal/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer(provider UseCaseProvider) *grpc.Server {
	grpcServer := grpc.NewServer()

	userHandler := NewUserHandler(provider)
	api.RegisterUserServiceServer(grpcServer, userHandler)

	reflection.Register(grpcServer)

	return grpcServer
}
