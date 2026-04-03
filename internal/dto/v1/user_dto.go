package v1dto

import (
	"github.com/google/uuid"
	"github.com/seminhnva/gin-layered-architecture/internal/db/sqlc"
)

type UserDTO struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type CreateUserRequest struct {
	UserName string `json:"userName" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (c *CreateUserRequest) ToCreateUserParams() sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		UserName:     c.UserName,
		Name:         c.Name,
		Email:        c.Email,
		PasswordHash: c.Password,
	}
}

type UpdateUserRequest struct {
	UserName *string `json:"userName" binding:"omitempty,required"`
	Name     *string `json:"name" binding:"omitempty,required"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,required,min=6"`
}

func (u *UpdateUserRequest) ToUpdateUserParams(ID uuid.UUID) sqlc.UpdateUserByIDParams {
	return sqlc.UpdateUserByIDParams{
		UserID:       ID,
		UserName:     u.UserName,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.Password,
	}
}

type GetUserIdParams struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}

type ListUsersQuery struct {
	Page   int     `form:"page" binding:"omitempty,min=1"`
	Limit  int     `form:"limit" binding:"omitempty,min=1,max=100"`
	SortBy string  `form:"sort_by" binding:"omitempty,oneof=created_at updated_at name email"`
	Order  string  `form:"order" binding:"omitempty,oneof=asc desc"`
	Search *string `form:"search" binding:"omitempty,max=100"`
}

func (q *ListUsersQuery) Normalize() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Order == "" {
		q.Order = "asc"
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
}

func ToUserResponse(user sqlc.User) *UserDTO {
	return &UserDTO{
		UUID:     user.UserID.String(),
		Name:     user.Name,
		Username: user.UserName,
		Email:    user.Email,
	}
}
