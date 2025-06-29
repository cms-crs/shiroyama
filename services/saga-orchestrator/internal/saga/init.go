package saga

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"saga-orchestrator/internal/config"
	"saga-orchestrator/internal/events"
	"saga-orchestrator/internal/kafka"
)

type Storage interface {
	SaveSagaState(ctx context.Context, state *events.SagaState) error
	GetSagaState(ctx context.Context, sagaID string) (*events.SagaState, error)
	DeleteSagaState(ctx context.Context, sagaID string) error
	GetExpiredSagas(ctx context.Context) ([]*events.SagaState, error)
	GetSagasByStatus(ctx context.Context, status events.SagaStatus) ([]*events.SagaState, error)
	UpdateSagaStep(ctx context.Context, sagaID string, step string, status events.SagaStatus) error
	IncrementRetryCount(ctx context.Context, sagaID string) error
	Close() error
}

type Orchestrator struct {
	producer *kafka.Producer
	storage  Storage
	logger   *slog.Logger
	config   config.SagaConfig
	steps    map[string][]Step
}

type Step struct {
	Name           string
	Topic          string
	EventType      events.EventType
	CompensateType events.EventType
	Timeout        time.Duration
}

func NewOrchestrator(producer *kafka.Producer, storage Storage, log *slog.Logger, cfg config.SagaConfig) *Orchestrator {
	orchestrator := &Orchestrator{
		producer: producer,
		storage:  storage,
		logger:   log,
		config:   cfg,
		steps:    make(map[string][]Step),
	}

	orchestrator.initializeSagaSteps()
	return orchestrator
}

func (o *Orchestrator) initializeSagaSteps() {
	o.steps["user_deletion"] = []Step{
		{
			Name:           "delete_auth_user",
			Topic:          "auth-service-commands",
			EventType:      events.AuthUserDeleteRequested,
			CompensateType: events.AuthUserDeleteRollback,
			Timeout:        30 * time.Second,
		},
		{
			Name:           "delete_team_user",
			Topic:          "team-service-commands",
			EventType:      events.TeamUserDeleteRequested,
			CompensateType: events.TeamUserDeleteRollback,
			Timeout:        30 * time.Second,
		},
		{
			Name:           "delete_board_user",
			Topic:          "board-service-commands",
			EventType:      events.BoardUserDeleteRequested,
			CompensateType: events.BoardUserDeleteRollback,
			Timeout:        30 * time.Second,
		},
		{
			Name:           "delete_task_user",
			Topic:          "task-service-commands",
			EventType:      events.TaskUserDeleteRequested,
			CompensateType: events.TaskUserDeleteRollback,
			Timeout:        30 * time.Second,
		},
	}
}

func (o *Orchestrator) HandleEvent(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.UserDeletionRequested:
		return o.handleUserDeletionRequested(ctx, event)
	case events.AuthUserDeleted:
		return o.handleStepCompleted(ctx, event, "delete_auth_user")
	case events.TeamUserDeleted:
		return o.handleStepCompleted(ctx, event, "delete_team_user")
	case events.BoardUserDeleted:
		return o.handleStepCompleted(ctx, event, "delete_board_user")
	case events.TaskUserDeleted:
		return o.handleStepCompleted(ctx, event, "delete_task_user")
	case events.AuthUserDeleteFailed:
		return o.handleStepFailed(ctx, event, "delete_auth_user")
	case events.TeamUserDeleteFailed:
		return o.handleStepFailed(ctx, event, "delete_team_user")
	case events.BoardUserDeleteFailed:
		return o.handleStepFailed(ctx, event, "delete_board_user")
	case events.TaskUserDeleteFailed:
		return o.handleStepFailed(ctx, event, "delete_task_user")
	default:
		o.logger.Warn("Unknown event type: %s", event.Type)
		return nil
	}
}

func (o *Orchestrator) handleUserDeletionRequested(ctx context.Context, event events.Event) error {
	sagaID := uuid.New().String()

	sagaState := &events.SagaState{
		ID:             sagaID,
		UserID:         event.UserID,
		Status:         events.SagaStatusPending,
		CurrentStep:    "",
		CompletedSteps: []string{},
		RetryCount:     0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(o.config.Timeout),
		Metadata:       make(map[string]string),
	}

	if err := o.storage.SaveSagaState(ctx, sagaState); err != nil {
		return fmt.Errorf("failed to save saga state: %w", err)
	}

	//metrics.SagasStarted.WithLabelValues("user_deletion").Inc()
	o.logger.Info("Started user deletion saga: %s for user: %s", sagaID, event.UserID)

	return o.executeNextStep(ctx, sagaState)
}

