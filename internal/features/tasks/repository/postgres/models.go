package tasks_postgres_repository

import (
	"time"

	"github.com/nullphosphor/golang-todoapp/internal/core/domain"
)

type TaskModel struct {
	ID, Version int

	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time

	AuthorUserID int
}

func taskDomainsFromModels(taskModels []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(taskModels))

	for i, model := range taskModels {
		domains[i] = domain.NewTask(
			model.ID,
			model.Version,
			model.Title,
			model.Description,
			model.Completed,
			model.CreatedAt,
			model.CompletedAt,
			model.AuthorUserID,
		)
	}

	return domains
}
