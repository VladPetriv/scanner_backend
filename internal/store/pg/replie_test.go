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
			want: nil,
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id", "title"})

				mock.ExpectQuery("SELECT * FROM replie WHERE id=$1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetReplie(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)

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
					WithArgs("lost").WillReturnRows(rows)
			},
			input: "lost",
			want:  nil,
		},
		{
			name: "empty field",
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

func TestRepliePg_GetFullRepliesByMessageID(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewReplieRepo(&DB{DB: db})

	tests := []struct {
		name  string
		mock  func()
		input int
		want  []model.FullReplie
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "id", "fullname", "photourl"}).
					AddRow(1, "test1", 1, "test1 test", "test1.jpg").
					AddRow(2, "test2", 2, "test2 test", "test2.jpg")

				mock.ExpectQuery("SELECT r.id, r.title, u.id, u.fullname, u.photourl FROM replie r LEFT JOIN tg_user u ON r.user_id = u.id WHERE r.message_id = $1 ORDER BY r.id DESC NULLS LAST;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want: []model.FullReplie{
				{ID: 1, Title: "test1", UserID: 1, FullName: "test1 test", PhotoURL: "test1.jpg"},
				{ID: 2, Title: "test2", UserID: 2, FullName: "test2 test", PhotoURL: "test2.jpg"},
			},
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "id", "fullname", "photourl"})

				mock.ExpectQuery("SELECT r.id, r.title, u.id, u.fullname, u.photourl FROM replie r LEFT JOIN tg_user u ON r.user_id = u.id WHERE r.message_id = $1 ORDER BY r.id DESC NULLS LAST;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "replies not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "id", "fullname", "photourl"})

				mock.ExpectQuery("SELECT r.id, r.title, u.id, u.fullname, u.photourl FROM replie r LEFT JOIN tg_user u ON r.user_id = u.id WHERE r.message_id = $1 ORDER BY r.id DESC NULLS LAST;").
					WithArgs(404).WillReturnRows(rows)
			},
			input: 404,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullRepliesByMessageID(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
