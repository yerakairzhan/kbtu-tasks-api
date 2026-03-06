package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"tasks_assignment/internal/models"
	usersrepo "tasks_assignment/internal/repository/_postgres/users"
)

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidUserInput = errors.New("invalid user input")

type UserRepository interface {
	List(ctx context.Context, page, pageSize int) ([]models.User, int, error)
	GetByID(ctx context.Context, id models.UUID) (*models.User, error)
	Create(ctx context.Context, name, email, gender string, birthDate time.Time) (*models.User, error)
}

type UserUsecase struct {
	repo UserRepository
}

func NewUserUsecase(repo UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) GetUsers(ctx context.Context, page, pageSize int) (models.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	users, total, err := u.repo.List(ctx, page, pageSize)
	if err != nil {
		return models.PaginatedResponse{}, err
	}
	return models.NewPaginatedResponse(users, total, page, pageSize), nil
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrInvalidUserInput
	}

	user, err := u.repo.GetByID(ctx, models.UUID(id))
	if err != nil {
		if errors.Is(err, usersrepo.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.TrimSpace(req.Email)
	gender := strings.TrimSpace(req.Gender)

	if name == "" || email == "" || gender == "" || req.BirthDate.IsZero() {
		return nil, ErrInvalidUserInput
	}

	return u.repo.Create(ctx, name, email, gender, req.BirthDate)
}
