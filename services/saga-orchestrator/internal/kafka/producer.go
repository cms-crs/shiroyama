package kafka

import (
	"encoding/json"
	"fmt"
	"saga-orchestrator/internal/events"
	"time"

	"github.com/IBM/sarama"
	"saga-orchestrator/internal/config"
)

type Producer struct {
	producer sarama.SyncProducer
	cfg      config.ProducerConfig
}

func NewProducer(brokers []string, producerCfg config.ProducerConfig) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = producerCfg.RetryMax
	cfg.Producer.Flush.Frequency = producerCfg.FlushTimeout
	cfg.Producer.Flush.Messages = producerCfg.BatchSize

	switch producerCfg.Compression {
	case "gzip":
		cfg.Producer.Compression = sarama.CompressionGZIP
	case "snappy":
		cfg.Producer.Compression = sarama.CompressionSnappy
	case "lz4":
		cfg.Producer.Compression = sarama.CompressionLZ4
	case "zstd":
		cfg.Producer.Compression = sarama.CompressionZSTD
	default:
		cfg.Producer.Compression = sarama.CompressionSnappy
	}

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{
		producer: producer,
		cfg:      producerCfg,
	}, nil
}
func (p *Producer) PublishEvent(topic string, event events.Event) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		//metrics.EventsProcessed.WithLabelValues(string(event.Type), "marshal_error").Inc()
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	message := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(event.UserID),
		Value:     sarama.ByteEncoder(eventBytes),
		Timestamp: event.Timestamp,
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("event_type"),
				Value: []byte(event.Type),
			},
			{
				Key:   []byte("saga_id"),
				Value: []byte(event.SagaID),
			},
		},
	}

	for key, value := range event.Headers {
		message.Headers = append(message.Headers, sarama.RecordHeader{
			Key:   []byte(key),
			Value: []byte(value),
		})
	}

	_, _, err = p.producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (p *Producer) PublishEventWithRetry(topic string, event events.Event, maxRetries int) error {
	var lastErr error
	backoff := time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff *= 2
			if backoff > 30*time.Second {
				backoff = 30 * time.Second
			}
		}

		if err := p.PublishEvent(topic, event); err != nil {
			lastErr = err
			continue
		}

		return nil
	}

	return fmt.Errorf("failed after %d attempts: %w", maxRetries+1, lastErr)
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
