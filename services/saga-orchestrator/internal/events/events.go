package events

type EventType string

const (
	UserDeletionRequested EventType = "UserDeletionRequested"
	UserDeletionCompleted EventType = "UserDeletionCompleted"
	UserDeletionRollback  EventType = "UserDeletionRollback"

	AuthUserDeleteRequested EventType = "AuthUserDeleteRequested"
	AuthUserDeleted         EventType = "AuthUserDeleted"
	AuthUserDeleteFailed    EventType = "AuthUserDeleteFailed"
	AuthUserDeleteRollback  EventType = "AuthUserDeleteRollback"

	TeamUserDeleteRequested EventType = "TeamUserDeleteRequested"
	TeamUserDeleted         EventType = "TeamUserDeleted"
	TeamUserDeleteFailed    EventType = "TeamUserDeleteFailed"
	TeamUserDeleteRollback  EventType = "TeamUserDeleteRollback"

	BoardUserDeleteRequested EventType = "BoardUserDeleteRequested"
	BoardUserDeleted         EventType = "BoardUserDeleted"
	BoardUserDeleteFailed    EventType = "BoardUserDeleteFailed"
	BoardUserDeleteRollback  EventType = "BoardUserDeleteRollback"

	TaskUserDeleteRequested EventType = "TaskUserDeleteRequested"
	TaskUserDeleted         EventType = "TaskUserDeleted"
	TaskUserDeleteFailed    EventType = "TaskUserDeleteFailed"
	TaskUserDeleteRollback  EventType = "TaskUserDeleteRollback"
)

type SagaStatus string

const (
	SagaStatusPending     SagaStatus = "pending"
	SagaStatusInProgress  SagaStatus = "in_progress"
	SagaStatusCompleted   SagaStatus = "completed"
	SagaStatusRollingBack SagaStatus = "rolling_back"
	SagaStatusRolledBack  SagaStatus = "rolled_back"
)
