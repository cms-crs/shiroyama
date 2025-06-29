package user

import (
	"context"
	"userservice/internal/entity"
)

func (repository *Repository) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	const op = "userRepository.CreateUser"

	query := `
		UPDATE users
		SET username = ($1)
		WHERE id = $2
		RETURNING id, email, username, created_at, updated_at
	`

	row := repository.db.QueryRowContext(ctx, query, user.Username, user.ID)

	var updatedUser entity.User

	err := row.Scan(
		&updatedUser.ID,
		&updatedUser.Email,
		&updatedUser.Username,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err != nil {
		repository.log.Error(op)
		return nil, err
	}

	return &updatedUser, nil
}
