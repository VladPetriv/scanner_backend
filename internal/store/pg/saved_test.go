package pg_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store/mocks"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
)

func Test_CreateSavedMessage(t *testing.T) {
	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewSavedRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         *model.Saved
		expectedError error
	}{
		{
			name: "CreateSavedMessage successful",
			mock: func() {
				mock.ExpectExec("INSERT INTO saved(user_id, message_id) VALUES ($1, $2);").
					WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: &model.Saved{WebUserID: 1, MessageID: 2},
		},
		{
			name: "CreateSavedMessage failed with some sql error",
			mock: func() {
				mock.ExpectExec("INSERT INTO saved(user_id, message_id) VALUES ($1, $2);").
					WithArgs(1, 2).WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         &model.Saved{WebUserID: 1, MessageID: 2},
			expectedError: fmt.Errorf("create saved message: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.CreateSavedMessage(tt.input)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
	t.Cleanup(func() {
		db.Close()
	})
}

func Test_GetSavedMessages(t *testing.T) {
	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewSavedRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         int
		want          []model.Saved
		expectedError error
	}{
		{
			name: "GetSavedMessages successful",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id"}).
					AddRow(1, 2, 1).
					AddRow(2, 2, 3)

				mock.ExpectQuery("SELECT * FROM saved WHERE user_id = $1;").
					WithArgs(2).WillReturnRows(rows)
			},
			input: 2,
			want: []model.Saved{
				{ID: 1, WebUserID: 2, MessageID: 1},
				{ID: 2, WebUserID: 2, MessageID: 3},
			},
		},
		{
			name: "GetSavedMessages failed with not found messages",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id"})

				mock.ExpectQuery("SELECT * FROM saved WHERE user_id = $1;").
					WithArgs(2).WillReturnRows(rows)
			},
			input:         2,
			expectedError: nil,
		},
		{
			name: "GetSavedMessages failed with some sql error",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM saved WHERE user_id = $1;").
					WithArgs(2).WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         2,
			expectedError: fmt.Errorf("get saved messages: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetSavedMessages(tt.input)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Cleanup(func() {
		db.Close()
	})
}

func Test_GetSavedMessageByID(t *testing.T) {
	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewSavedRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         int
		want          *model.Saved
		expectedError error
	}{
		{
			name: "GetSavedMessageByID successful",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id"}).
					AddRow(1, 2, 1)

				mock.ExpectQuery("SELECT * FROM saved WHERE message_id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  &model.Saved{ID: 1, WebUserID: 2, MessageID: 1},
		},
		{
			name: "GetSavedMessageByID failed with not found message",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "message_id"})

				mock.ExpectQuery("SELECT * FROM saved WHERE message_id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input:         1,
			expectedError: nil,
		},
		{
			name: "GetSavedMessageByID failed with some sql error",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM saved WHERE message_id = $1;").
					WithArgs(1).WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         1,
			expectedError: fmt.Errorf("get saved message by id: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetSavedMessageByID(tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Cleanup(func() {
		db.Close()
	})
}

func Test_DeleteSavedMessage(t *testing.T) {
	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewSavedRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         int
		expectedError error
	}{
		{
			name: "DeleteSavedMessage successful",
			mock: func() {
				mock.ExpectExec("DELETE FROM saved WHERE id = $1;").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: 1,
		},
		{
			name: "DelateSavedMessage failed with some sql error",
			mock: func() {
				mock.ExpectExec("DELETE FROM saved WHERE id = $1;").
					WithArgs(1).WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         1,
			expectedError: fmt.Errorf("delete saved message: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.DeleteSavedMessage(tt.input)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Cleanup(func() {
		db.Close()
	})
}
