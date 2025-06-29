package main

import (
	"boardservice/internal/app"
	"boardservice/internal/config"
	"boardservice/internal/infrastructure/database"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	db := database.MustLoad(cfg)

	log.Info("Starting board service application", slog.Any("config", cfg))
	application := app.New(log, cfg.Grpc.Port, db, cfg)
	go application.GRPCServer.MustRun()

	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sign := <-signalChan:
		log.Info("Stopping board service application", slog.String("signal", sign.String()))
	case <-ctx.Done():
	}

	cancel()
	application.GRPCServer.Stop()

	log.Info("Board service application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
