package models

import (
	"encoding/json"
	"time"
)

type UUID string

func (id UUID) String() string {
	return string(id)
}

type User struct {
	ID         UUID      `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	gender     string    `json:"gender"`
	birth_date time.Time `json:"birth_date"`
}

type CreateUserRequest struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

func NewUser(id UUID, name, email, gender string, birthDate time.Time) User {
	return User{
		ID:         id,
		Name:       name,
		Email:      email,
		gender:     gender,
		birth_date: birthDate,
	}
}

func (u User) Gender() string {
	return u.gender
}

func (u *User) SetGender(gender string) {
	u.gender = gender
}

func (u User) BirthDate() time.Time {
	return u.birth_date
}

func (u *User) SetBirthDate(birthDate time.Time) {
	u.birth_date = birthDate
}

func (u User) MarshalJSON() ([]byte, error) {
	type userJSON struct {
		ID        UUID      `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Gender    string    `json:"gender"`
		BirthDate time.Time `json:"birth_date"`
	}

	return json.Marshal(userJSON{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Gender:    u.gender,
		BirthDate: u.birth_date,
	})
}

func (u *User) UnmarshalJSON(data []byte) error {
	type userJSON struct {
		ID        UUID      `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Gender    string    `json:"gender"`
		BirthDate time.Time `json:"birth_date"`
	}

	var parsed userJSON
	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}

	u.ID = parsed.ID
	u.Name = parsed.Name
	u.Email = parsed.Email
	u.gender = parsed.Gender
	u.birth_date = parsed.BirthDate

	return nil
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
