package events

import "time"

type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	UserID    string                 `json:"user_id"`
	SagaID    string                 `json:"saga_id"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	Headers   map[string]string      `json:"headers"`
}
