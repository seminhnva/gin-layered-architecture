package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seminhnva/gin-layered-architecture/internal/db/sqlc"
)

type AuthRepo struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewAuthRepo(db *pgxpool.Pool, queries *sqlc.Queries) AuthRepository {
	return &AuthRepo{
		db:      db,
		queries: queries,
	}
}

func (ar *AuthRepo) GetByEmail(ctx context.Context, params sqlc.GetByEmailParams) (sqlc.GetByEmailRow, error) {
	userInfor, err := ar.queries.GetByEmail(ctx, params)
	return userInfor, err
}
