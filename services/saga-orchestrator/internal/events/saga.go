package events

import "time"

type SagaState struct {
	ID             string
	UserID         string
	Status         SagaStatus
	CurrentStep    string
	FailedStep     string
	CompletedSteps []string
	RetryCount     int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ExpiresAt      time.Time
	Metadata       map[string]string
}
