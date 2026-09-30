package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/nullphosphor/golang-todoapp/internal/core/domain"
	core_errors "github.com/nullphosphor/golang-todoapp/internal/core/errors"
)

func (r *UsersRepository) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
	SELECT * FROM todoapp.users
	WHERE id=$1		
	`
	row := r.pool.QueryRow(ctx, sqlQuery, id)

	var userModel UserModel
	if err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("not found user with id='%d': %w", id, core_errors.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("scan row: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)

	return userDomain, nil
}
