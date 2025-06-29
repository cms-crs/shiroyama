package team

import (
	"context"
	"fmt"
	"taskservice/internal/dto"
)

func (r *Repository) RemoveUserFromTeam(ctx context.Context, req *dto.RemoveUserFromTeamRequest) error {
	const op = "TeamRepository.RemoveUserFromTeam"

	query := "DELETE FROM team_members WHERE team_id = $1 AND user_id = $2"

	result, err := r.db.ExecContext(ctx, query, req.TeamID, req.UserID)
	if err != nil {
		r.log.Error(op, err)
		return fmt.Errorf("failted to remove user from team")
	}

	if count, _ := result.RowsAffected(); count == 0 {
		return fmt.Errorf("failted to remove user from team")
	}

	return nil
}
