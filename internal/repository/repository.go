package repository

import (
	"context"
	"time"

	"tasks_assignment/internal/models"
	"tasks_assignment/internal/repository/postgres"
	tasksrepo "tasks_assignment/internal/repository/postgres/tasks"
	usersrepo "tasks_assignment/internal/repository/postgres/users"
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
	Users UserRepository
}

type UserRepository interface {
	List(ctx context.Context, input models.ListUsersInput) ([]models.User, int, error)
	GetByID(ctx context.Context, id models.UUID) (*models.User, error)
	Create(ctx context.Context, name, email, gender string, birthDate time.Time) (*models.User, error)
	GetCommonFriends(ctx context.Context, user1, user2 string) ([]models.User, error)
}

func NewRepositories(pg *postgres.Dialect) *Repositories {
	return &Repositories{
		Tasks: tasksrepo.New(pg.DB),
		Users: usersrepo.New(pg.DB),
	}
}
