package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/nullphosphor/golang-todoapp/internal/core/domain"
	core_errors "github.com/nullphosphor/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/nullphosphor/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		INSERT INTO todoapp.tasks (title, description, completed, created_at, completed_at, author_user_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id;
	`

	row := r.pool.QueryRow(ctx,
		sqlQuery,
		task.Title, task.Description, task.Completed, task.CreatedAt, task.CompletedAt, task.AuthorUserID,
	)

	var taskModel TaskModel
	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Task{}, fmt.Errorf(
				"%v: user with id:'%d': %w",
				err, task.AuthorUserID,
				core_errors.ErrNotFound,
			)
		}
		return domain.Task{}, fmt.Errorf("scan error create_task from repositrory: %w", err)
	}
	taskDomain := domain.NewTask(taskModel.ID, taskModel.Version, taskModel.Title, taskModel.Description, taskModel.Completed, taskModel.CreatedAt, taskModel.CompletedAt, taskModel.AuthorUserID)

	return taskDomain, nil
}
