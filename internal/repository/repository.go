package repository

import (
	"context"

	"tasks_assignment/internal/models"
	"tasks_assignment/internal/repository/_postgres"
	tasksrepo "tasks_assignment/internal/repository/_postgres/tasks"
)

type TaskRepository interface {
	List(ctx context.Context, done *bool) ([]models.Task, error)
	GetByID(ctx context.Context, id int) (*models.Task, error)
	Create(ctx context.Context, title string) (*models.Task, error)
	UpdateDone(ctx context.Context, id int, done bool) error
	Delete(ctx context.Context, id int) error
}

type Repositories struct {
	Tasks TaskRepository
}

func NewRepositories(postgres *_postgres.Dialect) *Repositories {
	return &Repositories{
		Tasks: tasksrepo.New(postgres.DB),
	}
}
