package user

import (
	"context"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/pkg/logging"
	"dynamic-user-segmentation/pkg/postgres"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"math"
)

type Repository interface {
	GetNPercentOfUsersIDs(ctx context.Context, percent float64) ([]int64, error)
}

type repo struct {
	conn postgres.PgxPool
	l    logging.Logger
}

func New(conn *pgxpool.Pool, l logging.Logger) Repository {
	return &repo{conn: conn, l: l}
}

func (r *repo) GetNPercentOfUsersIDs(ctx context.Context, percent float64) ([]int64, error) {
	countOfUsers, err := r.getCountOfUsers(ctx)
	if err != nil {
		r.l.Error(fmt.Errorf("user.GetNPercentOfUsersIDs - r.getCountOfUsers - %w", err))
		return nil, repository.SqlErrorWrapper(err)
	}
	q := `SELECT id FROM users LIMIT $1 OFFSET floor(random() * $2);`
	rows, err := r.conn.Query(ctx, q, math.Round(percent/100*float64(countOfUsers)), countOfUsers)
	if err != nil {
		r.l.Error(fmt.Errorf("user.GetNPercentOfUsersIDs - r.conn.Query - %w", err))
		return nil, repository.SqlErrorWrapper(err)
	}
	ids := make([]int64, 0)
	for rows.Next() {
		var curId int64
		err = rows.Scan(&curId)
		if err != nil {
			r.l.Error(fmt.Errorf("user.GetNPercentOfUsersIDs - rows.Scan - %w", err))
			return nil, repository.SqlErrorWrapper(err)
		}
		ids = append(ids, curId)
	}
	return ids, nil
}

func (r *repo) getCountOfUsers(ctx context.Context) (int, error) {
	var countOfUsers int
	q := `SELECT count(*) FROM users`
	err := r.conn.QueryRow(ctx, q).Scan(&countOfUsers)
	if err != nil {
		return 0, fmt.Errorf("r.conn.QueryRow.Scan - %w", err)
	}
	return countOfUsers, nil
}
