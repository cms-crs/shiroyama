package models

import "time"

type Board struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TeamID      string    `json:"team_id"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type List struct {
	ID        string    `json:"id"`
	BoardID   string    `json:"board_id"`
	Name      string    `json:"name"`
	Position  int32     `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListWithTasks struct {
	List  List   `json:"list"`
	Tasks []Task `json:"tasks"`
}

type BoardWithLists struct {
	Board Board           `json:"board"`
	Lists []ListWithTasks `json:"lists"`
}

type Label struct {
	ID      string `json:"id"`
	BoardID string `json:"board_id"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

type CreateBoardRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	TeamID      string `json:"team_id" binding:"required"`
	CreatedBy   string `json:"created_by" binding:"required"`
}

type UpdateBoardRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

type CreateListRequest struct {
	BoardID  string `json:"board_id" binding:"required"`
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Position int32  `json:"position" binding:"min=0"`
}

type UpdateListRequest struct {
	Name string `json:"name" binding:"omitempty,min=1,max=100"`
}

type ListPosition struct {
	ListID   string `json:"list_id" binding:"required"`
	Position int32  `json:"position" binding:"min=0"`
}

type ReorderListsRequest struct {
	Positions []ListPosition `json:"positions" binding:"required"`
}

type GetUserBoardsResponse struct {
	Boards []Board `json:"boards"`
}

type GetTeamBoardsResponse struct {
	Boards []Board `json:"boards"`
}

type GetBoardLabelsResponse struct {
	Labels []Label `json:"labels"`
}
