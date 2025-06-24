package team

import (
	"context"
	"userservice/internal/entity"
)

func (repository *Repository) CreateTeam(ctx context.Context, team *entity.Team) (*entity.Team, error) {
	const op = "TeamRepository.CreateTeam"

	query := `
		INSERT INTO teams (name, description) VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`

	row := repository.db.QueryRowContext(ctx, query, team.Name, team.Description)

	if err := row.Err(); err != nil {
		repository.log.Error(op, err)
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

	return team, nil
}
