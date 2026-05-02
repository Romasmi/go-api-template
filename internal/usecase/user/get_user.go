package user

import (
	"context"

	"github.com/Romasmi/s-shop-microservices/internal/domain/user"
	"github.com/google/uuid"
)

type GetUserUseCase struct {
	repo user.Repository
}

func NewGetUserUseCase(repo user.Repository) *GetUserUseCase {
	return &GetUserUseCase{
		repo: repo,
	}
}

func (uc *GetUserUseCase) Do(ctx context.Context, id uuid.UUID) (*user.User, error) {
	return uc.repo.GetByID(ctx, id)
}
