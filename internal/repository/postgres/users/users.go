package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

func (r *Repository) List(ctx context.Context, input models.ListUsersInput) ([]models.User, int, error) {
	var total int
	whereClauses := make([]string, 0, 5)
	args := make([]interface{}, 0, 7)
	idx := 1

	filters := input.Filters
	if trimmed := strings.TrimSpace(filters.ID); trimmed != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("id = $%d", idx))
		args = append(args, trimmed)
		idx++
	}
	if trimmed := strings.TrimSpace(filters.Name); trimmed != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("name ILIKE $%d", idx))
		args = append(args, "%"+trimmed+"%")
		idx++
	}
	if trimmed := strings.TrimSpace(filters.Email); trimmed != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("email ILIKE $%d", idx))
		args = append(args, "%"+trimmed+"%")
		idx++
	}
	if trimmed := strings.TrimSpace(filters.Gender); trimmed != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("gender = $%d", idx))
		args = append(args, trimmed)
		idx++
	}
	if filters.BirthDate != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("DATE(birth_date) = $%d", idx))
		args = append(args, filters.BirthDate.Format("2006-01-02"))
		idx++
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM users" + whereSQL
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	offset := (input.Page - 1) * input.PageSize
	allowedOrderBy := map[string]string{
		"id":         "id",
		"name":       "name",
		"email":      "email",
		"gender":     "gender",
		"birth_date": "birth_date",
	}

	orderBy := "id"
	if normalized := strings.ToLower(strings.TrimSpace(input.OrderBy)); normalized != "" {
		if val, ok := allowedOrderBy[normalized]; ok {
			orderBy = val
		}
	}

	users := make([]models.User, 0, input.PageSize)
	selectQuery := fmt.Sprintf(
		"SELECT id, name, email, gender, birth_date FROM users%s ORDER BY %s LIMIT $%d OFFSET $%d",
		whereSQL,
		orderBy,
		len(args)+1,
		len(args)+2,
	)

	selectArgs := append(append(make([]interface{}, 0, len(args)+2), args...), input.PageSize, offset)
	if err := r.db.SelectContext(ctx, &users, selectQuery, selectArgs...); err != nil {
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

func (r *Repository) GetCommonFriends(ctx context.Context, user1, user2 string) ([]models.User, error) {
	friends := make([]models.User, 0)
	query := `
SELECT u.id, u.name, u.email, u.gender, u.birth_date
FROM users u
JOIN user_friends f1 ON u.id = f1.friend_id AND f1.user_id = $1
JOIN user_friends f2 ON u.id = f2.friend_id AND f2.user_id = $2
ORDER BY u.id
`
	if err := r.db.SelectContext(ctx, &friends, query, user1, user2); err != nil {
		return nil, err
	}
	return friends, nil
}
