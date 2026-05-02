package user

import (
	"context"
	"time"

	"github.com/Romasmi/s-shop-microservices/internal/domain/user"
	"github.com/google/uuid"
)

type CreateUserInput struct {
	Name  string
	Email string
}

type CreateUserUseCase struct {
	repo          user.Repository
	eventProducer EventProducer
}

func NewCreateUserUseCase(repo user.Repository, eventProducer EventProducer) *CreateUserUseCase {
	return &CreateUserUseCase{
		repo:          repo,
		eventProducer: eventProducer,
	}
}

func (uc *CreateUserUseCase) Do(ctx context.Context, input CreateUserInput) (*user.User, error) {
	u := &user.User{
		ID:        uuid.New(),
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: time.Now(),
	}

	if err := uc.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	if uc.eventProducer != nil {
		_ = uc.eventProducer.UserCreated(ctx, u)
	}

	return u, nil
}
