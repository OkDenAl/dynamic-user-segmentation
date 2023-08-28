package db

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/pkg/logging"
	"errors"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository_GetAllUserSegments(t *testing.T) {
	type args struct {
		ctx    context.Context
		userId int64
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantRes      []entity.Segment
		wantErr      bool
	}{
		{
			name: "OK", args: args{ctx: context.Background(), userId: 1},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"segment_name"}).
					AddRow("test1").AddRow("test2")

				m.ExpectQuery("SELECT segment_name").
					WithArgs(args.userId).
					WillReturnRows(rows)
			}, wantRes: []entity.Segment{{Name: "test1"}, {Name: "test2"}}, wantErr: false,
		},
		{
			name: "Unexpected error", args: args{ctx: context.Background(), userId: 1},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {

				m.ExpectQuery("SELECT segment_name").
					WithArgs(args.userId).
					WillReturnError(errors.New("some error"))
			}, wantRes: nil, wantErr: true,
		},
		{
			name: "Invalid data type from bd", args: args{ctx: context.Background(), userId: 1},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"segment_name"}).
					AddRow(1)

				m.ExpectQuery("SELECT segment_name").
					WithArgs(args.userId).
					WillReturnRows(rows)
			}, wantRes: nil, wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)
			segmentRepoMock := NewUserSegment(poolMock, logging.NewForMocks())
			got, err := segmentRepoMock.GetAllUserSegments(tc.args.ctx, tc.args.userId)
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
