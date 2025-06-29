package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"saga-orchestrator/internal/config"
	"saga-orchestrator/internal/events"
)

type EventHandler interface {
	HandleEvent(ctx context.Context, event events.Event) error
}

type ConsumerGroup struct {
	consumer sarama.ConsumerGroup
	handler  EventHandler
	log      *slog.Logger
	topics   []string
	cfg      config.ConsumerConfig
	wg       sync.WaitGroup
}

func NewConsumerGroup(
	brokers []string,
	groupID string,
	consumerConfig config.ConsumerConfig,
	handler EventHandler,
	log *slog.Logger,
) (*ConsumerGroup, error) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Group.Session.Timeout = consumerConfig.SessionTimeout
	cfg.Consumer.Group.Heartbeat.Interval = consumerConfig.HeartbeatInterval
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	if consumerConfig.AutoOffsetReset == "earliest" {
		cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	topics := []string{
		"user-deletion-saga",
		"auth-service-events",
		"team-service-events",
		"board-service-events",
		"task-service-events",
	}

	return &ConsumerGroup{
		consumer: consumer,
		handler:  handler,
		log:      log,
		topics:   topics,
		cfg:      consumerConfig,
	}, nil
}

func (cg *ConsumerGroup) Start(ctx context.Context) error {
	cg.wg.Add(1)
	go func() {
		defer cg.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-cg.consumer.Errors():
				if err != nil {
					cg.log.Error("Consumer error", err)
				}
			}
		}
	}()

	cg.wg.Add(1)
	go func() {
		defer cg.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := cg.consumer.Consume(ctx, cg.topics, cg); err != nil {
					cg.log.Info("Info from consumer", err)
					time.Sleep(time.Second)
				}
			}
		}
	}()

	return nil
}

func (cg *ConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	cg.log.Info("Consumer group session started")
	return nil
}

func (cg *ConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	cg.log.Info("Consumer group session ended")
	return nil
}

func (cg *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			if err := cg.processMessage(session.Context(), message); err != nil {
				cg.log.Error("Error processing message", err)
				continue
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func (cg *ConsumerGroup) processMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	start := time.Now()

	processingCtx, cancel := context.WithTimeout(ctx, cg.cfg.MaxProcessingTime)
	defer cancel()

	var event events.Event
	if err := json.Unmarshal(message.Value, &event); err != nil {
		//metrics.EventsProcessed.WithLabelValues("unknown", "unmarshal_error").Inc()
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	cg.log.Debug("Processing event",
		"type", event.Type,
		"saga", event.SagaID,
		"user", event.UserID,
	)

	if err := cg.handler.HandleEvent(processingCtx, event); err != nil {
		//metrics.EventsProcessed.WithLabelValues(string(event.Type), "handler_error").Inc()
		return fmt.Errorf("handler error: %w", err)
	}

	//metrics.EventsProcessed.WithLabelValues(string(event.Type), "success").Inc()

	duration := time.Since(start)
	cg.log.Debug("Event processed in", duration, ":", event.Type)

	return nil
}

func (cg *ConsumerGroup) Close() error {
	cg.log.Info("Closing consumer group...")

	if err := cg.consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %w", err)
	}

	cg.wg.Wait()
	cg.log.Info("Consumer group closed successfully")
	return nil
}