func (o *Orchestrator) handleStepCompleted(ctx context.Context, event events.Event, stepName string) error {
	sagaState, err := o.storage.GetSagaState(ctx, event.SagaID)
	if err != nil {
		return fmt.Errorf("failed to get saga state: %w", err)
	}

	if sagaState.Status != events.SagaStatusInProgress {
		o.logger.Warn("Received step completion for saga in wrong state: %s, status: %s",
			event.SagaID, sagaState.Status)
		return nil
	}

	sagaState.CompletedSteps = append(sagaState.CompletedSteps, stepName)
	sagaState.UpdatedAt = time.Now()

	if err := o.storage.SaveSagaState(ctx, sagaState); err != nil {
		return fmt.Errorf("failed to update saga state: %w", err)
	}

	o.logger.Info("Step %s completed for saga: %s", stepName, event.SagaID)

	return o.executeNextStep(ctx, sagaState)
}

func (o *Orchestrator) handleStepFailed(ctx context.Context, event events.Event, stepName string) error {
	sagaState, err := o.storage.GetSagaState(ctx, event.SagaID)
	if err != nil {
		return fmt.Errorf("failed to get saga state: %w", err)
	}

	sagaState.FailedStep = stepName
	sagaState.UpdatedAt = time.Now()

	if sagaState.RetryCount < o.config.MaxRetries {
		return o.retrySagaStep(ctx, sagaState, stepName)
	}

	o.logger.Error("Step %s failed for saga: %s, starting compensation", stepName, event.SagaID)
	//metrics.SagasFailed.WithLabelValues("user_deletion", stepName).Inc()

	return o.startCompensation(ctx, sagaState)
}

func (o *Orchestrator) executeNextStep(ctx context.Context, sagaState *events.SagaState) error {
	steps := o.steps["user_deletion"]

	// Находим следующий шаг для выполнения
	for _, step := range steps {
		if !o.isStepCompleted(sagaState, step.Name) {
			return o.executeStep(ctx, sagaState, step)
		}
	}

	// Все шаги выполнены успешно
	return o.completeSaga(ctx, sagaState)
}

func (o *Orchestrator) executeStep(ctx context.Context, sagaState *events.SagaState, step Step) error {
	sagaState.CurrentStep = step.Name
	sagaState.Status = events.SagaStatusInProgress
	sagaState.UpdatedAt = time.Now()

	if err := o.storage.SaveSagaState(ctx, sagaState); err != nil {
		return fmt.Errorf("failed to update saga state: %w", err)
	}

	event := events.Event{
		ID:        uuid.New().String(),
		Type:      step.EventType,
		UserID:    sagaState.UserID,
		Timestamp: time.Now(),
		SagaID:    sagaState.ID,
		Data:      make(map[string]interface{}),
	}

	if err := o.producer.PublishEvent(step.Topic, event); err != nil {
		return fmt.Errorf("failed to publish step event: %w", err)
	}

	o.logger.Info("Executed step %s for saga: %s", step.Name, sagaState.ID)
	return nil
}

func (o *Orchestrator) retrySagaStep(ctx context.Context, sagaState *events.SagaState, stepName string) error {
	if err := o.storage.IncrementRetryCount(ctx, sagaState.ID); err != nil {
		return fmt.Errorf("failed to increment retry count: %w", err)
	}

	o.logger.Info("Retrying step %s for saga: %s (attempt %d)",
		stepName, sagaState.ID, sagaState.RetryCount+1)

	// Находим шаг для повтора
	steps := o.steps["user_deletion"]
	for _, step := range steps {
		if step.Name == stepName {
			time.Sleep(o.config.RetryInterval)
			return o.executeStep(ctx, sagaState, step)
		}
	}

	return fmt.Errorf("step not found: %s", stepName)
}

func (o *Orchestrator) startCompensation(ctx context.Context, sagaState *events.SagaState) error {
	sagaState.Status = events.SagaStatusRollingBack
	sagaState.UpdatedAt = time.Now()

	if err := o.storage.SaveSagaState(ctx, sagaState); err != nil {
		return fmt.Errorf("failed to update saga state: %w", err)
	}

	steps := o.steps["user_deletion"]

	// Выполняем компенсацию в обратном порядке
	for i := len(steps) - 1; i >= 0; i-- {
		step := steps[i]
		if o.isStepCompleted(sagaState, step.Name) {
			if err := o.compensateStep(ctx, sagaState, step); err != nil {
				o.logger.Error("Failed to compensate step %s: %v", step.Name, err)
				// Продолжаем компенсацию даже если один шаг не удался
			}
		}
	}

	return o.finalizeSagaRollback(ctx, sagaState)
}

