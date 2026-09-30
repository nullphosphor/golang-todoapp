package users_service

import (
	"context"
	"fmt"
)

func (s *UsersService) DeleteUser(
	ctx context.Context,
	id int,
) error {
	if err := s.usersRepository.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user with id=%d: %w", id, err)
	}
	return nil
}
