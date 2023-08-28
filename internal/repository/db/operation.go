package db

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository/dberrors"
	"dynamic-user-segmentation/pkg/logging"
	"dynamic-user-segmentation/pkg/postgres"
	"fmt"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=OperationRepository --output=../../mocks/repo/operationrepo --outpkg=operationrepo_mocks

type OperationRepository interface {
	GetAllOperationsSortedByUserId(ctx context.Context, month, year int) ([]entity.Operation, error)
}

type operationRepo struct {
	conn postgres.PgxPool
	l    logging.Logger
}

func NewOperationRepo(conn postgres.PgxPool, l logging.Logger) OperationRepository {
	return &operationRepo{conn: conn, l: l}
}

func (r *operationRepo) GetAllOperationsSortedByUserId(ctx context.Context, month, year int) ([]entity.Operation, error) {
	q := `SELECT user_id,segment_name,type,created_at FROM operations WHERE extract
	(month from created_at) = $1 and extract(year from created_at) = $2 ORDER BY user_id`
	rows, err := r.conn.Query(ctx, q, month, year)
	if err != nil {
		r.l.Error(fmt.Errorf("operation.GetAllOperationsSortedByUserId - r.conn.Query - %w", err))
		return nil, dberrors.SqlErrorWrapper(err)
	}
	operations := make([]entity.Operation, 0)
	for rows.Next() {
		var oper entity.Operation
		err = rows.Scan(&oper.UserId, &oper.SegmentName, &oper.Type, &oper.CreatedAt)
		if err != nil {
			r.l.Error(fmt.Errorf("operation.GetAllOperationsSortedByUserId - r.conn.Query - %w", err))
			return nil, dberrors.SqlErrorWrapper(err)
		}
		operations = append(operations, oper)
	}
	return operations, nil
}
