package v1dto

import (
	"net/url"

	"github.com/google/uuid"
	"github.com/seminhnva/gin-layered-architecture/internal/db/sqlc"
)

var ListUsersAllowedQueryKeys = map[string]struct{}{
	"page":    {},
	"limit":   {},
	"sort_by": {},
	"order":   {},
	"search":  {},
}

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
	Page   int32   `form:"page" binding:"omitempty,min=1"`
	Limit  int32   `form:"limit" binding:"omitempty,min=1,max=100"`
	SortBy string  `form:"sort_by" binding:"omitempty,oneof=created_at name email"`
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
	allowed := map[string]bool{
		"name":       true,
		"email":      true,
		"created_at": true,
	}
	if !allowed[q.SortBy] {
		q.SortBy = "created_at"
	}
}

func (q ListUsersQuery) ToQueryValues() url.Values {
	values := url.Values{}

	if q.SortBy != "" {
		values.Set("sort_by", q.SortBy)
	}

	if q.Order != "" {
		values.Set("order", q.Order)
	}

	if q.Search != nil && *q.Search != "" {
		values.Set("search", *q.Search)
	}

	return values
}

type UserListItem struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type PaginationResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int32 `json:"page"`
	Limit      int32 `json:"limit"`
	TotalPages int32 `json:"total_pages"`
}

func toUserDTO(id uuid.UUID, name, username, email string) *UserDTO {
	return &UserDTO{
		UUID:     id.String(),
		Name:     name,
		Username: username,
		Email:    email,
	}
}

func ToUserDTOFromList(u sqlc.ListUsersRow) *UserDTO {
	return toUserDTO(u.UserID, u.Name, u.UserName, u.Email)
}
func ToUserDTOFromCreate(u sqlc.CreateUserRow) *UserDTO {
	return toUserDTO(u.UserID, u.Name, u.UserName, u.Email)
}
func ToUserDTOFromFind(u sqlc.FindUserByIDRow) *UserDTO {
	return toUserDTO(u.UserID, u.Name, u.UserName, u.Email)
}
func ToUserDTOFromUpdate(u sqlc.UpdateUserByIDRow) *UserDTO {
	return toUserDTO(u.UserID, u.Name, u.UserName, u.Email)
}
