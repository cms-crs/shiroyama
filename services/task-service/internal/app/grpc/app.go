package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log/slog"
	"net"
	"taskservice/internal/clients"
	"taskservice/internal/config"
	handler "taskservice/internal/handler/grpc/taskservice"
	postgresrepo "taskservice/internal/repository"
	service "taskservice/internal/service/task"
)

type App struct {
	log  *slog.Logger
	gRPC *grpc.Server
	port int
	db   *gorm.DB
}

func New(log *slog.Logger, port int, db *gorm.DB, cfg *config.Config) *App {
	gRPCServer := grpc.NewServer()
	serviceClients, err := clients.NewServiceClients(&clients.ClientConfig{
		UserServiceAddr:  cfg.Clients.BoardServiceAddr,
		BoardServiceAddr: cfg.Clients.TeamServiceAddr,
		DialTimeout:      cfg.Clients.DialTimeout,
	})

	if err != nil {
		panic(err)
	}

	taskRepo := postgresrepo.NewTaskRepository(db)
	taskService := service.NewTaskService(log, taskRepo, serviceClients.UserClient, serviceClients.BoardClient)

	handler.Register(gRPCServer, log, taskService)

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
