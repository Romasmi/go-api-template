package app

import (
	"context"
	"fmt"

	"github.com/Romasmi/s-shop-microservices/internal/config"
	"github.com/Romasmi/s-shop-microservices/internal/infrastructure/db/postgres"
	"github.com/Romasmi/s-shop-microservices/internal/infrastructure/kafka"
	"github.com/Romasmi/s-shop-microservices/internal/usecase"
	useruc "github.com/Romasmi/s-shop-microservices/internal/usecase/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Cfg      *config.Config
	Pool     *pgxpool.Pool
	Producer useruc.EventProducer
	Handlers map[usecase.UseCaseID]usecase.Handler
}

func NewApp(cfg *config.Config) (*App, error) {
	ctx := context.Background()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.Name)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	producer := kafka.NewUserProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)

	app := &App{
		Cfg:      cfg,
		Pool:     pool,
		Producer: producer,
		Handlers: make(map[usecase.UseCaseID]usecase.Handler),
	}

	app.registerHandlers()

	return app, nil
}

func (a *App) registerHandlers() {
	repo := postgres.NewUserRepository(a.Pool)

	a.Handlers[usecase.UseCaseCreateUser] = usecase.NewHandler(useruc.NewCreateUserUseCase(repo, a.Producer))
	a.Handlers[usecase.UseCaseGetUser] = usecase.NewHandler(useruc.NewGetUserUseCase(repo))
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
}

func (a *App) GetHandler(id usecase.UseCaseID) usecase.Handler {
	return a.Handlers[id]
}

func (a *App) GetConfig() *config.Config {
	return a.Cfg
}

func (a *App) Ping(ctx context.Context) error {
	if a.Pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	return a.Pool.Ping(ctx)
}
