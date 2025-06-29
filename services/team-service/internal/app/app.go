package app

import (
	"context"
	"database/sql"
	"log/slog"
	grpcapp "taskservice/internal/app/grpc"
	"taskservice/internal/clients"
	"taskservice/internal/config"
	"taskservice/internal/kafka"
	teamRepo "taskservice/internal/repository/team"
	teamService "taskservice/internal/service/team"
)

type App struct {
	GRPCServer    *grpcapp.App
	KafkaConsumer *kafka.Consumer
	KafkaProducer *kafka.Producer
	logger        *slog.Logger
}

func New(
	logger *slog.Logger,
	grpcPort int,
	db *sql.DB,
	cfg *config.Config,
) *App {
	kafkaProducer, err := kafka.NewKafkaProducer(cfg.Kafka.Brokers)
	if err != nil {
		logger.Error("Failed to create Kafka producer", "error", err)
		panic(err)
	}

	userClient, err := clients.NewUserClient(cfg.UserService.Address)
	if err != nil {
		logger.Error("Failed to create user client", "error", err)
		userClient = nil
	}

	teamRepository := teamRepo.NewTeamRepository(logger, db)
	teamSvc := teamService.NewTeamService(logger, teamRepository, userClient)

	kafkaConsumer, err := kafka.NewConsumer(kafka.ConsumerConfig{
		Brokers:     cfg.Kafka.Brokers,
		GroupID:     "team-service-group",
		TeamService: teamSvc,
		Logger:      logger,
	})
	if err != nil {
		logger.Error("Failed to create Kafka consumer", "error", err)
		kafkaProducer.Close()
		panic(err)
	}

	grpcApp := grpcapp.New(logger, grpcPort, db, cfg, kafkaProducer)

	return &App{
		GRPCServer:    grpcApp,
		KafkaConsumer: kafkaConsumer,
		KafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (app *App) StartKafkaConsumer(ctx context.Context) error {
	topics := []string{"team-service-commands"}

	app.logger.Info("Starting Kafka consumer", "topics", topics)
	return app.KafkaConsumer.Start(ctx, topics)
}

func (app *App) Close() error {
	app.logger.Info("Closing application resources")

	var lastErr error

	if err := app.KafkaProducer.Close(); err != nil {
		app.logger.Error("Failed to close Kafka producer", "error", err)
		lastErr = err
	}

	if err := app.KafkaConsumer.Close(); err != nil {
		app.logger.Error("Failed to close Kafka consumer", "error", err)
		lastErr = err
	}

	return lastErr
}
