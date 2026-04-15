package validation

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/seminhnva/gin-layered-architecture/internal/common/domainerror"
)

const (
	PgUniqueViolation     = "23505"
	PgForeignKeyViolation = "23503"
	PgNotNullViolation    = "23502"
)

func MapDBError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case PgUniqueViolation:
		switch pgErr.ConstraintName {
		case "users_email_key":
			return domainerror.ErrEmailAlreadyExists
		case "users_user_name_key":
			return domainerror.ErrUserNameAlreadyExists
		default:
			return err
		}

	case PgNotNullViolation:
		return domainerror.ErrRequired

	case PgForeignKeyViolation:
		return domainerror.ErrInvalidReference

	default:
		return err
	}
}
