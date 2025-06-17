package team

import (
	"context"
)

func (repository *Repository) DeleteTeam(ctx context.Context, ID string) error {
	const op = "TeamRepository.DeleteTeam"

	query := `DELETE FROM teams WHERE id = $1`

	result, err := repository.db.ExecContext(ctx, query, ID)

	if err != nil {
		repository.log.Error(op, err.Error())
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		repository.log.Error(op, err.Error())
		return err
	}

	if affected == 0 {
		repository.log.Error(op, "Team not found")
		return err
	}

	return nil
}
