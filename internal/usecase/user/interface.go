package user

import (
	"context"

	"github.com/Romasmi/s-shop-microservices/internal/domain/user"
)

type EventProducer interface {
	UserCreated(ctx context.Context, u *user.User) error
}
