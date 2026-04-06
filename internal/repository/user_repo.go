package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/seminhnva/gin-layered-architecture/internal/api/validation"
	"github.com/seminhnva/gin-layered-architecture/internal/common/domainerror"
	"github.com/seminhnva/gin-layered-architecture/internal/db/sqlc"
)

type UserRepo struct {
	queries *sqlc.Queries
}

func NewUserRepo(queries *sqlc.Queries) UserRepository {
	return &UserRepo{
		queries: queries,
	}
}
func (ur *UserRepo) GetUsers(ctx context.Context, params sqlc.ListUsersParams) ([]sqlc.User, error) {
	users, error := ur.queries.ListUsers(ctx, params)
	return users, error
}

func (ur *UserRepo) CountUser(ctx context.Context, countParam sqlc.CountUsersParams) (int64, error) {
	total, err := ur.queries.CountUsers(ctx, sqlc.CountUsersParams{
		Search: countParam.Search,
	})
	return total, err
}

func (ur *UserRepo) FindByUUID(ctx context.Context, ID uuid.UUID) (sqlc.User, error) {
	user, err := ur.queries.FindUserByID(ctx, sqlc.FindUserByIDParams{
		UserID: ID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, domainerror.ErrUserNotFound
		}
		return sqlc.User{}, validation.MapDBError(err)
	}
	return user, nil
}

func (ur *UserRepo) Create(ctx context.Context, params sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := ur.queries.CreateUser(ctx, params)
	if err != nil {
		return sqlc.User{}, validation.MapDBError(err)
	}
	return user, nil
}
func (ur *UserRepo) Update(ctx context.Context, params sqlc.UpdateUserByIDParams) (sqlc.User, error) {
	user, err := ur.queries.UpdateUserByID(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, domainerror.ErrUserNotFound
		}
		return sqlc.User{}, validation.MapDBError(err)
	}
	return user, nil
}

func (ur *UserRepo) SoftDelete(ctx context.Context, ID uuid.UUID) error {
	rows, err := ur.queries.SoftDeleteUser(ctx, sqlc.SoftDeleteUserParams{
		UserID: ID,
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return domainerror.ErrUserNotFound
	}

	return nil
}
