package models

import "time"

type Comment struct {
	ID           string       `json:"id"`
	TaskID       string       `json:"task_id"`
	UserID       string       `json:"user_id"`
	Content      string       `json:"content"`
	Attachments  []Attachment `json:"attachments"`
	Mentions     []string     `json:"mentions"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	ParentID     string       `json:"parent_id,omitempty"`
	Reactions    []Reaction   `json:"reactions"`
	RepliesCount int32        `json:"replies_count"`
}

type Attachment struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

type Reaction struct {
	Emoji string   `json:"emoji"`
	Users []string `json:"users"`
}

type CreateCommentRequest struct {
	TaskID      string       `json:"task_id" binding:"required"`
	UserID      string       `json:"user_id" binding:"required"`
	Content     string       `json:"content" binding:"required,min=1,max=2000"`
	Attachments []Attachment `json:"attachments"`
	Mentions    []string     `json:"mentions"`
	ParentID    string       `json:"parent_id"`
}

type UpdateCommentRequest struct {
	Content     string       `json:"content" binding:"omitempty,min=1,max=2000"`
	Attachments []Attachment `json:"attachments"`
	Mentions    []string     `json:"mentions"`
}

type GetTaskCommentsRequest struct {
	IncludeReplies bool `json:"include_replies" form:"include_replies"`
	PaginationRequest
}

type GetTaskCommentsResponse struct {
	Comments []Comment `json:"comments"`
	PaginationResponse
}

type GetCommentRepliesResponse struct {
	Replies []Comment `json:"replies"`
	PaginationResponse
}

type AddReactionRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Emoji  string `json:"emoji" binding:"required"`
}

type RemoveReactionRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Emoji  string `json:"emoji" binding:"required"`
}
