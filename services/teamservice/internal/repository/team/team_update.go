package team

import (
	"context"
	"fmt"
	"strings"
	"userservice/internal/entity"
)

func (repository *Repository) UpdateTeam(ctx context.Context, team *entity.Team) (*entity.Team, error) {
	const op = "TeamRepository.UpdateTeam"

	var setClauses []string
	var args []interface{}
	argIdx := 1

	if team.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, team.Name)
		argIdx++
	}

	if team.Description != "" {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, team.Description)
		argIdx++
	}

	if len(setClauses) == 0 {
		return team, nil
	}

	query := fmt.Sprintf(`
		UPDATE teams
		SET %s
		WHERE id = $%d
		RETURNING id, name, description, created_at, updated_at
	`, strings.Join(setClauses, ", "), argIdx)

	args = append(args, team.ID)

	row := repository.db.QueryRowContext(ctx, query, args...)

	var updated entity.Team
	err := row.Scan(
		&updated.ID,
		&updated.Name,
		&updated.Description,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		repository.log.Error(op, err)
		return nil, fmt.Errorf("failed to update team")
	}

	return &updated, nil
}
