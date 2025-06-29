package user

import (
	"context"
	"database/sql"
)

func (repository *Repository) SoftDeleteUserTx(ctx context.Context, tx *sql.Tx, userID string) error {
	query := `UPDATE users SET is_deleted = TRUE WHERE id = $1`
	result, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		repository.log.Error("Failed to soft delete user", "user_id", userID, "error", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		repository.log.Warn("No user found to soft delete", "user_id", userID)
	} else {
		repository.log.Info("User soft deleted successfully", "user_id", userID)
	}

	return nil
}
