package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seminhnva/gin-layered-architecture/internal/db/sqlc"
)

type AuthRepo struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewAuthRepo(db *pgxpool.Pool, queries *sqlc.Queries) UserRepository {
	return &UserRepo{
		db:      db,
		queries: queries,
	}
}
