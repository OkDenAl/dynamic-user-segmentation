package segment

import (
	"context"
	"dynamic-user-segmentation/pkg/logging"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

type args struct {
	ctx         context.Context
	segmentName string
}

type MockBehavior func(m pgxmock.PgxPoolIface, args args)

type testCase struct {
	name         string
	args         args
	mockBehavior MockBehavior
	wantErr      bool
}

func TestRepository_Create(t *testing.T) {
	testCases := []testCase{
		{
			name: "OK", args: args{ctx: context.Background(), segmentName: "test"},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec("INSERT INTO segments").
					WithArgs(args.segmentName).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			}, wantErr: false,
		},
		{
			name: "Already exists", args: args{ctx: context.Background(), segmentName: "test"},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec("INSERT INTO segments").
					WithArgs(args.segmentName).
					WillReturnError(&pgconn.PgError{
						Code: "23505",
					})
			}, wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)
			segmentRepoMock := New(poolMock, logging.NewForMocks())
			err := segmentRepoMock.Create(tc.args.ctx, tc.args.segmentName)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	testCases := []testCase{
		{
			name: "OK", args: args{ctx: context.Background(), segmentName: "test"},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec("DELETE FROM segments").
					WithArgs(args.segmentName).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			}, wantErr: false,
		},
		{
			name: "Unexpected error", args: args{ctx: context.Background(), segmentName: "test"},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec("DELETE FROM segments").
					WithArgs(args.segmentName).
					WillReturnError(errors.New("some error"))
			}, wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)
			segmentRepoMock := New(poolMock, logging.NewForMocks())
			err := segmentRepoMock.Delete(tc.args.ctx, tc.args.segmentName)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
