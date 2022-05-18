package pg_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
	"github.com/VladPetriv/scanner_backend/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestSavedPg_GetSavedMessgaes(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewSavedRepo(&pg.DB{DB: db})

	tests := []struct {
		name  string
		mock  func()
		input int
		want  []model.Saved
	}{
		{
			name: "OK: [Messages found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id"}).
					AddRow(1, 2, 1).
					AddRow(2, 2, 3)

				mock.ExpectQuery("SELECT * FROM saved WHERE user_id=$1;").
					WithArgs(2).WillReturnRows(rows)
			},
			input: 2,
			want: []model.Saved{
				{ID: 1, WebUserID: 2, MessageID: 1},
				{ID: 2, WebUserID: 2, MessageID: 3},
			},
		},
		{
			name: "Error: [Messages not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id"})

				mock.ExpectQuery("SELECT * FROM saved WHERE user_id=$1;").
					WithArgs(2).WillReturnRows(rows)
			},
			input: 2,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetSavedMessages(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSavedPg_GetSavedMessageByMessageID(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewSavedRepo(&pg.DB{DB: db})

	tests := []struct {
		name  string
		mock  func()
		input int
		want  *model.Saved
	}{
		{
			name: "OK: [Message found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id"}).
					AddRow(1, 2, 1)

				mock.ExpectQuery("SELECT * FROM saved WHERE message_id=$1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  &model.Saved{ID: 1, WebUserID: 2, MessageID: 1},
		},
		{
			name: "Error: [Messages not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id"})

				mock.ExpectQuery("SELECT * FROM saved WHERE message_id=$1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetSavedMessageByMessageID(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSavedPg_CreateSavedMessage(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewSavedRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   *model.Saved
		want    int
		wantErr bool
	}{
		{
			name: "OK: [Messages created]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectQuery("INSERT INTO saved(user_id, message_id) VALUES ($1, $2) RETURNING id;").
					WithArgs(1, 2).WillReturnRows(rows)
			},
			input: &model.Saved{WebUserID: 1, MessageID: 2},
			want:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateSavedMessage(tt.input)

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
