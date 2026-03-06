package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"tasks_assignment/internal/models"
	"tasks_assignment/internal/usecase"
)

type UserHandler struct {
	usecase UserUsecase
}

type UserUsecase interface {
	GetUsers(ctx context.Context, page, pageSize int) (models.PaginatedResponse, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error)
}

func NewUserHandler(uc UserUsecase) *UserHandler {
	return &UserHandler{usecase: uc}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id != "" {
		user, err := h.usecase.GetUserByID(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, usecase.ErrInvalidUserInput):
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid id"})
			case errors.Is(err, usecase.ErrUserNotFound):
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(models.ErrorResponse{Error: "user not found"})
			default:
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models.ErrorResponse{Error: "internal error"})
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
		return
	}

	page := 1
	if raw := r.URL.Query().Get("page"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid page"})
			return
		}
		page = parsed
	}

	pageSize := 10
	if raw := r.URL.Query().Get("page_size"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid page_size"})
			return
		}
		pageSize = parsed
	}

	resp, err := h.usecase.GetUsers(r.Context(), page, pageSize)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "internal error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid request body"})
		return
	}

	user, err := h.usecase.CreateUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidUserInput) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid user data"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "internal error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
