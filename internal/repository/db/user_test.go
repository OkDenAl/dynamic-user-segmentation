package db

import (
	"context"
	"dynamic-user-segmentation/pkg/logging"
	"errors"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestRepository_GetNPercentOfUsersIDs(t *testing.T) {
	type args struct {
		ctx          context.Context
		percent      float64
		countOfUsers int
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantRes      []int64
		wantErr      bool
	}{
		{
			name: "OK", args: args{ctx: context.Background(), percent: 50, countOfUsers: 10},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				count := pgxmock.NewRows([]string{""}).AddRow(10)
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow(int64(1)).AddRow(int64(2)).AddRow(int64(3)).AddRow(int64(4)).AddRow(int64(5))

				m.ExpectQuery("SELECT count").WillReturnRows(count)

				m.ExpectQuery("SELECT id").
					WithArgs(math.Round(args.percent/100*float64(args.countOfUsers)), args.countOfUsers).
					WillReturnRows(rows)
			}, wantRes: []int64{1, 2, 3, 4, 5}, wantErr: false,
		},
		{
			name: "Invalid data types from BD", args: args{ctx: context.Background(), percent: 50, countOfUsers: 10},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				count := pgxmock.NewRows([]string{""}).AddRow(10)
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow(1).AddRow(int64(2)).AddRow(int64(3)).AddRow(int64(4)).AddRow(int64(5))

				m.ExpectQuery("SELECT count").WillReturnRows(count)

				m.ExpectQuery("SELECT id").
					WithArgs(math.Round(args.percent/100*float64(args.countOfUsers)), args.countOfUsers).
					WillReturnRows(rows)
			}, wantRes: nil, wantErr: true,
		},
		{
			name: "Unexpected error", args: args{ctx: context.Background(), percent: 50, countOfUsers: 10},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				count := pgxmock.NewRows([]string{""}).AddRow(10)
				m.ExpectQuery("SELECT count").WillReturnRows(count)

				m.ExpectQuery("SELECT id").
					WithArgs(math.Round(args.percent/100*float64(args.countOfUsers)), args.countOfUsers).
					WillReturnError(errors.New("some error"))
			}, wantRes: nil, wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)
			userRepoMock := NewUserRepo(poolMock, logging.NewForMocks())
			got, err := userRepoMock.GetNPercentOfUsersIDs(tc.args.ctx, tc.args.percent)
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
