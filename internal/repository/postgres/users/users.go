package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"tasks_assignment/internal/models"
)

var ErrNotFound = errors.New("user not found")

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context, page, pageSize int) ([]models.User, int, error) {
	var total int
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM users`); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	users := make([]models.User, 0, pageSize)
	if err := r.db.SelectContext(
		ctx,
		&users,
		`SELECT id, name, email, gender, birth_date FROM users ORDER BY id LIMIT $1 OFFSET $2`,
		pageSize,
		offset,
	); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *Repository) GetByID(ctx context.Context, id models.UUID) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(ctx, &user, `SELECT id, name, email, gender, birth_date FROM users WHERE id = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) Create(ctx context.Context, name, email, gender string, birthDate time.Time) (*models.User, error) {
	var user models.User
	err := r.db.GetContext(
		ctx,
		&user,
		`INSERT INTO users (name, email, gender, birth_date) VALUES ($1, $2, $3, $4)
		 RETURNING id, name, email, gender, birth_date`,
		name,
		email,
		gender,
		birthDate,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