func (o *Orchestrator) compensateStep(ctx context.Context, sagaState *events.SagaState, step Step) error {
	event := events.Event{
		ID:        uuid.New().String(),
		Type:      step.CompensateType,
		UserID:    sagaState.UserID,
		Timestamp: time.Now(),
		SagaID:    sagaState.ID,
		Data:      make(map[string]interface{}),
	}

	if err := o.producer.PublishEvent(step.Topic, event); err != nil {
		return fmt.Errorf("failed to publish compensation event: %w", err)
	}

	o.logger.Info("Compensated step %s for saga: %s", step.Name, sagaState.ID)
	return nil
}

func (o *Orchestrator) completeSaga(ctx context.Context, sagaState *events.SagaState) error {
	sagaState.Status = events.SagaStatusCompleted
	sagaState.UpdatedAt = time.Now()

	if err := o.storage.SaveSagaState(ctx, sagaState); err != nil {
		return fmt.Errorf("failed to complete saga: %w", err)
	}

	// Публикуем событие о завершении саги
	event := events.Event{
		ID:        uuid.New().String(),
		Type:      events.UserDeletionCompleted,
		UserID:    sagaState.UserID,
		Timestamp: time.Now(),
		SagaID:    sagaState.ID,
	}

	if err := o.producer.PublishEvent("user-deletion-saga", event); err != nil {
		return fmt.Errorf("failed to publish completion event: %w", err)
	}

	duration := time.Since(sagaState.CreatedAt)
	//metrics.SagasCompleted.WithLabelValues("user_deletion").Inc()
	//metrics.SagaDuration.WithLabelValues("user_deletion", "completed").Observe(duration.Seconds())

	o.logger.Info("Saga completed successfully: %s in %v", sagaState.ID, duration)
	return nil
}

func (o *Orchestrator) finalizeSagaRollback(ctx context.Context, sagaState *events.SagaState) error {
	sagaState.Status = events.SagaStatusRolledBack
	sagaState.UpdatedAt = time.Now()

	if err := o.storage.SaveSagaState(ctx, sagaState); err != nil {
		return fmt.Errorf("failed to finalize saga rollback: %w", err)
	}

	// Публикуем событие об откате саги
	event := events.Event{
		ID:        uuid.New().String(),
		Type:      events.UserDeletionRollback,
		UserID:    sagaState.UserID,
		Timestamp: time.Now(),
		SagaID:    sagaState.ID,
	}

	if err := o.producer.PublishEvent("user-deletion-saga", event); err != nil {
		return fmt.Errorf("failed to publish rollback event: %w", err)
	}

	duration := time.Since(sagaState.CreatedAt)
	//metrics.SagaDuration.WithLabelValues("user_deletion", "rolled_back").Observe(duration.Seconds())

	o.logger.Info("Saga rolled back: %s in %v", sagaState.ID, duration)
	return nil
}

func (o *Orchestrator) isStepCompleted(sagaState *events.SagaState, stepName string) bool {
	for _, completedStep := range sagaState.CompletedSteps {
		if completedStep == stepName {
			return true
		}
	}
	return false
}

func (o *Orchestrator) StartTimeoutMonitor(ctx context.Context) {
	ticker := time.NewTicker(o.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			o.logger.Info("Timeout monitor stopped")
			return
		case <-ticker.C:
			if err := o.handleTimeouts(ctx); err != nil {
				o.logger.Error("Error handling timeouts: %v", err)
			}
		}
	}
}

func (o *Orchestrator) handleTimeouts(ctx context.Context) error {
	expiredSagas, err := o.storage.GetExpiredSagas(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expired sagas: %w", err)
	}

	for _, saga := range expiredSagas {
		o.logger.Warn("Saga timed out: %s (UserID: %s)", saga.ID, saga.UserID)

		if saga.Status == events.SagaStatusCompleted || saga.Status == events.SagaStatusRolledBack {
			continue
		}

		saga.Status = events.SagaStatusRollingBack
		saga.UpdatedAt = time.Now()

		if err := o.storage.SaveSagaState(ctx, saga); err != nil {
			o.logger.Error("Failed to mark saga %s as rolling back: %v", saga.ID, err)
			continue
		}

		if err := o.startCompensation(ctx, saga); err != nil {
			o.logger.Error("Failed to compensate expired saga %s: %v", saga.ID, err)
			continue
		}

		o.logger.Info("Compensated expired saga: %s", saga.ID)
	}

	return nil
}
