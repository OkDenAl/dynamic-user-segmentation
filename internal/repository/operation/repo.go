package operation

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/pkg/logging"
	"dynamic-user-segmentation/pkg/postgres"
	"fmt"
)

type Repository interface {
	GetAllOperationsSortedByUserId(ctx context.Context, month, year int) ([]entity.Operation, error)
}

type repo struct {
	conn postgres.PgxPool
	l    logging.Logger
}

func New(conn postgres.PgxPool, l logging.Logger) Repository {
	return &repo{conn: conn, l: l}
}

func (r *repo) GetAllOperationsSortedByUserId(ctx context.Context, month, year int) ([]entity.Operation, error) {
	q := `SELECT user_id,segment_name,type,created_at FROM operations WHERE extract
	(month from created_at) = $1 and extract(year from created_at) = $2 ORDER BY user_id`
	rows, err := r.conn.Query(ctx, q, month, year)
	if err != nil {
		r.l.Error(fmt.Errorf("operation.GetAllOperationsSortedByUserId - r.conn.Query - %w", err))
		return nil, repository.SqlErrorWrapper(err)
	}
	operations := make([]entity.Operation, 0)
	for rows.Next() {
		var oper entity.Operation
		err = rows.Scan(&oper.UserId, &oper.SegmentName, &oper.Type, &oper.CreatedAt)
		if err != nil {
			r.l.Error(fmt.Errorf("operation.GetAllOperationsSortedByUserId - r.conn.Query - %w", err))
			return nil, repository.SqlErrorWrapper(err)
		}
		operations = append(operations, oper)
	}
	return operations, nil
}
