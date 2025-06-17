package dto

type UpdateUserRoleRequest struct {
	TeamID string
	UserID string
	Role   string
}
