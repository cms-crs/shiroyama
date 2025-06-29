package team

import (
	"context"
	"fmt"
	"taskservice/internal/entity"
)

func (r *Repository) AddUserToTeam(ctx context.Context, req *entity.TeamMember) error {
	const op = "TeamRepository.AddUserToTeam"

	query := `
		INSERT INTO team_members (team_id, user_id, role) VALUES ($1, $2, $3)
	`

	result, err := r.db.ExecContext(ctx, query, req.TeamID, req.UserID, req.Role)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.Error(op, err)
		return err
	}

	if rowsAffected == 0 {
		r.log.Warn(op, "User already added to team or insertion skipped")
		return fmt.Errorf("user already added to team or insertion skipped")
	}

	return nil
}
