package app

import (
	"gorm.io/gorm"
	"log/slog"
	grpcapp "taskservice/internal/app/grpc"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	logger *slog.Logger,
	grpcPort int,
	db *gorm.DB,
) *App {
	return &App{
		GRPCServer: grpcapp.New(logger, grpcPort, db),
	}
}
