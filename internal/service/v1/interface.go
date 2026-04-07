package v1service

import (
	"context"

	"github.com/google/uuid"
	"github.com/seminhnva/gin-layered-architecture/internal/db/sqlc"
	v1dto "github.com/seminhnva/gin-layered-architecture/internal/dto/v1"
)

type UserService interface {
	GetUsers(ctx context.Context, param v1dto.ListUsersQuery) (*v1dto.PaginationResponse[v1dto.UserDTO], error)
	CreateUser(ctx context.Context, req v1dto.CreateUserRequest) (sqlc.CreateUserRow, error)
	GetUserByUUID(ctx context.Context, ID uuid.UUID) (sqlc.FindUserByIDRow, error)
	UpdateUser(ctx context.Context, params sqlc.UpdateUserByIDParams) (sqlc.UpdateUserByIDRow, error)
	DeleteUser(ctx context.Context, ID uuid.UUID) error
}
