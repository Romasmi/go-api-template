package app

import (
	"context"
	"fmt"

	"github.com/Romasmi/s-shop-microservices/internal/infrastructure/db/postgres"
	infrakafka "github.com/Romasmi/s-shop-microservices/internal/infrastructure/kafka"
	kafkaint "github.com/Romasmi/s-shop-microservices/internal/interface/kafka"
	"github.com/Romasmi/s-shop-microservices/internal/usecase"
	useruc "github.com/Romasmi/s-shop-microservices/internal/usecase/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string

	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string
}

type App struct {
	Cfg       Config
	Pool      *pgxpool.Pool
	Producer  useruc.EventProducer
	Handlers  map[usecase.UseCaseID]usecase.Handler
	Consumers []kafkaint.Consumer
}

func NewApp(cfg Config) (*App, error) {
	ctx := context.Background()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	producer := infrakafka.NewUserProducer(cfg.KafkaBrokers, cfg.KafkaTopic)

	app := &App{
		Cfg:       cfg,
		Pool:      pool,
		Producer:  producer,
		Handlers:  make(map[usecase.UseCaseID]usecase.Handler),
		Consumers: make([]kafkaint.Consumer, 0),
	}

	app.registerHandlers()
	app.registerConsumers()

	return app, nil
}

func (a *App) registerHandlers() {
	repo := postgres.NewUserRepository(a.Pool)

	a.Handlers[usecase.UseCaseCreateUser] = usecase.NewHandler(useruc.NewCreateUserUseCase(repo, a.Producer))
	a.Handlers[usecase.UseCaseGetUser] = usecase.NewHandler(useruc.NewGetUserUseCase(repo))
}

func (a *App) registerConsumers() {
	userConsumer := kafkaint.NewUserConsumer(a.Cfg.KafkaBrokers, a.Cfg.KafkaTopic, a.Cfg.KafkaGroupID)
	a.Consumers = append(a.Consumers, userConsumer)
}

func (a *App) Close() {
	if a.Pool != nil {
		a.Pool.Close()
	}
	if a.Producer != nil {
		// We know it's kafka.userProducer which has Close()
		// But it's better to cast or use interface
		if closer, ok := a.Producer.(interface{ Close() error }); ok {
			closer.Close()
		}
	}
	for _, c := range a.Consumers {
		c.Close()
	}
}

func (a *App) GetHandler(id usecase.UseCaseID) usecase.Handler {
	return a.Handlers[id]
}

func (a *App) GetConfig() Config {
	return a.Cfg
}

func (a *App) Ping(ctx context.Context) error {
	if a.Pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	return a.Pool.Ping(ctx)
}
