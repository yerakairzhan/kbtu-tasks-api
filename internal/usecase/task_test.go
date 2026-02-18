package usecase

import (
	"context"
	"errors"
	"testing"

	"tasks_assignment/internal/models"
	tasksrepo "tasks_assignment/internal/repository/_postgres/tasks"
)

type mockTaskRepo struct {
	listFn       func(ctx context.Context, done *bool) ([]models.Task, error)
	getByIDFn    func(ctx context.Context, id int) (*models.Task, error)
	createFn     func(ctx context.Context, title string) (*models.Task, error)
	updateDoneFn func(ctx context.Context, id int, done bool) error
	deleteFn     func(ctx context.Context, id int) error
}

func (m *mockTaskRepo) List(ctx context.Context, done *bool) ([]models.Task, error) {
	return m.listFn(ctx, done)
}
func (m *mockTaskRepo) GetByID(ctx context.Context, id int) (*models.Task, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockTaskRepo) Create(ctx context.Context, title string) (*models.Task, error) {
	return m.createFn(ctx, title)
}
func (m *mockTaskRepo) UpdateDone(ctx context.Context, id int, done bool) error {
	return m.updateDoneFn(ctx, id, done)
}
func (m *mockTaskRepo) Delete(ctx context.Context, id int) error {
	return m.deleteFn(ctx, id)
}

func TestCreateTask_ValidatesAndTrims(t *testing.T) {
	repo := &mockTaskRepo{
		createFn: func(_ context.Context, title string) (*models.Task, error) {
			if title != "hello" {
				t.Fatalf("expected trimmed title, got %q", title)
			}
			return &models.Task{ID: 1, Title: title, Done: false}, nil
		},
	}
	uc := NewTaskUsecase(repo)

	task, err := uc.CreateTask(context.Background(), "  hello  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.Title != "hello" {
		t.Fatalf("expected hello, got %s", task.Title)
	}

	_, err = uc.CreateTask(context.Background(), "   ")
	if !errors.Is(err, ErrInvalidTitle) {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}
}

func TestGetTaskByID_MapsNotFound(t *testing.T) {
	repo := &mockTaskRepo{
		getByIDFn: func(_ context.Context, _ int) (*models.Task, error) {
			return nil, tasksrepo.ErrNotFound
		},
	}
	uc := NewTaskUsecase(repo)

	_, err := uc.GetTaskByID(context.Background(), 10)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}
