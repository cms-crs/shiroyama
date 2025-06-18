package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"userservice/internal/app"
	"userservice/internal/config"
	"userservice/internal/infrastructure/database"
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

	log.Info("Starting application", slog.Any("config", cfg))
	application := app.New(log, cfg.Grpc.Port, db, cfg)
	go application.GRPCServer.MustRun()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	sign := <-signalChan

	log.Info("Stopping application", slog.String("signal", sign.String()))

	application.GRPCServer.Stop()

	log.Info("Application stopped")

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
