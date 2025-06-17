package main

import (
	"authservice/src/app"
	"authservice/src/config"
	"authservice/src/database/postgres"
	"authservice/src/database/redis"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// load env to get secret key
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	
	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	db := postgres.MustConnect(cfg)
	rdb := redis.MustConnect(cfg)

	application := app.New(db, rdb, cfg, log)
	go application.MustRun()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	sign := <-signalChan

	log.Info("Stopping application", slog.String("signal", sign.String()))

	application.Stop()

	log.Info("Application stopped")
}
