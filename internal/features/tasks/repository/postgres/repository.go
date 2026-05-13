package tasks_postgres_repository

import core_postgres_pool "github.com/glebateee/todoapp/internal/core/repository/postgres/pool"

type TasksRepository struct {
	core_postgres_pool.Pool
}

func NewTasksrepository(
	pool core_postgres_pool.Pool,
) *TasksRepository {
	return &TasksRepository{
		Pool: pool,
	}
}
