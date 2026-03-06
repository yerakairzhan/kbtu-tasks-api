package models

import "time"

type UUID string

func (id UUID) String() string {
	return string(id)
}

type User struct {
	ID        UUID      `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Gender    string    `json:"gender" db:"gender"`
	BirthDate time.Time `json:"birth_date" db:"birth_date"`
}

type CreateUserRequest struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

type ListUsersInput struct {
	Page     int
	PageSize int
}

type PaginatedResponse struct {
	Data       []User `json:"data"`
	TotalCount int    `json:"totalCount"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
}

func NewPaginatedResponse(data []User, totalCount, page, pageSize int) PaginatedResponse {
	return PaginatedResponse{
		Data:       data,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}
}
