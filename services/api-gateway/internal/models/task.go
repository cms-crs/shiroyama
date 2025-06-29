package models

import "time"

type Task struct {
	ID            string    `json:"id"`
	ListID        string    `json:"list_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Position      int32     `json:"position"`
	DueDate       time.Time `json:"due_date,omitempty"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	AssignedUsers []string  `json:"assigned_users"`
	Labels        []Label   `json:"labels"`
}

type CreateTaskRequest struct {
	ListID      string    `json:"list_id" binding:"required"`
	Title       string    `json:"title" binding:"required,min=1,max=200"`
	Description string    `json:"description" binding:"omitempty,max=2000"`
	Position    int32     `json:"position" binding:"min=0"`
	DueDate     time.Time `json:"due_date"`
	CreatedBy   string    `json:"created_by" binding:"required"`
}

type UpdateTaskRequest struct {
	Title       string    `json:"title" binding:"omitempty,min=1,max=200"`
	Description string    `json:"description" binding:"omitempty,max=2000"`
	DueDate     time.Time `json:"due_date"`
}

type MoveTaskRequest struct {
	ToListID string `json:"to_list_id" binding:"required"`
	Position int32  `json:"position" binding:"min=0"`
}

type AssignUserRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type CreateLabelRequest struct {
	BoardID string `json:"board_id" binding:"required"`
	Name    string `json:"name" binding:"required,min=1,max=50"`
	Color   string `json:"color" binding:"required,hexcolor"`
}

type UpdateLabelRequest struct {
	Name  string `json:"name" binding:"omitempty,min=1,max=50"`
	Color string `json:"color" binding:"omitempty,hexcolor"`
}

type GetTasksForUserResponse struct {
	Tasks []Task `json:"tasks"`
}

type GetTasksForListsRequest struct {
	ListIDs []string `json:"list_ids" binding:"required"`
}

type GetTasksForListsResponse struct {
	Tasks []Task `json:"tasks"`
}
