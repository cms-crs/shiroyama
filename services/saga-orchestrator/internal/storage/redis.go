package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"

	"saga-orchestrator/internal/config"
	"saga-orchestrator/internal/events"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(cfg config.RedisConfig) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStorage{
		client: client,
	}, nil
}

func (rs *RedisStorage) SaveSagaState(ctx context.Context, state *events.SagaState) error {
	key := fmt.Sprintf("saga:%s", state.ID)

	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal saga state: %w", err)
	}

	ttl := time.Until(state.ExpiresAt)
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}

	if err := rs.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save saga state: %w", err)
	}

	return nil
}

func (rs *RedisStorage) GetSagaState(ctx context.Context, sagaID string) (*events.SagaState, error) {
	key := fmt.Sprintf("saga:%s", sagaID)

	data, err := rs.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("saga not found: %s", sagaID)
		}
		return nil, fmt.Errorf("failed to get saga state: %w", err)
	}

	var state events.SagaState
	if err := json.Unmarshal([]byte(data), &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal saga state: %w", err)
	}

	return &state, nil
}

func (rs *RedisStorage) DeleteSagaState(ctx context.Context, sagaID string) error {
	key := fmt.Sprintf("saga:%s", sagaID)

	if err := rs.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete saga state: %w", err)
	}

	return nil
}

func (rs *RedisStorage) GetExpiredSagas(ctx context.Context) ([]*events.SagaState, error) {
	pattern := "saga:*"
	var cursor uint64
	var expiredSagas []*events.SagaState

	for {
		keys, nextCursor, err := rs.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to scan keys: %w", err)
		}

		for _, key := range keys {
			data, err := rs.client.Get(ctx, key).Result()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					continue
				}
				return nil, fmt.Errorf("failed to get saga state: %w", err)
			}

			var state events.SagaState
			if err := json.Unmarshal([]byte(data), &state); err != nil {
				continue
			}

			if time.Now().After(state.ExpiresAt) &&
				(state.Status == events.SagaStatusPending || state.Status == events.SagaStatusInProgress) {
				expiredSagas = append(expiredSagas, &state)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return expiredSagas, nil
}

func (rs *RedisStorage) GetSagasByStatus(ctx context.Context, status events.SagaStatus) ([]*events.SagaState, error) {
	pattern := "saga:*"
	var cursor uint64
	var sagas []*events.SagaState

	for {
		keys, nextCursor, err := rs.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to scan keys: %w", err)
		}

		for _, key := range keys {
			data, err := rs.client.Get(ctx, key).Result()
			if err != nil {
				if err == redis.Nil {
					continue
				}
				return nil, fmt.Errorf("failed to get saga state: %w", err)
			}

			var state events.SagaState
			if err := json.Unmarshal([]byte(data), &state); err != nil {
				continue
			}

			if state.Status == status {
				sagas = append(sagas, &state)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return sagas, nil
}

func (rs *RedisStorage) UpdateSagaStep(ctx context.Context, sagaID string, step string, status events.SagaStatus) error {
	state, err := rs.GetSagaState(ctx, sagaID)
	if err != nil {
		return err
	}

	state.CurrentStep = step
	state.Status = status
	state.UpdatedAt = time.Now()

	if status == events.SagaStatusInProgress {
		found := false
		for _, completedStep := range state.CompletedSteps {
			if completedStep == step {
				found = true
				break
			}
		}
		if !found {
			state.CompletedSteps = append(state.CompletedSteps, step)
		}
	}

	return rs.SaveSagaState(ctx, state)
}

func (rs *RedisStorage) IncrementRetryCount(ctx context.Context, sagaID string) error {
	state, err := rs.GetSagaState(ctx, sagaID)
	if err != nil {
		return err
	}

	state.RetryCount++
	state.UpdatedAt = time.Now()

	return rs.SaveSagaState(ctx, state)
}

func (rs *RedisStorage) Close() error {
	return rs.client.Close()
}
