package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Romasmi/s-shop-microservices/internal/domain/user"
	"github.com/segmentio/kafka-go"
)

type UserProducer struct {
	Writer *kafka.Writer
}

func NewUserProducer(brokers []string, topic string) *UserProducer {
	return &UserProducer{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *UserProducer) UserCreated(ctx context.Context, u *user.User) error {
	payload, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	err = p.Writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(u.ID.String()),
		Value: payload,
	})
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (p *UserProducer) Close() error {
	return p.Writer.Close()
}
