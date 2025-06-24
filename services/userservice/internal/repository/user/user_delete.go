package user

import (
	"context"
	"fmt"
)

func (repository *Repository) DeleteUser(ctx context.Context, ID string) error {
	const op = "userRepository.DeleteUser"

	query := `DELETE FROM users WHERE id = $1`

	result, err := repository.db.ExecContext(ctx, query, ID)
	if err != nil {
		repository.Log.Error(op, err.Error())
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.Log.Error(op, err.Error())
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
