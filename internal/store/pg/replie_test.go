package pg

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestRepliePg_GetReplie(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewReplieRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    *model.Replie
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id", "title"}).
					AddRow(1, 1, 1, "test")

				mock.ExpectQuery("SELECT * FROM replie WHERE id=$1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  &model.Replie{ID: 1, UserID: 1, MessageID: 1, Title: "test"},
		},
		{
			name: "replie not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id", "title"})

				mock.ExpectQuery("SELECT * FROM replie WHERE id=$1;").
					WithArgs().WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetReplie(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepliePg_GetReplies(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewReplieRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		want    []model.Replie
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id", "title"}).
					AddRow(1, 1, 1, "test1").
					AddRow(2, 2, 2, "test2")

				mock.ExpectQuery("SELECT * FROM replie;").
					WillReturnRows(rows)
			},
			want: []model.Replie{
				{ID: 1, UserID: 1, MessageID: 1, Title: "test1"},
				{ID: 2, UserID: 2, MessageID: 2, Title: "test2"},
			},
		},
		{
			name: "replies not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "message_id", "title"})

				mock.ExpectQuery("SELECT * FROM replie;").
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetReplies()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepliePg_GetReplieByName(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewReplieRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    *model.Replie
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id", "title"}).
					AddRow(1, 1, 1, "test")

				mock.ExpectQuery("SELECT * FROM replie WHERE title=$1;").
					WithArgs("test").WillReturnRows(rows)
			},
			input: "test",
			want:  &model.Replie{ID: 1, UserID: 1, MessageID: 1, Title: "test"},
		},
		{
			name: "replie not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id", "title"})

				mock.ExpectQuery("SELECT * FROM replie WHERE title=$1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetReplieByName(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
