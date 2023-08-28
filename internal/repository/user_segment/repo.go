package user_segment

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/pkg/logging"
	"dynamic-user-segmentation/pkg/postgres"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name=Repository --output=../../mocks/repo/usersegmentrepo --outpkg=usersegmentrepo_mocks

type Repository interface {
	CreateMultSegsForOneUser(ctx context.Context, userId int64, segments []entity.Segment, ttl entity.TTL) error
	CreateOneSegForMultUsers(ctx context.Context, userIds []int64, segment entity.Segment) error
	DeleteSegmentsFromSpecUser(ctx context.Context, userId int64, segments []entity.Segment) error
	GetAllUserSegments(ctx context.Context, userId int64) ([]entity.Segment, error)
}

type repo struct {
	conn postgres.PgxPool
	l    logging.Logger
}

func New(conn postgres.PgxPool, l logging.Logger) Repository {
	return &repo{conn: conn, l: l}
}

func (r *repo) CreateMultSegsForOneUser(ctx context.Context, userId int64, segments []entity.Segment, ttl entity.TTL) error {
	var (
		insertUsersSegmentBatch = &pgx.Batch{}
		insertUsersSegmentQ     string
		insertOperationBatch    = &pgx.Batch{}
		insertOperationQ        = `INSERT INTO operations(user_id,segment_name,type,created_at) VALUES ($1,$2,$3,$4)`
	)
	intervalString := fmt.Sprintf("P%dY%dM%dDT%dH%dM%dS", ttl.Years, ttl.Months, ttl.Days, ttl.Hours, ttl.Minutes, ttl.Seconds)
	if intervalString == "P0Y0M0DT0H0M0S" {
		insertUsersSegmentQ = `INSERT INTO users_segments (user_id,segment_name) VALUES ($1,$2)`
	} else {
		insertUsersSegmentQ = fmt.Sprintf("INSERT INTO users_segments (user_id,segment_name,expires_at) VALUES ($1,$2,now()+interval '%s')", intervalString)
	}

	tx, err := r.conn.Begin(ctx)
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.CreateMultSegsForOneUser - r.conn.Begin - %w", err))
		return repository.SqlErrorWrapper(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, segment := range segments {
		insertUsersSegmentBatch.Queue(insertUsersSegmentQ, userId, segment.Name)
		insertOperationBatch.Queue(insertOperationQ, userId, segment.Name, entity.AddOperation, time.Now())
	}
	err = tx.SendBatch(ctx, insertUsersSegmentBatch).Close()
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.CreateMultSegsForOneUser - r.conn.SendBatch.Close - insertUsersSegmentBatch - %w", err))
		return repository.SqlErrorWrapper(err)
	}
	err = tx.SendBatch(ctx, insertOperationBatch).Close()
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.CreateMultSegsForOneUser - r.conn.SendBatch.Close - insertOperationBatch - %w", err))
		return repository.SqlErrorWrapper(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.CreateMultSegsForOneUser - tx.Commit - %w", err))
		return repository.SqlErrorWrapper(err)
	}

	return nil
}

func (r *repo) CreateOneSegForMultUsers(ctx context.Context, userIds []int64, segment entity.Segment) error {
	var (
		insertUsersSegmentBatch = &pgx.Batch{}
		insertUsersSegmentQ     = `INSERT INTO users_segments (user_id,segment_name) VALUES ($1,$2)`
		insertOperationBatch    = &pgx.Batch{}
		insertOperationQ        = `INSERT INTO operations(user_id,segment_name,type,created_at) VALUES ($1,$2,$3,$4)`
	)
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.CreateOneSegForMultUsers - r.conn.Begin - %w", err))
		return repository.SqlErrorWrapper(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, userId := range userIds {
		insertUsersSegmentBatch.Queue(insertUsersSegmentQ, userId, segment.Name)
		insertOperationBatch.Queue(insertOperationQ, userId, segment.Name, entity.AddOperation, time.Now())
	}
	err = tx.SendBatch(ctx, insertUsersSegmentBatch).Close()
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.CreateOneSegForMultUsers - r.conn.SendBatch.Close - insertUsersSegmentBatch  - %w", err))
		return repository.SqlErrorWrapper(err)
	}
	err = tx.SendBatch(ctx, insertOperationBatch).Close()
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.CreateOneSegForMultUsers - r.conn.SendBatch.Close - insertOperationBatch - %w", err))
		return repository.SqlErrorWrapper(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.CreateOneSegForMultUsers - tx.Commit - %w", err))
		return repository.SqlErrorWrapper(err)
	}

	return nil
}

func (r *repo) DeleteSegmentsFromSpecUser(ctx context.Context, userId int64, segments []entity.Segment) error {
	var (
		deleteBatch          = &pgx.Batch{}
		deleteQ              = `DELETE FROM users_segments WHERE user_id=$1 AND segment_name=$2`
		insertOperationBatch = &pgx.Batch{}
		insertOperationQ     = `INSERT INTO operations(user_id,segment_name,type,created_at) VALUES ($1,$2,$3,$4)`
	)
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.DeleteSegmentsFromSpecUser - r.conn.Begin - %w", err))
		return repository.SqlErrorWrapper(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, segment := range segments {
		deleteBatch.Queue(deleteQ, userId, segment.Name)
		insertOperationBatch.Queue(insertOperationQ, userId, segment.Name, entity.DelOperation, time.Now())
	}
	err = tx.SendBatch(ctx, deleteBatch).Close()
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.DeleteSegmentsFromSpecUser - r.conn.SendBatch.Close - deleteBatch - %w", err))
		return repository.SqlErrorWrapper(err)
	}
	err = tx.SendBatch(ctx, insertOperationBatch).Close()
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.DeleteSegmentsFromSpecUser - r.conn.SendBatch.Close - insertOperationBatch - %w", err))
		return repository.SqlErrorWrapper(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.DeleteSegmentsFromSpecUser - tx.Commit - %w", err))
		return repository.SqlErrorWrapper(err)
	}

	return nil
}

func (r *repo) GetAllUserSegments(ctx context.Context, userId int64) ([]entity.Segment, error) {
	q := `SELECT segment_name FROM users_segments WHERE user_id=$1 AND(expires_at > now() OR expires_at IS NULL)`
	rows, err := r.conn.Query(ctx, q, userId)
	if err != nil {
		r.l.Error(fmt.Errorf("user_segment.GetAllUserSegments - r.conn.Query - %w", err))
		return nil, repository.SqlErrorWrapper(err)
	}
	defer rows.Close()
	usersSegments := make([]entity.Segment, 0)
	for rows.Next() {
		var curSegment entity.Segment
		err = rows.Scan(&curSegment.Name)
		if err != nil {
			r.l.Error(fmt.Errorf("user_segment.GetAllUserSegments - rows.Scan - %w", err))
			return nil, repository.SqlErrorWrapper(err)
		}
		usersSegments = append(usersSegments, curSegment)
	}
	return usersSegments, nil
}
