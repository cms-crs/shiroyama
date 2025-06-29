package team

import (
	"context"
	"fmt"
	"taskservice/internal/dto"
)

func (r *Repository) UpdateUserRole(ctx context.Context, req *dto.UpdateUserRoleRequest) error {
	const op = "TeamRepository.UpdateUserRole"

	query := `
		UPDATE team_members
		SET role = $1
		WHERE user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, req.Role, req.UserID)
	if err != nil {
		r.log.Error(op, err)
		return fmt.Errorf("failed to remove user from team")
	}

	if count, _ := result.RowsAffected(); count == 0 {
		return fmt.Errorf("failed to remove user from team")
	}

	return nil
}
