package team

import (
	"context"
	"taskservice/internal/entity"
)

func (r *Repository) GetTeam(ctx context.Context, ID string) (*entity.Team, error) {
	const op = "repository.GetTeam"

	query := `
		SELECT id, name, description, created_at, updated_at 
		FROM teams WHERE id = $1
	`

	var team entity.Team

	row := r.db.QueryRowContext(ctx, query, ID)

	if err := row.Err(); err != nil {
		r.log.Error(op, row.Err().Error())
		return nil, err
	}

	err := row.Scan(
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

	return &team, nil
}
