package kafka

import (
	"time"
)

type EventType string

const (
	TeamUserDeleteRequested EventType = "team.user.delete.requested"
	TeamUserDeleted         EventType = "team.user.deleted"
	TeamUserDeleteFailed    EventType = "team.user.delete.failed"
	TeamUserDeleteRollback  EventType = "team.user.delete.rollback"

	UserDeletionRequested EventType = "user.deletion.requested"
	UserDeletionCompleted EventType = "user.deletion.completed"
	UserDeletionRollback  EventType = "user.deletion.rollback"
)

type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	UserID    string                 `json:"user_id"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data,omitempty"`
	SagaID    string                 `json:"saga_id"`
}

type SagaState struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Status         string    `json:"status"`
	CompletedSteps []string  `json:"completed_steps"`
	FailedStep     string    `json:"failed_step,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type TeamDeletionData struct {
	Teams []TeamMembershipData `json:"teams"`
}

type TeamMembershipData struct {
	TeamID   string    `json:"team_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}
