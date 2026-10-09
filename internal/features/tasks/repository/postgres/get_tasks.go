package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/nullphosphor/golang-todoapp/internal/core/domain"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	userID, limit, offset *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		SELECT id, version, title, description, completed, created_at, completed_at, author_user_id 
		FROM todoapp.tasks
		%s
		ORDER BY id ASC
		LIMIT  $1
		OFFSET $2;
	`

	args := []any{limit, offset}

	if userID != nil {
		sqlQuery = fmt.Sprintf(sqlQuery, "WHERE author_user_id=$3")
		args = append(args, *userID)
	} else {
		sqlQuery = fmt.Sprintf(sqlQuery, "")
	}

	rows, err := r.pool.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel

	for rows.Next() {
		var taskModel TaskModel
		err := rows.Scan(
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
			return nil, fmt.Errorf("scan tasks: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("")
	}

	taskDomains := taskDomainsFromModels(taskModels)

	return taskDomains, nil
}
