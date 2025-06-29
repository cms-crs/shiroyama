package team

import (
	"context"
	"fmt"
	"taskservice/internal/entity"
)

func (r *Repository) GetTeamMembers(ctx context.Context, ID string) ([]*entity.TeamMember, error) {
	const op = "TeamRepository.GetTeamMembers"

	query := `SELECT * FROM team_members WHERE team_id=$1`

	rows, err := r.db.QueryContext(ctx, query, ID)

	if err != nil {
		r.log.Error(op, err)
		return nil, fmt.Errorf("failed to get team members")
	}

	var teamMembers []*entity.TeamMember
	if rows.Next() {
		var teamMember entity.TeamMember
		err := rows.Scan(
			&teamMember.UserID,
			&teamMember.TeamID,
			&teamMember.Role,
			&teamMember.CreatedAt,
			&teamMember.UpdatedAt,
		)

		teamMembers = append(teamMembers, &teamMember)

		if err != nil {
			r.log.Error(op, err)
			return nil, fmt.Errorf("failed to get team members")
		}
	}

	fmt.Println(teamMembers)

	return teamMembers, nil
}
