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

func Test_CreateUser(t *testing.T) {
	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewUserRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         *model.User
		want          int
		expectedError error
	}{
		{
			name: "CreateUser successful",
			mock: func() {
				row := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(`
					INSERT INTO tg_user(username, fullname, image_url) VALUES ($1, $2, $3) RETURNING id;`,
				).WithArgs("test", "test test", "test.jpg").WillReturnRows(row)
			},
			input: &model.User{Username: "test", FullName: "test test", ImageURL: "test.jpg"},
			want:  1,
		},
		{
			name: "CreateUser failed with some sql error",
			mock: func() {
				mock.ExpectQuery(`
					INSERT INTO tg_user(username, fullname, image_url) VALUES ($1, $2, $3) RETURNING id;`,
				).WithArgs("test", "test test", "test.jpg").WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         &model.User{Username: "test", FullName: "test test", ImageURL: "test.jpg"},
			expectedError: fmt.Errorf("some sql error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateUser(tt.input)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Cleanup(func() {
		db.Close()
	})
}

func Test_GetUserByUsername(t *testing.T) {
	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewUserRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         string
		want          *model.User
		expectedError error
	}{
		{
			name: "GetUserByUsername successful",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "image_url"}).
					AddRow(1, "test", "test test", "test.jpg")

				mock.ExpectQuery("SELECT * FROM tg_user WHERE username = $1;").
					WithArgs("test").WillReturnRows(rows)
			},
			input: "test",
			want:  &model.User{ID: 1, Username: "test", FullName: "test test", ImageURL: "test.jpg"},
		},
		{
			name: "GetUserByUserName failed with not found user",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "image_url"})

				mock.ExpectQuery("SELECT * FROM tg_user WHERE username = $1;").
					WithArgs().WillReturnRows(rows)
			},
			input:         "test",
			expectedError: nil,
		},
		{
			name: "GetUserByUsername failed with some sql error",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM tg_user WHERE username = $1;").
					WithArgs("test").WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         "test",
			expectedError: fmt.Errorf("some sql error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUserByUsername(tt.input)
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

func Test_GetUserByID(t *testing.T) {
	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewUserRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         int
		want          *model.User
		expectedError error
	}{
		{
			name: "GetUserByID successfully",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "image_url"}).
					AddRow(1, "test", "test test", "test.jpg")

				mock.ExpectQuery("SELECT * FROM tg_user WHERE id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  &model.User{ID: 1, Username: "test", FullName: "test test", ImageURL: "test.jpg"},
		},
		{
			name: "GetUserByID failed with not found user",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "image_url"})

				mock.ExpectQuery("SELECT * FROM tg_user WHERE id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input:         1,
			expectedError: nil,
		},
		{
			name: "GetUserByID failed with some sql error",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM tg_user WHERE id = $1;").
					WithArgs(1).WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         1,
			expectedError: fmt.Errorf("some sql error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUserByID(tt.input)
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
