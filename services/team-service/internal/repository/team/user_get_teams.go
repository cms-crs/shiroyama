package team

import (
	"context"
	"database/sql"
	"taskservice/internal/entity"
)

func (r *Repository) GetUserTeams(ctx context.Context, UserID string) ([]*entity.Team, error) {
	const op = "TeamRepository.GetUserTeams"

	query := `
		SELECT t.id, t.name, t.description, t.created_at, t.updated_at
		FROM team_members as tm
		JOIN teams as t ON tm.team_id = t.id
		WHERE tm.user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, UserID)
	if err != nil {
		r.log.Error(op, err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.log.Error(op, err)
		}
	}(rows)

	var teams []*entity.Team
	for rows.Next() {
		var team entity.Team

		err := rows.Scan(
			&team.ID,
			&team.Name,
			&team.Description,
			&team.CreatedAt,
			&team.UpdatedAt,
		)

		if err != nil {
			r.log.Error(op, err)
			return nil, err
		}

		teams = append(teams, &team)
	}

	if rows.Err() != nil {
		r.log.Error(op, rows.Err())
		return nil, err
	}

	return teams, nil
}
