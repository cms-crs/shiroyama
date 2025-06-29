package app

import (
	grpcapp "boardservice/internal/app/grpc"
	"boardservice/internal/config"
	"database/sql"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	logger *slog.Logger,
	grpcPort int,
	db *sql.DB,
	cfg *config.Config,
) *App {

	return &App{
		GRPCServer: grpcapp.New(logger, grpcPort, db, cfg),
	}
}
