package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"saga-orchestrator/internal/storage"
	"sync"
	"syscall"
	"time"

	"saga-orchestrator/internal/config"
	"saga-orchestrator/internal/kafka"
	"saga-orchestrator/internal/saga"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.LogLevel)

	//metricsServer := metrics.NewServer(cfg.MetricsPort)
	//go func() {
	//	if err := metricsServer.Start(); err != nil {
	//		log.Error("Failed to start metrics server: %v", err)
	//	}
	//}()

	sagaStorage, err := storage.NewRedisStorage(cfg.Redis)
	if err != nil {
		log.Error("Saga Orchestrator started successfully", "port", cfg.HTTPPort)
	}
	defer sagaStorage.Close()

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Producer)
	if err != nil {
		log.Error("Failed to create Kafka producer", err)
	}
	defer producer.Close()

	orchestrator := saga.NewOrchestrator(
		producer,
		sagaStorage,
		log,
		cfg.Saga,
	)

	consumer, err := kafka.NewConsumerGroup(
		cfg.Kafka.Brokers,
		cfg.Kafka.Consumer.GroupID,
		cfg.Kafka.Consumer,
		orchestrator,
		log,
	)
	if err != nil {
		log.Error("Failed to create Kafka consumer", err)
	}
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := consumer.Start(ctx); err != nil {
			log.Error("Consumer error", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		orchestrator.StartTimeoutMonitor(ctx)
	}()

	healthServer := NewHealthServer(cfg.HTTPPort, log)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := healthServer.Start(); err != nil {
			log.Error("Health server error", err)
		}
	}()

	log.Error("Saga Orchestrator started successfully on port", cfg.HTTPPort)
	log.Error("Metrics server started on port", cfg.MetricsPort)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	log.Info("Received shutdown signal, starting graceful shutdown...")

	cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info("Graceful shutdown completed")
	case <-time.After(30 * time.Second):
		log.Warn("Shutdown timeout exceeded, forcing exit")
	}

	healthServer.Stop()

	log.Info("Saga Orchestrator stopped")
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

type HealthServer struct {
	port   int
	log    *slog.Logger
	server *http.Server
}

func NewHealthServer(port int, log *slog.Logger) *HealthServer {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	hs := &HealthServer{
		port:   port,
		log:    log,
		server: server,
	}

	mux.HandleFunc("/health", hs.healthHandler)
	mux.HandleFunc("/ready", hs.readyHandler)

	return hs
}

func (hs *HealthServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (hs *HealthServer) readyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}

func (hs *HealthServer) Start() error {
	hs.log.Error("Starting health server on port", hs.port)
	return hs.server.ListenAndServe()
}

func (hs *HealthServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := hs.server.Shutdown(ctx); err != nil {
		hs.log.Error("Error stopping health server", err)
	}
}
