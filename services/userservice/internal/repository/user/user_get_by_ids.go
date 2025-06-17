package user

import (
	"context"
	"github.com/lib/pq"
	"userservice/internal/entity"
)

func (repository *Repository) GetUsersByIds(ctx context.Context, IDs []string) ([]*entity.User, error) {
	const op = "userRepository.GetUsersByIds"

	query := `
		SELECT * FROM users WHERE id = ANY ($1) 
	`

	rows, err := repository.db.QueryContext(ctx, query, pq.Array(IDs))

	if err != nil {
		repository.Log.Error(op, err.Error())
		return nil, err
	}

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		err = rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			break
		}
		users = append(users, &user)
	}

	if closeErr := rows.Close(); closeErr != nil {
		repository.Log.Error(op, closeErr.Error())
		return nil, closeErr
	}

	if err != nil {
		repository.Log.Error(op, err.Error())
		return nil, err
	}

	if err = rows.Err(); err != nil {
		repository.Log.Error(op, err.Error())
		return nil, err
	}

	return users, nil
}
