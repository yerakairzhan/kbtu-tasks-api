package usecase

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"tasks_assignment/internal/models"
	usersrepo "tasks_assignment/internal/repository/postgres/users"
)

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidUserInput = errors.New("invalid user input")

const (
	defaultUsersPage     = 1
	defaultUsersPageSize = 10
	maxUsersPageSize     = 100
)

type UserRepository interface {
	List(ctx context.Context, input models.ListUsersInput) ([]models.User, int, error)
	GetByID(ctx context.Context, id models.UUID) (*models.User, error)
	Create(ctx context.Context, name, email, gender string, birthDate time.Time) (*models.User, error)
	GetCommonFriends(ctx context.Context, user1, user2 string) ([]models.User, error)
}

type UserUsecase struct {
	repo UserRepository
}

func NewUserUsecase(repo UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) GetUsers(ctx context.Context, input models.ListUsersInput) (models.PaginatedResponse, error) {
	page, pageSize := normalizePagination(input.Page, input.PageSize)
	input.Page = page
	input.PageSize = pageSize

	users, total, err := u.repo.List(ctx, input)
	if err != nil {
		return models.PaginatedResponse{}, err
	}
	return models.NewPaginatedResponse(users, total, page, pageSize), nil
}

func (u *UserUsecase) GetCommonFriends(ctx context.Context, user1ID, user2ID string) ([]models.User, error) {
	user1 := strings.TrimSpace(user1ID)
	user2 := strings.TrimSpace(user2ID)
	if user1 == "" || user2 == "" {
		return nil, ErrInvalidUserInput
	}
	return u.repo.GetCommonFriends(ctx, user1, user2)
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	trimmedID := strings.TrimSpace(id)
	if trimmedID == "" {
		return nil, ErrInvalidUserInput
	}

	user, err := u.repo.GetByID(ctx, models.UUID(trimmedID))
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
	email := strings.ToLower(strings.TrimSpace(req.Email))
	gender := strings.ToLower(strings.TrimSpace(req.Gender))

	if name == "" || email == "" || gender == "" || req.BirthDate.IsZero() {
		return nil, ErrInvalidUserInput
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, ErrInvalidUserInput
	}

	return u.repo.Create(ctx, name, email, gender, req.BirthDate)
}

func normalizePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = defaultUsersPage
	}
	if pageSize < 1 {
		pageSize = defaultUsersPageSize
	}
	if pageSize > maxUsersPageSize {
		pageSize = maxUsersPageSize
	}

	return page, pageSize
}
