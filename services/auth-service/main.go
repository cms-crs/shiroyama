package main

import (
	"authservice/src/app"
	"authservice/src/config"
	"authservice/src/database"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	db := database.MustConnect(cfg)

	log.Info("Starting application", slog.Any("config", cfg))

	application := app.New(db, cfg, log)
	go application.MustRun()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	sign := <-signalChan

	log.Info("Stopping application", slog.String("signal", sign.String()))

	application.Stop()

	log.Info("Application stopped")
}
