package main

import (
	"authservice/src/app"
	"authservice/src/config"
	"authservice/src/database/postgres"
	"authservice/src/database/redis"
	"authservice/src/kafka"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// load env to get secret key
	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	db := postgres.MustConnect(cfg)
	rdb := redis.MustConnect(cfg)

	application := app.New(db, rdb, cfg, log)
	go application.MustRun()

	ctx, cancel := context.WithCancel(context.Background())

	consumer, err := kafka.NewConsumer(ctx, cfg, log, rdb, db)
	if err != nil {
		log.Error("Failed to create Kafka consumer", "error", err)
		return
	}

	go func() {
		if err := consumer.Start(ctx, []string{"user-deletion-saga"}); err != nil {
			log.Error("Kafka consumer stopped with error", "error", err)
			cancel()
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sign := <-signalChan:
		log.Info("Stopping application", slog.String("signal", sign.String()))
	case <-ctx.Done():
	}

	cancel()

	application.Stop()

	log.Info("Application stopped")
}
