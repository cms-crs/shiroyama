package app

import (
	"log/slog"
	grpcapp "taskservice/internal/app/grpc"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	logger *slog.Logger,
	grpcPort int,
) *App {

	return &App{
		GRPCServer: grpcapp.New(logger, grpcPort),
	}
}
