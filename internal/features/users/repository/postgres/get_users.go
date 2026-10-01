package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/nullphosphor/golang-todoapp/internal/core/domain"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
	SELECT id, version, full_name, phone_number
	FROM todoapp.users
	ORDER BY id ASC
	LIMIT $1	
	OFFSET $2;
	`

	rows, err := r.pool.Query(ctx, sqlQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	var usersModels []UserModel
	for rows.Next() {
		var userModel UserModel

		if err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.FullName,
			&userModel.PhoneNumber,
		); err != nil {
			return nil, fmt.Errorf("scan users: %w", err)
		}
		usersModels = append(usersModels, userModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	userDomain := userDomainsFromModels(usersModels)

	return userDomain, nil
}
