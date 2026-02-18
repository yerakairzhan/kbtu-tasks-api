package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"tasks_assignment/internal/models"
	"tasks_assignment/internal/usecase"
)

type mockTaskUsecase struct {
	getTasksFn       func(ctx context.Context, done *bool) ([]models.Task, error)
	getTaskByIDFn    func(ctx context.Context, id int) (*models.Task, error)
	createTaskFn     func(ctx context.Context, title string) (*models.Task, error)
	updateTaskDoneFn func(ctx context.Context, id int, done bool) error
	deleteTaskFn     func(ctx context.Context, id int) error
}

func (m *mockTaskUsecase) GetTasks(ctx context.Context, done *bool) ([]models.Task, error) {
	return m.getTasksFn(ctx, done)
}
func (m *mockTaskUsecase) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	return m.getTaskByIDFn(ctx, id)
}
func (m *mockTaskUsecase) CreateTask(ctx context.Context, title string) (*models.Task, error) {
	return m.createTaskFn(ctx, title)
}
func (m *mockTaskUsecase) UpdateTaskDone(ctx context.Context, id int, done bool) error {
	return m.updateTaskDoneFn(ctx, id, done)
}
func (m *mockTaskUsecase) DeleteTask(ctx context.Context, id int) error {
	return m.deleteTaskFn(ctx, id)
}

func TestCreateTask_Success(t *testing.T) {
	h := NewTaskHandler(&mockTaskUsecase{
		createTaskFn: func(_ context.Context, title string) (*models.Task, error) {
			if title != "write tests" {
				t.Fatalf("unexpected title: %s", title)
			}
			return &models.Task{ID: 1, Title: title, Done: false}, nil
		},
	})

	body := bytes.NewBufferString(`{"title":"write tests"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/tasks", body)
	w := httptest.NewRecorder()

	h.CreateTask(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}

	var got models.Task
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.ID != 1 || got.Title != "write tests" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestGetTasks_InvalidID(t *testing.T) {
	h := NewTaskHandler(&mockTaskUsecase{
		getTaskByIDFn: func(_ context.Context, _ int) (*models.Task, error) {
			return nil, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/tasks?id=abc", nil)
	w := httptest.NewRecorder()

	h.GetTasks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	h := NewTaskHandler(&mockTaskUsecase{
		deleteTaskFn: func(_ context.Context, _ int) error {
			return usecase.ErrTaskNotFound
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/v1/tasks?id=7", nil)
	w := httptest.NewRecorder()

	h.DeleteTask(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}

	var body map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "task not found" {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestCreateTask_InvalidTitle(t *testing.T) {
	h := NewTaskHandler(&mockTaskUsecase{
		createTaskFn: func(_ context.Context, _ string) (*models.Task, error) {
			return nil, usecase.ErrInvalidTitle
		},
	})

	body := bytes.NewBufferString(`{"title":""}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/tasks", body)
	w := httptest.NewRecorder()

	h.CreateTask(w, req)
	if !errors.Is(usecase.ErrInvalidTitle, usecase.ErrInvalidTitle) {
		t.Fatalf("sanity check failed")
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
