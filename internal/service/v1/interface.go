package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/seminhnva/gin-layered-architecture/internal/db/sqlc"
	v1dto "github.com/seminhnva/gin-layered-architecture/internal/dto/v1"
)

type UserService interface {
	GetUsers(ctx context.Context, query v1dto.ListUsersQuery) (*v1dto.PaginationResponse[v1dto.UserDTO], error)
	CreateUser(ctx context.Context, req v1dto.CreateUserRequest) (sqlc.User, error)
	GetUserByUUID(ctx context.Context, ID uuid.UUID) (sqlc.User, error)
	UpdateUser(ctx context.Context, params sqlc.UpdateUserByIDParams) (sqlc.User, error)
	DeleteUser(ctx context.Context, ID uuid.UUID) error
}
