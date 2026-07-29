package sqlerr

import (
	"blog_platform/internal/errs"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func HandleError(err error) error {

	// Check if it's already a custom error
	var httpErr *errs.HTTPError
	if errors.As(err, &httpErr) {
		return err
	}

	// Check GORM specific errors
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.NewNotFoundError("Resource not found")
	}

	// Check Postgres driver (pgx) errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique violation
			return errs.NewBadRequestError("A record with this identifier already exists")
		case "23503": // foreign_key_violation
			return errs.NewBadRequestError("The referenced record does not exist")
		case "23502": // not_null_violation
			return errs.NewBadRequestError("A required field is missing")
		default:
			return errs.NewInternalServerError()
		}
	}

	return errs.NewInternalServerError()
}
