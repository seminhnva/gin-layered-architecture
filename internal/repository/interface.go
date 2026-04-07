package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/seminhnva/gin-layered-architecture/internal/db/sqlc"
)

type UserRepository interface {
	GetUsers(ctx context.Context, q sqlc.ListUsersParams) ([]sqlc.ListUsersRow, error)
	GetUsersV2(ctx context.Context, q UserListFilter) ([]sqlc.ListUsersRow, error)
	CountUser(ctx context.Context, countParam sqlc.CountUsersParams) (int64, error)
	Create(ctx context.Context, params sqlc.CreateUserParams) (sqlc.CreateUserRow, error)
	FindByUUID(ctx context.Context, ID uuid.UUID) (sqlc.FindUserByIDRow, error)
	Update(ctx context.Context, params sqlc.UpdateUserByIDParams) (sqlc.UpdateUserByIDRow, error)
	SoftDelete(ctx context.Context, ID uuid.UUID) error
}
type UserListFilter struct {
	Page   int32
	Limit  int32
	SortBy string
	Order  string
	Search *string
}

type AuthRepository interface {
}
