package app

import (
	"database/sql"
	"log/slog"
	grpcapp "userservice/internal/app/grpc"
	"userservice/internal/config"
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
