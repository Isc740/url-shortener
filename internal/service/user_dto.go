package service

import (
	"time"
)

type UserResponse struct {
	ID        int64     `json:"id"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}
