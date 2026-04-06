package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/seminhnva/gin-layered-architecture/internal/db/sqlc"
)

type UserRepository interface {
	GetUsers(ctx context.Context, params sqlc.ListUsersParams) ([]sqlc.User, error)
	CountUser(ctx context.Context, params sqlc.CountUsersParams) (int64, error)
	Create(ctx context.Context, params sqlc.CreateUserParams) (sqlc.User, error)
	FindByUUID(ctx context.Context, ID uuid.UUID) (sqlc.User, error)
	Update(ctx context.Context, params sqlc.UpdateUserByIDParams) (sqlc.User, error)
	SoftDelete(ctx context.Context, ID uuid.UUID) error
}
