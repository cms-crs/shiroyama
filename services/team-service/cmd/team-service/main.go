package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"taskservice/internal/app"
	"taskservice/internal/config"
	"taskservice/internal/infrastructure/database"
	"time"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Starting gRPC server")
		if err := application.GRPCServer.Run(); err != nil {
			log.Error("gRPC server failed", "error", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Starting Kafka consumer")
		if err := application.StartKafkaConsumer(ctx); err != nil {
			log.Error("Kafka consumer failed", "error", err)
			cancel()
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sign := <-signalChan:
		log.Info("Received shutdown signal", slog.String("signal", sign.String()))
	case <-ctx.Done():
		log.Info("Context cancelled, shutting down")
	}

	log.Info("Starting graceful shutdown...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)

		log.Info("Stopping gRPC server...")
		application.GRPCServer.Stop()

		log.Info("Closing application resources...")
		if err := application.Close(); err != nil {
			log.Error("Error closing application resources", "error", err)
		}

		log.Info("Waiting for goroutines to finish...")
		wg.Wait()
	}()

	select {
	case <-shutdownDone:
		log.Info("Graceful shutdown completed")
	case <-shutdownCtx.Done():
		log.Warn("Shutdown timeout exceeded, forcing exit")
	}

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
