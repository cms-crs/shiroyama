package list

import (
	"context"
	"fmt"
)

func (r *Repository) DeleteList(ctx context.Context, listID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	//_, err = tx.ExecContext(ctx, "DELETE FROM tasks WHERE list_id = $1", listID)
	//if err != nil {
	//	return fmt.Errorf("failed to delete tasks: %w", err)
	//}

	result, err := tx.ExecContext(ctx, "DELETE FROM lists WHERE id = $1", listID)
	if err != nil {
		return fmt.Errorf("failed to delete list: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("list not found")
	}

	return tx.Commit()
}
