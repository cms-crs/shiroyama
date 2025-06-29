package models

import "time"

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type PaginationRequest struct {
	Page  int32     `json:"page" form:"page" binding:"min=1"`
	Limit int32     `json:"limit" form:"limit" binding:"min=1,max=100"`
	From  time.Time `json:"from" form:"from"`
	To    time.Time `json:"to" form:"to"`
}

type PaginationResponse struct {
	CurrentPage int32 `json:"current_page"`
	PerPage     int32 `json:"per_page"`
	TotalCount  int32 `json:"total_count"`
	TotalPages  int32 `json:"total_pages"`
}

type ListResponse struct {
	Data       interface{}         `json:"data"`
	Pagination *PaginationResponse `json:"pagination"`
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

func (p *PaginationRequest) SetDefaultPagination() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 20
	}
}

func CalculatePaginationResponse(page, limit, totalCount int32) *PaginationResponse {
	totalPages := (totalCount + limit - 1) / limit
	return &PaginationResponse{
		CurrentPage: page,
		PerPage:     limit,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
	}
}
