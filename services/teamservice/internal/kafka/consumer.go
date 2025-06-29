package kafka

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	handler       *TeamConsumer
	logger        *slog.Logger
}

type ConsumerConfig struct {
	Brokers     []string
	GroupID     string
	TeamService TeamService
	Logger      *slog.Logger
}

func NewConsumer(cfg ConsumerConfig) (*Consumer, error) {
	producer, err := NewKafkaProducer(cfg.Brokers)
	if err != nil {
		return nil, err
	}

	consumerHandler := NewTeamConsumer(cfg.TeamService, producer, cfg.Logger)

	configSarama := sarama.NewConfig()
	configSarama.Version = sarama.V2_8_0_0
	configSarama.Consumer.Offsets.Initial = sarama.OffsetNewest
	configSarama.Consumer.MaxProcessingTime = 20 * time.Second
	configSarama.Consumer.Group.Heartbeat.Interval = 6 * time.Second
	configSarama.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, configSarama)
	if err != nil {
		producer.Close()
		return nil, err
	}

	return &Consumer{
		consumerGroup: consumerGroup,
		handler:       consumerHandler,
		logger:        cfg.Logger,
	}, nil
}

func (c *Consumer) Start(ctx context.Context, topics []string) error {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			if err := c.consumerGroup.Consume(ctx, topics, c.handler); err != nil {
				c.logger.Error("Error from consumer", "error", err)
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Second):
					continue
				}
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case err, ok := <-c.consumerGroup.Errors():
				if !ok {
					return
				}
				c.logger.Error("Consumer group error", "error", err)
			}
		}
	}()

	c.logger.Info("Kafka consumer started", "topics", topics)

	<-ctx.Done()
	c.logger.Info("Shutting down Kafka consumer...")

	wg.Wait()

	if err := c.consumerGroup.Close(); err != nil {
		c.logger.Error("Error closing consumer group", "error", err)
		return err
	}

	c.logger.Info("Kafka consumer stopped")
	return nil
}

func (c *Consumer) Close() error {
	return c.consumerGroup.Close()
}
