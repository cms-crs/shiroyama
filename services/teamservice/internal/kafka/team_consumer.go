package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type TeamService interface {
	DeleteUserFromAllTeams(ctx context.Context, userID string) (*TeamDeletionData, error)

	RestoreUserTeams(ctx context.Context, userID string, data *TeamDeletionData) error

	GetUserTeamMemberships(ctx context.Context, userID string) (*TeamDeletionData, error)
}

type TeamConsumer struct {
	teamService TeamService
	producer    *Producer
	logger      *slog.Logger
}

func NewTeamConsumer(teamService TeamService, producer *Producer, logger *slog.Logger) *TeamConsumer {
	return &TeamConsumer{
		teamService: teamService,
		producer:    producer,
		logger:      logger,
	}
}

func (tc *TeamConsumer) Setup(sarama.ConsumerGroupSession) error {
	tc.logger.Info("Team consumer setup completed")
	return nil
}

func (tc *TeamConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	tc.logger.Info("Team consumer cleanup completed")
	return nil
}

func (tc *TeamConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			if err := tc.handleMessage(session.Context(), message); err != nil {
				tc.logger.Error("Failed to handle message",
					"error", err,
					"topic", message.Topic,
					"partition", message.Partition,
					"offset", message.Offset)
			} else {
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			tc.logger.Info("Consumer session context cancelled")
			return nil
		}
	}
}

func (tc *TeamConsumer) handleMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	tc.logger.Debug("Received message",
		"topic", message.Topic,
		"partition", message.Partition,
		"offset", message.Offset)

	var event Event
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	tc.logger.Info("Processing event",
		"event_id", event.ID,
		"event_type", event.Type,
		"user_id", event.UserID,
		"saga_id", event.SagaID)

	switch event.Type {
	case TeamUserDeleteRequested:
		return tc.handleTeamUserDeleteRequested(ctx, event)
	case TeamUserDeleteRollback:
		return tc.handleTeamUserDeleteRollback(ctx, event)
	default:
		tc.logger.Warn("Unknown event type", "event_type", event.Type)
		return nil
	}
}

func (tc *TeamConsumer) handleTeamUserDeleteRequested(ctx context.Context, event Event) error {
	tc.logger.Info("Handling team user delete requested",
		"user_id", event.UserID,
		"saga_id", event.SagaID)

	deletionData, err := tc.teamService.DeleteUserFromAllTeams(ctx, event.UserID)
	if err != nil {
		tc.logger.Error("Failed to delete user from teams",
			"user_id", event.UserID,
			"saga_id", event.SagaID,
			"error", err)

		return tc.publishFailureEvent(event, err)
	}

	tc.logger.Info("Successfully deleted user from all teams",
		"user_id", event.UserID,
		"saga_id", event.SagaID,
		"teams_count", len(deletionData.Teams))

	return tc.publishSuccessEvent(event, deletionData)
}

func (tc *TeamConsumer) handleTeamUserDeleteRollback(ctx context.Context, event Event) error {
	tc.logger.Info("Handling team user delete rollback",
		"user_id", event.UserID,
		"saga_id", event.SagaID)

	var deletionData *TeamDeletionData
	if rollbackData, exists := event.Data["rollback_data"]; exists {
		if rollbackBytes, err := json.Marshal(rollbackData); err == nil {
			json.Unmarshal(rollbackBytes, &deletionData)
		}
	}

	if deletionData == nil {
		var err error
		deletionData, err = tc.teamService.GetUserTeamMemberships(ctx, event.UserID)
		if err != nil {
			tc.logger.Error("Failed to get user team memberships for rollback",
				"user_id", event.UserID,
				"saga_id", event.SagaID,
				"error", err)
			return err
		}
	}

	if err := tc.teamService.RestoreUserTeams(ctx, event.UserID, deletionData); err != nil {
		tc.logger.Error("Failed to restore user teams",
			"user_id", event.UserID,
			"saga_id", event.SagaID,
			"error", err)
		return err
	}

	tc.logger.Info("Successfully restored user teams",
		"user_id", event.UserID,
		"saga_id", event.SagaID,
		"teams_count", len(deletionData.Teams))

	return nil
}

func (tc *TeamConsumer) publishSuccessEvent(originalEvent Event, deletionData *TeamDeletionData) error {
	successEvent := Event{
		ID:        uuid.New().String(),
		Type:      TeamUserDeleted,
		UserID:    originalEvent.UserID,
		Timestamp: time.Now(),
		SagaID:    originalEvent.SagaID,
		Data: map[string]interface{}{
			"teams_deleted_count": len(deletionData.Teams),
			"rollback_data":       deletionData,
		},
	}

	if err := tc.producer.PublishEvent("user-deletion-saga", successEvent); err != nil {
		tc.logger.Error("Failed to publish success event",
			"user_id", originalEvent.UserID,
			"saga_id", originalEvent.SagaID,
			"error", err)
		return err
	}

	tc.logger.Info("Published team user deleted event",
		"user_id", originalEvent.UserID,
		"saga_id", originalEvent.SagaID)

	return nil
}

func (tc *TeamConsumer) publishFailureEvent(originalEvent Event, originalError error) error {
	failureEvent := Event{
		ID:        uuid.New().String(),
		Type:      TeamUserDeleteFailed,
		UserID:    originalEvent.UserID,
		Timestamp: time.Now(),
		SagaID:    originalEvent.SagaID,
		Data: map[string]interface{}{
			"error":  originalError.Error(),
			"reason": "failed_to_delete_user_from_teams",
		},
	}

	if err := tc.producer.PublishEvent("user-deletion-saga", failureEvent); err != nil {
		tc.logger.Error("Failed to publish failure event",
			"user_id", originalEvent.UserID,
			"saga_id", originalEvent.SagaID,
			"error", err)
		return err
	}

	tc.logger.Info("Published team user delete failed event",
		"user_id", originalEvent.UserID,
		"saga_id", originalEvent.SagaID,
		"original_error", originalError.Error())

	return nil
}
