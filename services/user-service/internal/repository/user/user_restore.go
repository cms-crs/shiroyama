package user

import (
	"context"
	"database/sql"
)

func (repository *Repository) RestoreUserTx(ctx context.Context, tx *sql.Tx, userID string) error {
	query := `UPDATE users SET is_deleted = FALSE WHERE id = $1`
	result, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		repository.log.Error("Failed to restore user", "user_id", userID, "error", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		repository.log.Warn("No user found to restore", "user_id", userID)
	} else {
		repository.log.Info("User restored successfully", "user_id", userID)
	}

	return nil
}
