package interview

import (
	"database/sql"
	"event-migration-script/helpers"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"testing"
)

func Test_interview_GetEventIDMap(t *testing.T) {
	ctx, mock := helpers.GetContextWithMockDB()

	tests := []struct {
		name      string
		want      map[int]string
		wantErr   bool
		execMocks func()
	}{
		{
			name: "success",
			want: map[int]string{1: "1000001", 2: "1000002"},
			execMocks: func() {
				mock.ExpectQuery(getEventIDMapQuery).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "event_id"}).
							AddRow(1, "1000001").
							AddRow(2, "1000002"))
			},
		},
		{
			name:    "failure - scan error",
			wantErr: true,
			execMocks: func() {
				mock.ExpectQuery(getEventIDMapQuery).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(1))
			},
		},
		{
			name:    "failure - query error",
			wantErr: true,
			execMocks: func() {
				mock.ExpectQuery(getEventIDMapQuery).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			i := New()

			tt.execMocks()

			got, err := i.GetEventIDMap(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventIDMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEventIDMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_interview_UpdateEventID(t *testing.T) {
	ctx, mock := helpers.GetContextWithMockDB()

	type args struct {
		oldEventID string
		newEventID string
	}

	tests := []struct {
		name      string
		args      args
		wantErr   bool
		execMocks func()
	}{
		{
			name: "success",
			args: args{oldEventID: "100001", newEventID: "100002"},
			execMocks: func() {
				mock.ExpectExec(updateEventIDQuery).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name:    "failure - sql error",
			args:    args{oldEventID: "100001", newEventID: "100002"},
			wantErr: true,
			execMocks: func() {
				mock.ExpectExec(updateEventIDQuery).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			i := New()

			tt.execMocks()

			if err := i.UpdateEventID(ctx, tt.args.oldEventID, tt.args.newEventID); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEventID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
