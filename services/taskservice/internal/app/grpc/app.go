package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log/slog"
	"net"
	postgresrepo "taskservice/internal/adapter/postgresrepo/task"
	"taskservice/internal/handler/grpc/taskservice"
	service "taskservice/internal/service/task"
)

type App struct {
	log  *slog.Logger
	gRPC *grpc.Server
	port int
	db   *gorm.DB
}

func New(log *slog.Logger, port int, db *gorm.DB) *App {
	gRPCServer := grpc.NewServer()
	taskRepo := postgresrepo.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)

	taskservicegrpc.Register(gRPCServer, taskService, log)

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
