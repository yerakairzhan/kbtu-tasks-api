package tasks

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"

	"tasks_assignment/internal/models"
)

var ErrNotFound = errors.New("task not found")

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context, done *bool) ([]models.Task, error) {
	const baseQuery = `SELECT id, title, done FROM tasks`

	var tasks []models.Task
	if done == nil {
		if err := r.db.SelectContext(ctx, &tasks, baseQuery+` ORDER BY id`); err != nil {
			return nil, err
		}
		return tasks, nil
	}

	if err := r.db.SelectContext(ctx, &tasks, baseQuery+` WHERE done = $1 ORDER BY id`, *done); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (*models.Task, error) {
	var task models.Task
	err := r.db.GetContext(ctx, &task, `SELECT id, title, done FROM tasks WHERE id = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &task, nil
}

func (r *Repository) Create(ctx context.Context, title string) (*models.Task, error) {
	var task models.Task
	err := r.db.GetContext(
		ctx,
		&task,
		`INSERT INTO tasks (title) VALUES ($1) RETURNING id, title, done`,
		title,
	)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) UpdateDone(ctx context.Context, id int, done bool) error {
	res, err := r.db.ExecContext(ctx, `UPDATE tasks SET done = $1 WHERE id = $2`, done, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
