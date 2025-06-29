package grpcapp

import (
	"boardservice/internal/clients"
	"boardservice/internal/config"
	"boardservice/internal/handler"
	boardRepo "boardservice/internal/repository/board"
	listRepo "boardservice/internal/repository/list"
	boardService "boardservice/internal/service/board"
	listService "boardservice/internal/service/list"
	"database/sql"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log  *slog.Logger
	gRPC *grpc.Server
	port int
	db   *sql.DB
}

func New(log *slog.Logger, port int, db *sql.DB, cfg *config.Config) *App {
	gRPCServer := grpc.NewServer()

	boardRepository := boardRepo.NewRepository(db)
	listRepository := listRepo.NewRepository(db)
	serviceClients, err := clients.NewServiceClients(&clients.ClientConfig{
		UserServiceAddr: cfg.Clients.UserServiceAddr,
		TeamServiceAddr: cfg.Clients.TeamServiceAddr,
		DialTimeout:     cfg.Clients.ClientTimeout,
	})

	if err != nil {
		log.Error("failed to initialize service clients", "error", err)
	}
	defer serviceClients.Close()

	boardSvc := boardService.NewBoardService(log, boardRepository, serviceClients.UserClient, serviceClients.TeamClient)
	listSvc := listService.NewListService(log, listRepository, boardRepository)

	handler.Register(gRPCServer, log, boardSvc, listSvc)

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

	log.Info("Starting gRPC server", "address", listener.Addr().String())

	if err := app.gRPC.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *App) Stop() {
	const op = "grpcapp.Stop"

	log := app.log.With(
		slog.String("op", op),
	)

	log.Info("Stopping gRPC server", slog.Int("port", app.port))

	app.gRPC.GracefulStop()
}
