package pg_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
	"github.com/VladPetriv/scanner_backend/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestRepliePg_GetFullRepliesByMessageID(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewReplieRepo(&pg.DB{DB: db})

	tests := []struct {
		name  string
		mock  func()
		input int
		want  []model.FullReplie
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "imageurl", "id", "fullname", "imageurl"}).
					AddRow(1, "test1", "testr1.jpg", 1, "test1 test", "test1.jpg").
					AddRow(2, "test2", "testr2.jpg", 2, "test2 test", "test2.jpg")

				mock.ExpectQuery("SELECT r.id, r.title, r.imageurl, u.id, u.fullname, u.imageurl FROM replie r LEFT JOIN tg_user u ON r.user_id = u.id WHERE r.message_id = $1 ORDER BY r.id DESC NULLS LAST;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want: []model.FullReplie{
				{ID: 1, Title: "test1", ImageURL: "testr1.jpg", UserID: 1, FullName: "test1 test", UserImageURL: "test1.jpg"},
				{ID: 2, Title: "test2", ImageURL: "testr2.jpg", UserID: 2, FullName: "test2 test", UserImageURL: "test2.jpg"},
			},
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "imageurl", "id", "fullname", "imageurl"})

				mock.ExpectQuery("SELECT r.id, r.title, r.imageurl, u.id, u.fullname, u.imageurl FROM replie r LEFT JOIN tg_user u ON r.user_id = u.id WHERE r.message_id = $1 ORDER BY r.id DESC NULLS LAST;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "replies not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "imageurl", "id", "fullname", "imageurl"})

				mock.ExpectQuery("SELECT r.id, r.title, r.imageurl, u.id, u.fullname, u.imageurl FROM replie r LEFT JOIN tg_user u ON r.user_id = u.id WHERE r.message_id = $1 ORDER BY r.id DESC NULLS LAST;").
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
