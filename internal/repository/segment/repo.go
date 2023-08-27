package segment

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/pkg/logging"
	"dynamic-user-segmentation/pkg/postgres"
	"fmt"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=Repository --output=../../mocks/repo/segmentrepo --outpkg=segmentrepo_mocks

type Repository interface {
	Create(ctx context.Context, segment entity.Segment) error
	Delete(ctx context.Context, segment entity.Segment) error
}

type repo struct {
	conn postgres.PgxPool
	l    logging.Logger
}

func New(conn postgres.PgxPool, l logging.Logger) Repository {
	return &repo{conn: conn, l: l}
}

func (r *repo) Create(ctx context.Context, segment entity.Segment) error {
	q := `INSERT INTO segments (name) VALUES ($1)`
	_, err := r.conn.Exec(ctx, q, segment.Name)
	if err != nil {
		r.l.Error(fmt.Errorf("segment.Create - r.conn.Exec - %w", err))
		return repository.SqlErrorWrapper(err)
	}
	return nil
}

func (r *repo) Delete(ctx context.Context, segment entity.Segment) error {
	q := `DELETE FROM segments WHERE name=$1`
	_, err := r.conn.Exec(ctx, q, segment.Name)
	if err != nil {
		r.l.Error(fmt.Errorf("segment.Delete - r.conn.Exec - %w", err))
		return repository.SqlErrorWrapper(err)
	}
	return nil
}
