package dto

import "time"

type EventType string

type DeleteUserEvent struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	UserID    string                 `json:"user_id"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data,omitempty"`
	SagaID    string                 `json:"saga_id"`
}
