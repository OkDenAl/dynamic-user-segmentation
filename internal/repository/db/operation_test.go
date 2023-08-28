package db

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/pkg/logging"
	"errors"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRepository_GetAllOperationsSortedByUserId(t *testing.T) {
	type args struct {
		ctx   context.Context
		year  int
		month int
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)
	testTime := time.Now()
	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantRes      []entity.Operation
		wantErr      bool
	}{
		{
			name: "OK", args: args{ctx: context.Background(), month: 8, year: 2023},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"user_id", "segment_name", "type", "created_at"}).
					AddRows([]any{int64(1), "test", entity.AddOperation, testTime}).
					AddRows([]any{int64(1), "avito", entity.DelOperation, testTime})

				m.ExpectQuery("SELECT user_id,segment_name,type,created_at").
					WithArgs(args.month, args.year).
					WillReturnRows(rows)
			}, wantRes: []entity.Operation{{UserId: int64(1), SegmentName: "test", Type: entity.AddOperation, CreatedAt: testTime},
			{UserId: int64(1), SegmentName: "avito", Type: entity.DelOperation, CreatedAt: testTime}}, wantErr: false,
		},
		{
			name: "Unexpected error", args: args{ctx: context.Background(), month: 8, year: 2023},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("SELECT user_id,segment_name,type,created_at").
					WithArgs(args.month, args.year).
					WillReturnError(errors.New("some error"))
			}, wantRes: nil, wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)
			operationRepoMock := NewOperationRepo(poolMock, logging.NewForMocks())
			got, err := operationRepoMock.GetAllOperationsSortedByUserId(tc.args.ctx, tc.args.month, tc.args.year)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, got, len(tc.wantRes))
			assert.Equal(t, got, tc.wantRes)
			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
