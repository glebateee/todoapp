package tasks_service

import (
	"context"

	"github.com/glebateee/todoapp/internal/core/domain"
)

type TaskService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)
}

func NewTaskService(
	tasksRepository TasksRepository,
) *TaskService {
	return &TaskService{
		tasksRepository: tasksRepository,
	}

}
