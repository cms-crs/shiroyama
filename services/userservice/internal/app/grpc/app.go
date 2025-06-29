package grpcapp

import (
	"database/sql"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"userservice/internal/config"
	"userservice/internal/handler"
	"userservice/internal/kafka"
	userRepo "userservice/internal/repository/user"
	userService "userservice/internal/service/user"
)

type App struct {
	log  *slog.Logger
	gRPC *grpc.Server
	port int
	db   *sql.DB
}

func New(log *slog.Logger, port int, db *sql.DB, cfg *config.Config) *App {
	gRPCServer := grpc.NewServer()

	userRepository := userRepo.NewUserRepository(log, db)
	service := userService.NewUserService(log, userRepository)
	kafkaProducer, err := kafka.NewKafkaProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Error("Failed to create kafka producer", "error", err)
		panic(err)
	}

	handler.Register(gRPCServer, log, service, kafkaProducer)

	return &App{
		log:  log,
		gRPC: gRPCServer,
		port: port,
		db:   db,
	}
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err.Error())
	}
}

func (app *App) Run() error {
	const op = "grpcapp.Run"

	log := app.log.With(
		slog.String("op", op),
		slog.Int("port", app.port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Starting gRPC server", listener.Addr().String())

	if err := app.gRPC.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) Stop() {
	const op = "grpc.Stop"

	log := app.log.With(
		slog.String("op", op),
	)

	log.Info("Stopping gRPC server", slog.Int("port", app.port))

	app.gRPC.GracefulStop()
}
