package dto

import "time"

type UserURI struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100"`
	Email string `json:"email" binding:"required,email,max=255"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name" binding:"omitempty,min=2,max=100"`
	Email *string `json:"email" binding:"omitempty,email,max=255"`
}

type UserResponse struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
