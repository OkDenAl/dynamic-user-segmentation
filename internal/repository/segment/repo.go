package segment

import (
	"context"
	"dynamic-user-segmentation/pkg/postgres"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrAlreadyExists = errors.New("segment with this name already exists")
)

type Repository interface {
	Create(ctx context.Context, name string) error
	Delete(ctx context.Context, name string) error
}

type repo struct {
	conn postgres.PgxPool
}

func New(conn *pgxpool.Pool) Repository {
	return &repo{conn: conn}
}

func (r *repo) Create(ctx context.Context, name string) error {
	q := `INSERT INTO segments (name) VALUES ($1)`
	_, err := r.conn.Exec(ctx, q, name)
	if err != nil {
		if pgError := err.(*pgconn.PgError); errors.Is(err, pgError) {
			switch pgError.Code {
			case "23505":
				return ErrAlreadyExists
			default:
				return fmt.Errorf("r.conn.Exec - %w", err)
			}
		}
	}
	return nil
}

func (r *repo) Delete(ctx context.Context, name string) error {
	q := `DELETE FROM segments WHERE name=$1`
	_, err := r.conn.Exec(ctx, q, name)
	if err != nil {
		return fmt.Errorf("r.conn.Exec - %w", err)
	}
	return nil
}
