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

type dbUserRow struct {
	ID        string       `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	Gender    string       `db:"gender"`
	BirthDate sql.NullTime `db:"birth_date"`
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func toModel(row dbUserRow) models.User {
	birthDate := row.BirthDate.Time
	if !row.BirthDate.Valid {
		birthDate = sql.NullTime{}.Time
	}
	return models.NewUser(models.UUID(row.ID), row.Name, row.Email, row.Gender, birthDate)
}

func (r *Repository) List(ctx context.Context, page, pageSize int) ([]models.User, int, error) {
	var total int
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM users`); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows := make([]dbUserRow, 0, pageSize)
	if err := r.db.SelectContext(
		ctx,
		&rows,
		`SELECT id, name, email, gender, birth_date FROM users ORDER BY name LIMIT $1 OFFSET $2`,
		pageSize,
		offset,
	); err != nil {
		return nil, 0, err
	}

	users := make([]models.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, toModel(row))
	}

	return users, total, nil
}

func (r *Repository) GetByID(ctx context.Context, id models.UUID) (*models.User, error) {
	var row dbUserRow
	err := r.db.GetContext(ctx, &row, `SELECT id, name, email, gender, birth_date FROM users WHERE id = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	user := toModel(row)
	return &user, nil
}

func (r *Repository) Create(ctx context.Context, name, email, gender string, birthDate time.Time) (*models.User, error) {
	var row dbUserRow
	err := r.db.GetContext(
		ctx,
		&row,
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

	user := toModel(row)
	return &user, nil
}
