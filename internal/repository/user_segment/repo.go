package user_segment

import (
	"context"
	"dynamic-user-segmentation/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, userId int64, segments []string) error
	DeleteSegmentsFromSpecUser(ctx context.Context, userId int64, segments []string) error
	DeleteSpecSegmentFromUsers(ctx context.Context, segmentName string) error
	GetAllUserSegments(ctx context.Context, userId int64) ([]string, error)
}

type repo struct {
	conn postgres.PgxPool
}

func New(conn *pgxpool.Pool) Repository {
	return &repo{conn: conn}
}

func (r *repo) Create(ctx context.Context, userId int64, segments []string) error {
	batch := &pgx.Batch{}
	q := `INSERT INTO users_segments (user_id,segment_name) VALUES ($1,$2)`
	for _, segment := range segments {
		batch.Queue(q, userId, segment)
	}
	res := r.conn.SendBatch(ctx, batch)
	return res.Close()
}

func (r *repo) DeleteSegmentsFromSpecUser(ctx context.Context, userId int64, segments []string) error {
	batch := &pgx.Batch{}
	q := `DELETE FROM users_segments WHERE user_id=$1 AND segment_name=$2`
	for _, segment := range segments {
		batch.Queue(q, userId, segment)
	}
	res := r.conn.SendBatch(ctx, batch)
	return res.Close()
}

func (r *repo) DeleteSpecSegmentFromUsers(ctx context.Context, segmentName string) error {
	q := `DELETE FROM users_segments WHERE segment_name=$1`
	_, err := r.conn.Exec(ctx, q, segmentName)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) GetAllUserSegments(ctx context.Context, userId int64) ([]string, error) {
	q := `SELECT segment_name FROM users_segments WHERE user_id=$1`
	rows, err := r.conn.Query(ctx, q, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	usersSegments := make([]string, 0)
	for rows.Next() {
		var curSegment string
		err = rows.Scan(&curSegment)
		if err != nil {
			return nil, err
		}
		usersSegments = append(usersSegments, curSegment)
	}
	return usersSegments, nil
}
