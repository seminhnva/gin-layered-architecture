package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seminhnva/gin-layered-architecture/internal/api/validation"
	"github.com/seminhnva/gin-layered-architecture/internal/common/domainerror"
	"github.com/seminhnva/gin-layered-architecture/internal/db/sqlc"
	"github.com/seminhnva/gin-layered-architecture/internal/utils"
)

type UserRepo struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewUserRepo(db *pgxpool.Pool, queries *sqlc.Queries) UserRepository {
	return &UserRepo{
		db:      db,
		queries: queries,
	}
}
func (ur *UserRepo) GetUsers(ctx context.Context, q sqlc.ListUsersParams) ([]sqlc.ListUsersRow, error) {
	users, error := ur.queries.ListUsers(ctx, q)
	return users, error
}
func (ur *UserRepo) GetUsersV2(ctx context.Context, q UserListFilter) ([]sqlc.ListUsersRow, error) {
	sortBy := utils.NormalizeSortBy(q.SortBy)
	order := utils.NormalizeOrder(q.Order)

	query := `
			SELECT
				user_id,
				user_name,
				name,
				email
			FROM users
			WHERE deleted_at IS NULL
			`

	args := make([]any, 0, 3)
	argPos := 1

	if q.Search != nil && *q.Search != "" {
		query += fmt.Sprintf(`
			AND (
				user_name ILIKE $%d OR
				email ILIKE $%d OR
				name ILIKE $%d
			)`,
			argPos, argPos, argPos)

		args = append(args, "%"+*q.Search+"%")
		argPos++
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)

	offset := (q.Page - 1) * q.Limit
	args = append(args, q.Limit, offset)

	rows, err := ur.db.Query(ctx, query, args...)
	if err != nil {
		return nil, validation.MapDBError(err)
	}
	defer rows.Close()

	users := make([]sqlc.ListUsersRow, 0)
	for rows.Next() {
		var user sqlc.ListUsersRow

		if err := rows.Scan(
			&user.UserID,
			&user.UserName,
			&user.Name,
			&user.Email,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepo) CountUser(ctx context.Context, countParam sqlc.CountUsersParams) (int64, error) {
	total, err := ur.queries.CountUsers(ctx, sqlc.CountUsersParams{
		Search: countParam.Search,
	})
	return total, err
}

func (ur *UserRepo) FindByUUID(ctx context.Context, ID uuid.UUID) (sqlc.FindUserByIDRow, error) {
	user, err := ur.queries.FindUserByID(ctx, sqlc.FindUserByIDParams{
		UserID: ID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.FindUserByIDRow{}, domainerror.ErrUserNotFound
		}
		return sqlc.FindUserByIDRow{}, validation.MapDBError(err)
	}
	return user, nil
}

func (ur *UserRepo) Create(ctx context.Context, params sqlc.CreateUserParams) (sqlc.CreateUserRow, error) {
	user, err := ur.queries.CreateUser(ctx, params)
	if err != nil {
		return sqlc.CreateUserRow{}, validation.MapDBError(err)
	}
	return user, nil
}
func (ur *UserRepo) Update(ctx context.Context, params sqlc.UpdateUserByIDParams) (sqlc.UpdateUserByIDRow, error) {
	user, err := ur.queries.UpdateUserByID(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.UpdateUserByIDRow{}, domainerror.ErrUserNotFound
		}
		return sqlc.UpdateUserByIDRow{}, validation.MapDBError(err)
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
