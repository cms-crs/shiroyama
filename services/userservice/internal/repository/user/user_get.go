package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"userservice/internal/entity"
)

func (repository *Repository) GetUser(ctx context.Context, id string) (*entity.User, error) {
	const op = "userRepository.GetUser"

	query := `
		SELECT * FROM users WHERE id=$1 
	`

	row := repository.db.QueryRowContext(ctx, query, id)

	var user entity.User

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.UpdatedAt,
		&user.CreatedAt,
		&user.IsDeleted,
	)

	if err != nil {
		repository.log.Error(op, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}
