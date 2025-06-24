package user

import (
	"context"
	"userservice/internal/entity"
)

func (repository *Repository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	const op = "userRepository.CreateUser"

	query := `
		INSERT INTO users(email, username)
		VALUES ($1, $2)
		RETURNING id, email, username, created_at, updated_at
	`

	row := repository.db.QueryRowContext(ctx, query, user.Email, user.Username)

	var createdUser entity.User

	err := row.Scan(
		&createdUser.ID,
		&createdUser.Email,
		&createdUser.Username,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)

	if err != nil {
		repository.Log.Error(op, err.Error())
		return nil, err
	}

	return &createdUser, nil
}
