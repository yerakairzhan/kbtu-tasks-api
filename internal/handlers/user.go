package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"tasks_assignment/internal/models"
	"tasks_assignment/internal/usecase"
)

type UserHandler struct {
	usecase UserUsecase
}

type UserUsecase interface {
	GetUsers(ctx context.Context, input models.ListUsersInput) (models.PaginatedResponse, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error)
	GetCommonFriends(ctx context.Context, user1, user2 string) ([]models.User, error)
}

func NewUserHandler(uc UserUsecase) *UserHandler {
	return &UserHandler{usecase: uc}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if id := r.URL.Query().Get("id"); id != "" {
		h.getUserByID(w, r, id)
		return
	}

	page, err := parsePositiveInt(r.URL.Query().Get("page"), "page", 1)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	pageSize, err := parsePositiveInt(r.URL.Query().Get("page_size"), "page_size", 10)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	filters, err := buildUserFilters(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	orderBy := r.URL.Query().Get("order_by")

	input := models.ListUsersInput{
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
		Filters:  filters,
	}

	resp, err := h.usecase.GetUsers(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.CreateUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.usecase.CreateUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidUserInput) {
			writeError(w, http.StatusBadRequest, "invalid user data")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) getUserByID(w http.ResponseWriter, r *http.Request, id string) {
	user, err := h.usecase.GetUserByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidUserInput):
			writeError(w, http.StatusBadRequest, "invalid id")
		case errors.Is(err, usecase.ErrUserNotFound):
			writeError(w, http.StatusNotFound, "user not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CommonFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user1 := strings.TrimSpace(r.URL.Query().Get("user1_id"))
	user2 := strings.TrimSpace(r.URL.Query().Get("user2_id"))
	friends, err := h.usecase.GetCommonFriends(r.Context(), user1, user2)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidUserInput) {
			writeError(w, http.StatusBadRequest, "invalid user ids")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(friends)
}

func parsePositiveInt(raw, field string, fallback int) (int, error) {
	if raw == "" {
		return fallback, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return 0, errors.New("invalid " + field)
	}

	return value, nil
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}

func buildUserFilters(r *http.Request) (models.UserFilters, error) {
	q := r.URL.Query()
	filters := models.UserFilters{
		ID:     strings.TrimSpace(q.Get("filter_id")),
		Name:   strings.TrimSpace(q.Get("name")),
		Email:  strings.TrimSpace(q.Get("email")),
		Gender: strings.TrimSpace(q.Get("gender")),
	}

	if birthDateRaw := strings.TrimSpace(q.Get("birth_date")); birthDateRaw != "" {
		birthDate, err := parseBirthDate(birthDateRaw)
		if err != nil {
			return filters, err
		}
		filters.BirthDate = &birthDate
	}

	return filters, nil
}

func parseBirthDate(raw string) (time.Time, error) {
	layouts := []string{time.RFC3339, "2006-01-02"}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, raw); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, errors.New("invalid birth_date")
}
