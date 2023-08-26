package repository

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrAlreadyExists = errors.New("this data already exists")
	ErrUnknownData   = errors.New("this data is unknown")
	ErrUnknown       = errors.New("unknown database error")
)

func SqlErrorWrapper(err error) error {
	if pgError := err.(*pgconn.PgError); errors.Is(err, pgError) {
		switch pgError.Code {
		case "23505":
			return ErrAlreadyExists
		case "23503":
			return ErrUnknownData
		default:
			return ErrUnknown
		}
	}
	return err
}
