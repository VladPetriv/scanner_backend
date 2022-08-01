package pg_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
	"github.com/VladPetriv/scanner_backend/pkg/util"
)

func TestRepliePg_CreateReplie(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewReplieRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		input   *model.DBReplie
		wantErr bool
	}{
		{
			name: "Ok: [replie created]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)

				mock.ExpectQuery(`
					INSERT INTO replie(user_id, message_id, title, imageurl) 
					VALUES ($1, $2, $3, $4) RETURNING id;`,
				).WithArgs(1, 1, "test", "test.jpg").WillReturnRows(rows)
			},
			input: &model.DBReplie{UserID: 1, MessageID: 1, Title: "test", ImageURL: "test.jpg"},
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery(`
					INSERT INTO replie(user_id, message_id, title, imageurl) 
					VALUES ($1, $2, $3, $4) RETURNING id;`,
				).WithArgs(1, 1, "test", "test.jpg").WillReturnError(fmt.Errorf("some sql error"))
			},
			input:   &model.DBReplie{UserID: 1, MessageID: 1, Title: "test", ImageURL: "test.jpg"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.CreateReplie(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
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

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewReplieRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    []model.FullReplie
		wantErr bool
	}{
		{
			name: "Ok: [replies found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "imageurl", "user_id", "fullname", "userimageurl"}).
					AddRow(1, "test1", "testr1.jpg", 1, "test1 test", "test1.jpg").
					AddRow(2, "test2", "testr2.jpg", 2, "test2 test", "test2.jpg")

				mock.ExpectQuery(
					`SELECT r.id, r.title, r.imageurl, 
					u.id as user_id, u.fullname, u.imageurl as userimageurl
					FROM replie r 
					LEFT JOIN tg_user u ON r.user_id = u.id 
					WHERE r.message_id = $1 
					ORDER BY r.id DESC NULLS LAST;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want: []model.FullReplie{
				{ID: 1, Title: "test1", ImageURL: "testr1.jpg", UserID: 1, FullName: "test1 test", UserImageURL: "test1.jpg"},
				{ID: 2, Title: "test2", ImageURL: "testr2.jpg", UserID: 2, FullName: "test2 test", UserImageURL: "test2.jpg"},
			},
		},
		{
			name: "Error: [replies not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "imageurl", "user_id", "fullname", "userimageurl"})

				mock.ExpectQuery(
					`SELECT r.id, r.title, r.imageurl, 
					u.id as user_id, u.fullname, u.imageurl as userimageurl
					FROM replie r 
					LEFT JOIN tg_user u ON r.user_id = u.id 
					WHERE r.message_id = $1 
					ORDER BY r.id DESC NULLS LAST;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input:   1,
			wantErr: true,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery(
					`SELECT r.id, r.title, r.imageurl, 
					u.id as user_id, u.fullname, u.imageurl as userimageurl
					FROM replie r 
					LEFT JOIN tg_user u ON r.user_id = u.id 
					WHERE r.message_id = $1 
					ORDER BY r.id DESC NULLS LAST;`,
				).WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			input:   1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullRepliesByMessageID(tt.input)
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
