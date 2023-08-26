package segment

import (
	"context"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/pkg/logging"
	"dynamic-user-segmentation/pkg/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, name string) error
	Delete(ctx context.Context, name string) error
}

type repo struct {
	conn postgres.PgxPool
	l    logging.Logger
}

func New(conn *pgxpool.Pool, l logging.Logger) Repository {
	return &repo{conn: conn, l: l}
}

func (r *repo) Create(ctx context.Context, name string) error {
	q := `INSERT INTO segments (name) VALUES ($1)`
	_, err := r.conn.Exec(ctx, q, name)
	if err != nil {
		r.l.Errorf("segment.Create - r.conn.Exec %w", err)
		return repository.SqlErrorWrapper(err)
	}
	return nil
}

func (r *repo) Delete(ctx context.Context, name string) error {
	q := `DELETE FROM segments WHERE name=$1`
	_, err := r.conn.Exec(ctx, q, name)
	if err != nil {
		r.l.Errorf("segment.Delete - r.conn.Exec %w", err)
		return repository.SqlErrorWrapper(err)
	}
	return nil
}
