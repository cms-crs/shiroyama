package grpcapp

import (
	"database/sql"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"userservice/internal/handler"
	teamRepo "userservice/internal/repository/team"
	teamService "userservice/internal/service/team"
)

type App struct {
	log  *slog.Logger
	gRPC *grpc.Server
	port int
	db   *sql.DB
}

func New(log *slog.Logger, port int, db *sql.DB) *App {
	gRPCServer := grpc.NewServer()

	teamRepository := teamRepo.NewTeamRepository(log, db)
	userService := teamService.NewTeamService(log, teamRepository)

	handler.Register(gRPCServer, log, userService)

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
