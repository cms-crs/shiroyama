package app

import (
	"authservice/src/config"
	"fmt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log/slog"
	"net"
)

type App struct {
	db   *gorm.DB
	log  *slog.Logger
	conf *config.Config
	gRPC *grpc.Server
}

func New(db *gorm.DB, cfg *config.Config, logger *slog.Logger) *App {
	gRPCServer := grpc.NewServer()
	return &App{
		db:   db,
		log:  logger,
		conf: cfg,
		gRPC: gRPCServer,
	}
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err.Error())
	}
}

func (app *App) Run() error {
	const op = "app.Run"

	log := app.log.With(
		slog.String("op", op),
		slog.Int("port", app.conf.Grpc.Port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", app.conf.Grpc.Port))
	if err != nil {
		return err
	}

	log.Info("Starting gRPC server", listener.Addr().String())

	if err := app.gRPC.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) Stop() {
	const op = "app.Stop"

	log := app.log.With(
		slog.String("op", op),
	)

	log.Info("Stopping gRPC server", slog.Int("port", app.conf.Grpc.Port))

	app.gRPC.GracefulStop()
}
