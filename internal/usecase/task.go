package usecase

import (
	"context"
	"errors"
	"strings"

	"tasks_assignment/internal/models"
	tasksrepo "tasks_assignment/internal/repository/_postgres/tasks"
)

var ErrTaskNotFound = errors.New("task not found")
var ErrInvalidTitle = errors.New("invalid title")

type TaskRepository interface {
	List(ctx context.Context, done *bool) ([]models.Task, error)
	GetByID(ctx context.Context, id int) (*models.Task, error)
	Create(ctx context.Context, title string) (*models.Task, error)
	UpdateDone(ctx context.Context, id int, done bool) error
	Delete(ctx context.Context, id int) error
}

type TaskUsecase struct {
	repo TaskRepository
}

func NewTaskUsecase(repo TaskRepository) *TaskUsecase {
	return &TaskUsecase{repo: repo}
}

func (u *TaskUsecase) GetTasks(ctx context.Context, done *bool) ([]models.Task, error) {
	return u.repo.List(ctx, done)
}

func (u *TaskUsecase) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	task, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, tasksrepo.ErrNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return task, nil
}

func (u *TaskUsecase) CreateTask(ctx context.Context, title string) (*models.Task, error) {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" || len(trimmed) > 100 {
		return nil, ErrInvalidTitle
	}
	return u.repo.Create(ctx, trimmed)
}

func (u *TaskUsecase) UpdateTaskDone(ctx context.Context, id int, done bool) error {
	err := u.repo.UpdateDone(ctx, id, done)
	if err != nil && errors.Is(err, tasksrepo.ErrNotFound) {
		return ErrTaskNotFound
	}
	return err
}

func (u *TaskUsecase) DeleteTask(ctx context.Context, id int) error {
	err := u.repo.Delete(ctx, id)
	if err != nil && errors.Is(err, tasksrepo.ErrNotFound) {
		return ErrTaskNotFound
	}
	return err
}
