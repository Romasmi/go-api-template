package grpc

import (
	"context"

	api "github.com/Romasmi/s-shop-microservices/internal/api"
	"github.com/Romasmi/s-shop-microservices/internal/domain/user"
	"github.com/Romasmi/s-shop-microservices/internal/usecase"
	useruc "github.com/Romasmi/s-shop-microservices/internal/usecase/user"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UseCaseProvider interface {
	GetHandler(id usecase.UseCaseID) usecase.Handler
}

type UserHandler struct {
	api.UnimplementedUserServiceServer
	provider UseCaseProvider
}

func NewUserHandler(provider UseCaseProvider) *UserHandler {
	return &UserHandler{
		provider: provider,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	handler := h.provider.GetHandler(usecase.UseCaseCreateUser)
	input := useruc.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	}

	resp, err := handler.Do(ctx, input)
	if err != nil {
		return nil, err
	}

	u := resp.(*user.User)

	return &api.CreateUserResponse{
		User: &api.User{
			Id:        u.ID.String(),
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: timestamppb.New(u.CreatedAt),
		},
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *api.GetUserRequest) (*api.GetUserResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	handler := h.provider.GetHandler(usecase.UseCaseGetUser)
	resp, err := handler.Do(ctx, id)
	if err != nil {
		return nil, err
	}

	u := resp.(*user.User)

	return &api.GetUserResponse{
		User: &api.User{
			Id:        u.ID.String(),
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: timestamppb.New(u.CreatedAt),
		},
	}, nil
}
