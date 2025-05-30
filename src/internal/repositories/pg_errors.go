package repositories

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrRecordNotFound = errors.New("record not found")
var ErrDuplicate = errors.New("duplicated key not allowed")
var ErrMissingWhereClause = errors.New("missing where clause in delete or update operation")
var ErrInvalidField = errors.New("invalid field")
var ErrInvalidValue = errors.New("invalid value")
var ErrOffsetMustBePositive = errors.New("offset must be positive")
var ErrLimitMustBePositive = errors.New("limit must be positive")

func parsePgError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return context.DeadlineExceeded
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "P0002":
			return ErrRecordNotFound
		case "23505":
			return ErrDuplicate
		case "2D000":
			return ErrMissingWhereClause
		case "42703":
			return ErrInvalidField
		case "22P02", "22007", "22003":
			return ErrInvalidValue
		case "2201X":
			return ErrOffsetMustBePositive
		case "2201W":
			return ErrLimitMustBePositive
		default:
			return err
		}
	}
	return err
}
