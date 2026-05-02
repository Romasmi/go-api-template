package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/Romasmi/s-shop-microservices/internal/domain/user"
	"github.com/segmentio/kafka-go"
)

type UserConsumer struct {
	reader *kafka.Reader
}

func NewUserConsumer(brokers []string, topic, groupID string) *UserConsumer {
	return &UserConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	}
}

func (c *UserConsumer) Start(ctx context.Context) {
	slog.Info("Starting UserConsumer", "topic", c.reader.Config().Topic)
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("Error reading message", "error", err)
			continue
		}

		var u user.User
		if err := json.Unmarshal(m.Value, &u); err != nil {
			slog.Error("Error unmarshaling message", "error", err)
			continue
		}

		slog.Info("Consumed UserCreated event", "user_id", u.ID)
	}
}

func (c *UserConsumer) Close() error {
	return c.reader.Close()
}
