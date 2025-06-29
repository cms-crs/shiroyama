package main

import (
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
)

type EventType string

const (
	UserDeletionRequested EventType = "user.deletion.requested"
	UserDeletionCompleted EventType = "user.deletion.completed"
	UserDeletionFailed    EventType = "user.deletion.failed"
	UserDeletionRollback  EventType = "user.deletion.rollback"
	AuthUserDeleted       EventType = "auth.user.deleted"
	TeamUserDeleted       EventType = "team.user.deleted"
	BoardUserDeleted      EventType = "board.user.deleted"
	TaskUserDeleted       EventType = "task.user.deleted"
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
	Status         string    `json:"status"` // pending, completed, failed, rolling_back
	CompletedSteps []string  `json:"completed_steps"`
	FailedStep     string    `json:"failed_step,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Compression = sarama.CompressionSnappy

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{producer: producer}, nil
}

func (kp *KafkaProducer) PublishEvent(topic string, event Event) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(event.UserID),
		Value: sarama.ByteEncoder(eventBytes),
	}

	_, _, err = kp.producer.SendMessage(msg)
	return err
}

//type TeamService struct {
//	producer *KafkaProducer
//}
//
//func NewTeamService(producer *KafkaProducer) *TeamService {
//	return &TeamService{producer: producer}
//}
//
//func (ts *TeamService) HandleUserDeletion(event Event) error {
//	log.Printf("TeamService: Removing user %s from all teams", event.UserID)
//
//	if err := ts.removeUserFromAllTeams(event.UserID); err != nil {
//		failureEvent := Event{
//			ID:        generateEventID(),
//			Type:      "team.user.delete.failed",
//			UserID:    event.UserID,
//			Timestamp: time.Now(),
//			SagaID:    event.SagaID,
//			Data: map[string]interface{}{
//				"error": err.Error(),
//			},
//		}
//		return ts.producer.PublishEvent("user-deletion-saga", failureEvent)
//	}
//
//	successEvent := Event{
//		ID:        generateEventID(),
//		Type:      TeamUserDeleted,
//		UserID:    event.UserID,
//		Timestamp: time.Now(),
//		SagaID:    event.SagaID,
//	}
//
//	return ts.producer.PublishEvent("user-deletion-saga", successEvent)
//}
//
//func (ts *TeamService) HandleRollback(event Event) error {
//	log.Printf("TeamService: Rolling back user removal for %s", event.UserID)
//	return ts.restoreUserInTeams(event.UserID)
//}
//
//func (ts *TeamService) removeUserFromAllTeams(userID string) error {
//	log.Printf("Removing user %s from all teams", userID)
//	return nil
//}
//
//func (ts *TeamService) restoreUserInTeams(userID string) error {
//	log.Printf("Restoring user %s in teams", userID)
//	return nil
//}
//
//// BoardService обработчик для board_service
//type BoardService struct {
//	producer *KafkaProducer
//}
//
//func NewBoardService(producer *KafkaProducer) *BoardService {
//	return &BoardService{producer: producer}
//}
//
//func (bs *BoardService) HandleUserDeletion(event Event) error {
//	log.Printf("BoardService: Processing user %s deletion", event.UserID)
//
//	if err := bs.handleUserBoardDeletion(event.UserID); err != nil {
//		failureEvent := Event{
//			ID:        generateEventID(),
//			Type:      "board.user.delete.failed",
//			UserID:    event.UserID,
//			Timestamp: time.Now(),
//			SagaID:    event.SagaID,
//			Data: map[string]interface{}{
//				"error": err.Error(),
//			},
//		}
//		return bs.producer.PublishEvent("user-deletion-saga", failureEvent)
//	}
//
//	successEvent := Event{
//		ID:        generateEventID(),
//		Type:      BoardUserDeleted,
//		UserID:    event.UserID,
//		Timestamp: time.Now(),
//		SagaID:    event.SagaID,
//	}
//
//	return bs.producer.PublishEvent("user-deletion-saga", successEvent)
//}
//
//func (bs *BoardService) HandleRoll
