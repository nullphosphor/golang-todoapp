package domain

import "time"

type Statistics struct {
	TaskCreated                int
	TasksCompleted             int
	TaskCompletedRate          *float64
	TasksAverageCompletionTime *time.Duration
}

func NewStatistics(
	taskCreated, tasksCompleted int,
	taskCompletedRate *float64,
	tasksAverageCompletionTime *time.Duration,
) Statistics {
	return Statistics{
		TaskCreated:                taskCreated,
		TasksCompleted:             tasksCompleted,
		TaskCompletedRate:          taskCompletedRate,
		TasksAverageCompletionTime: tasksAverageCompletionTime,
	}
}
