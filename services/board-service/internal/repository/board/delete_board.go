package board

import (
	"context"
	"fmt"
)

func (r *Repository) DeleteBoard(ctx context.Context, boardID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	//_, err = tx.ExecContext(ctx, `
	//	DELETE FROM tasks
	//	WHERE list_id IN (SELECT id FROM lists WHERE board_id = $1)
	//`, boardID)
	//if err != nil {
	//	return fmt.Errorf("failed to delete tasks: %w", err)
	//}

	_, err = tx.ExecContext(ctx, "DELETE FROM lists WHERE board_id = $1", boardID)
	if err != nil {
		return fmt.Errorf("failed to delete lists: %w", err)
	}

	//_, err = tx.ExecContext(ctx, "DELETE FROM labels WHERE board_id = $1", boardID)
	//if err != nil {
	//	return fmt.Errorf("failed to delete labels: %w", err)
	//}

	result, err := tx.ExecContext(ctx, "DELETE FROM boards WHERE id = $1", boardID)
	if err != nil {
		return fmt.Errorf("failed to delete board: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("board not found")
	}

	return tx.Commit()
}
