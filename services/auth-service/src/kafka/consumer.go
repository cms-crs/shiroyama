package kafka

import (
	"authservice/src/config"
	authRepo "authservice/src/repository"
	"context"
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"sync"
	"time"
)

type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	handler       *AuthConsumer
	logger        *slog.Logger
}

func NewConsumer(ctx context.Context, cfg *config.Config, logger *slog.Logger, rdb *redis.Client, db *gorm.DB) (*Consumer, error) {
	authRepository, err := authRepo.NewAuthRepository(db, rdb, cfg, logger)
	if err != nil {
		return nil, err
	}

	producer, err := NewKafkaProducer(cfg.Kafka.Brokers)
	if err != nil {
		return nil, err
	}

	consumerHandler := NewAuthConsumer(authRepository, producer, logger)

	configSarama := sarama.NewConfig()
	configSarama.Version = sarama.V2_8_0_0
	configSarama.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	configSarama.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, "user-service-group", configSarama)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumerGroup: consumerGroup,
		handler:       consumerHandler,
		logger:        logger,
	}, nil
}

func (c *Consumer) Start(ctx context.Context, topics []string) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if err := c.consumerGroup.Consume(ctx, topics, c.handler); err != nil {
				c.logger.Error("Error from consumer", "error", err)
				time.Sleep(time.Second)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-ctx.Done()

	wg.Wait()

	if err := c.consumerGroup.Close(); err != nil {
		c.logger.Error("Error closing consumer group", "error", err)
	}

	c.logger.Info("Kafka consumer stopped")
	return nil
}
