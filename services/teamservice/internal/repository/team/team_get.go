package team

import (
	"context"
	"userservice/internal/entity"
)

func (repository *Repository) GetTeam(ctx context.Context, ID string) (*entity.Team, error) {
	const op = "repository.GetTeam"

	query := `
		SELECT id, name, description, created_at, updated_at 
		FROM teams WHERE id = $1
	`

	var team entity.Team

	row := repository.db.QueryRowContext(ctx, query, ID)

	if err := row.Err(); err != nil {
		repository.log.Error(op, row.Err().Error())
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
		repository.log.Error(op, err)
		return nil, err
	}

	return &team, nil
}
