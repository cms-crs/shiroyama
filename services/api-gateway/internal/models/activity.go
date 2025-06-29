package models

import "time"

type Activity struct {
	ID         string                 `json:"id"`
	EntityType string                 `json:"entity_type"`
	EntityID   string                 `json:"entity_id"`
	Action     string                 `json:"action"`
	UserID     string                 `json:"user_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Details    map[string]interface{} `json:"details"`
	TaskID     string                 `json:"task_id,omitempty"`
	BoardID    string                 `json:"board_id,omitempty"`
	CommentID  string                 `json:"comment_id,omitempty"`
	TeamID     string                 `json:"team_id,omitempty"`
}

type LogActivityRequest struct {
	EntityType string                 `json:"entity_type" binding:"required"`
	EntityID   string                 `json:"entity_id" binding:"required"`
	Action     string                 `json:"action" binding:"required"`
	UserID     string                 `json:"user_id" binding:"required"`
	Details    map[string]interface{} `json:"details"`
	TaskID     string                 `json:"task_id"`
	BoardID    string                 `json:"board_id"`
	CommentID  string                 `json:"comment_id"`
	TeamID     string                 `json:"team_id"`
}

type GetEntityActivityRequest struct {
	EntityType string `json:"entity_type" form:"entity_type" binding:"required"`
	EntityID   string `json:"entity_id" form:"entity_id" binding:"required"`
	PaginationRequest
}

type GetActivityResponse struct {
	Activities []Activity `json:"activities"`
	PaginationResponse
}

type TaskMoveDetails struct {
	FromListID   string `json:"from_list_id"`
	ToListID     string `json:"to_list_id"`
	FromListName string `json:"from_list_name"`
	ToListName   string `json:"to_list_name"`
}

type TaskUpdateDetails struct {
	ChangedFields map[string]string `json:"changed_fields"`
	OldValues     map[string]string `json:"old_values"`
}

type UserAssignmentDetails struct {
	TaskTitle    string   `json:"task_title"`
	AddedUsers   []string `json:"added_users"`
	RemovedUsers []string `json:"removed_users"`
}
