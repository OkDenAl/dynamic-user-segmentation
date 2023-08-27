package user_segment

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/pkg/logging"
	"dynamic-user-segmentation/pkg/postgres"
	"fmt"
	"github.com/jackc/pgx/v5"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=Repository --output=../../mocks/repo/usersegmentrepo --outpkg=usersegmentrepo_mocks

type Repository interface {
	CreateMultSegsForOneUser(ctx context.Context, userId int64, segments []string, ttl entity.TTL) error
	CreateOneSegForMultUsers(ctx context.Context, userIds []int64, segment string) error
	DeleteSegmentsFromSpecUser(ctx context.Context, userId int64, segments []string) error
	GetAllUserSegments(ctx context.Context, userId int64) ([]string, error)
}

type repo struct {
	conn postgres.PgxPool
	l    logging.Logger
}

func New(conn postgres.PgxPool, l logging.Logger) Repository {
	return &repo{conn: conn, l: l}
}

func (r *repo) CreateMultSegsForOneUser(ctx context.Context, userId int64, segments []string, ttl entity.TTL) error {
	var (
		batch = &pgx.Batch{}
		q     string
	)
	intervalString := fmt.Sprintf("P%dY%dM%dDT%dH%dM%dS", ttl.Years, ttl.Months, ttl.Days, ttl.Hours, ttl.Minutes, ttl.Seconds)
	if intervalString == "P0Y0M0DT0H0M0S" {
		q = `INSERT INTO users_segments (user_id,segment_name) VALUES ($1,$2)`
	} else {
		q = fmt.Sprintf("INSERT INTO users_segments (user_id,segment_name,expires_at) VALUES ($1,$2,now()+interval '%s')", intervalString)
	}
	for _, segment := range segments {
		batch.Queue(q, userId, segment)
	}
	err := r.conn.SendBatch(ctx, batch).Close()
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.CreateMultSegsForOneUser - r.conn.SendBatch.Close - %w", err))
		return repository.SqlErrorWrapper(err)
	}
	return nil
}

func (r *repo) CreateOneSegForMultUsers(ctx context.Context, userIds []int64, segment string) error {
	batch := &pgx.Batch{}
	q := `INSERT INTO users_segments (user_id,segment_name) VALUES ($1,$2)`
	for _, userId := range userIds {
		batch.Queue(q, userId, segment)
	}
	err := r.conn.SendBatch(ctx, batch).Close()
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.CreateOneSegForMultUsers - r.conn.SendBatch.Close - %w", err))
		return repository.SqlErrorWrapper(err)
	}
	return nil
}

func (r *repo) DeleteSegmentsFromSpecUser(ctx context.Context, userId int64, segments []string) error {
	batch := &pgx.Batch{}
	q := `DELETE FROM users_segments WHERE user_id=$1 AND segment_name=$2`
	for _, segment := range segments {
		batch.Queue(q, userId, segment)
	}
	err := r.conn.SendBatch(ctx, batch).Close()
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.DeleteSegmentsFromSpecUser - r.conn.SendBatch.Close - %w", err))
		return repository.SqlErrorWrapper(err)
	}
	return nil
}

func (r *repo) GetAllUserSegments(ctx context.Context, userId int64) ([]string, error) {
	q := `SELECT segment_name FROM users_segments WHERE user_id=$1 AND(expires_at > now() OR expires_at IS NULL)`
	rows, err := r.conn.Query(ctx, q, userId)
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.GetAllUserSegments - r.conn.Query - %w", err))
		return nil, repository.SqlErrorWrapper(err)
	}
	defer rows.Close()
	usersSegments := make([]string, 0)
	for rows.Next() {
		var curSegment string
		err = rows.Scan(&curSegment)
		if err != nil {
			r.l.Error(fmt.Errorf("user_segment.GetAllUserSegments - rows.Scan - %w", err))
			return nil, repository.SqlErrorWrapper(err)
		}
		usersSegments = append(usersSegments, curSegment)
	}
	return usersSegments, nil
}
