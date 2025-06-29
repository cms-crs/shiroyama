package kafka

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

type AuthRepository interface {
	BeginTx(ctx context.Context) (*gorm.DB, error)
	SoftDeleteUserTx(ctx context.Context, tx *gorm.DB, userID string) error
	RestoreUserTx(ctx context.Context, tx *gorm.DB, userID string) error
}

type AuthConsumer struct {
	userRepo      AuthRepository
	kafkaProducer *Producer
	logger        *slog.Logger
}

func NewAuthConsumer(repo AuthRepository, producer *Producer, logger *slog.Logger) *AuthConsumer {
	return &AuthConsumer{
		userRepo:      repo,
		kafkaProducer: producer,
		logger:        logger,
	}
}

func (uc *AuthConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	uc.logger.Info("Consumer group session setup")
	return nil
}

func (uc *AuthConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	uc.logger.Info("Consumer group session cleanup")
	return nil
}

func (uc *AuthConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		uc.logger.Info("Consumed message", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset)

		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			uc.logger.Error("Failed to unmarshal event", "error", err)
			// можно skip или retry
			session.MarkMessage(msg, "")
			continue
		}

		err := uc.HandleEvent(session.Context(), event)
		if err != nil {
			uc.logger.Error("Failed to handle event", "error", err)
			// Тут можно не отмечать сообщение как обработанное, чтобы попытаться заново
			// или сделать dead-letter, зависит от логики
		} else {
			session.MarkMessage(msg, "")
		}
	}

	return nil
}

func (uc *AuthConsumer) HandleEvent(ctx context.Context, event Event) error {
	switch event.Type {
	case UserDeletionRequested:
		return uc.handleUserDeletion(ctx, event)
	case UserDeletionRollback:
		return uc.handleUserDeletionRollback(ctx, event)
	default:
		uc.logger.Warn("Unhandled event type", "type", event.Type)
		return nil
	}
}

func (uc *AuthConsumer) handleUserDeletion(ctx context.Context, event Event) error {
	tx, err := uc.userRepo.BeginTx(ctx)
	if err != nil {
		return err
	}

	if err := uc.userRepo.SoftDeleteUserTx(ctx, tx, event.UserID); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			uc.logger.Error("Failed to rollback transaction", "error", rbErr)
		}

		rollbackEvent := Event{
			ID:        generateUUID(),
			Type:      UserDeletionRollback,
			UserID:    event.UserID,
			Timestamp: time.Now(),
			SagaID:    event.SagaID,
			Data: map[string]interface{}{
				"error": err.Error(),
			},
		}
		_ = uc.kafkaProducer.PublishEvent("user-deletion-saga", rollbackEvent)

		return err
	}

	if err := tx.Commit(); err != nil {
		rollbackEvent := Event{
			ID:        generateUUID(),
			Type:      UserDeletionRollback,
			UserID:    event.UserID,
			Timestamp: time.Now(),
			SagaID:    event.SagaID,
			Data: map[string]interface{}{
				"error": err,
			},
		}
		_ = uc.kafkaProducer.PublishEvent("user-deletion-saga", rollbackEvent)

		return err.Error
	}

	completedEvent := Event{
		ID:        generateUUID(),
		Type:      UserDeletionCompleted,
		UserID:    event.UserID,
		Timestamp: time.Now(),
		SagaID:    event.SagaID,
	}

	if err := uc.kafkaProducer.PublishEvent("user-deletion-saga", completedEvent); err != nil {
		return err
	}

	uc.logger.Info("User deleted successfully", "user_id", event.UserID, "saga_id", event.SagaID)
	return nil
}

func (uc *AuthConsumer) handleUserDeletionRollback(ctx context.Context, event Event) error {
	tx, err := uc.userRepo.BeginTx(ctx)
	if err != nil {
		uc.logger.Error("Failed to begin transaction for rollback", "user_id", event.UserID, "error", err)
		return err
	}

	if err := uc.userRepo.RestoreUserTx(ctx, tx, event.UserID); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			uc.logger.Error("Failed to rollback transaction during rollback handling", "user_id", event.UserID, "error", rbErr)
		}
		uc.logger.Error("Failed to restore user during rollback", "user_id", event.UserID, "error", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		uc.logger.Error("Failed to commit transaction during rollback handling", "user_id", event.UserID, "error", err)
		return err.Error
	}

	uc.logger.Info("User restoration (rollback) completed successfully", "user_id", event.UserID, "saga_id", event.SagaID)
	return nil
}

func generateUUID() string {
	return uuid.New().String()
}
