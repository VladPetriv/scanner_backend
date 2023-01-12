package pg_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
	"github.com/VladPetriv/scanner_backend/pkg/util"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_CreateWebUser(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewWebUserRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         *model.WebUser
		expectedError error
	}{
		{
			name: "CreateWebUser successful",
			mock: func() {
				mock.ExpectExec("INSERT INTO web_user(email, password) VALUES ($1, $2);").
					WithArgs("test@test.com", "test").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: &model.WebUser{Email: "test@test.com", Password: "test"},
		},
		{
			name: "CreateWebUser failed with some sql error",
			mock: func() {
				mock.ExpectExec("INSERT INTO web_user(email, password) VALUES ($1, $2);").
					WithArgs("test@test.com", "test").WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         &model.WebUser{Email: "test@test.com", Password: "test"},
			expectedError: fmt.Errorf("create web user: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.CreateWebUser(tt.input)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
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

func Test_GetWebUserByID(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewWebUserRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         int
		want          *model.WebUser
		expectedError error
	}{
		{
			name: "GetWebUserByID successful",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"}).
					AddRow(1, "test@test.com", "test")

				mock.ExpectQuery("SELECT * FROM web_user WHERE id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  &model.WebUser{ID: 1, Email: "test@test.com", Password: "test"},
		},
		{
			name: "GetWebUserByID failed with not found user",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"})

				mock.ExpectQuery("SELECT * FROM web_user WHERE id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input:         1,
			expectedError: nil,
		},
		{
			name: "GetWebUserByID failed with some sql error",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM web_user WHERE id = $1;").
					WithArgs(1).WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         1,
			expectedError: fmt.Errorf("get web user by id: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetWebUserByID(tt.input)
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

func Test_GetWebUserByEmail(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewWebUserRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         string
		want          *model.WebUser
		expectedError error
	}{
		{
			name: "BetWebUserByEmail successful",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"}).
					AddRow(1, "test@test.com", "test")

				mock.ExpectQuery("SELECT * FROM web_user WHERE email = $1;").
					WithArgs("test@test.com").WillReturnRows(rows)
			},
			input: "test@test.com",
			want:  &model.WebUser{ID: 1, Email: "test@test.com", Password: "test"},
		},
		{
			name: "GetWebUserByEmail failed with not found user",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"})

				mock.ExpectQuery("SELECT * FROM web_user WHERE email = $1;").
					WithArgs("test@test.com").WillReturnRows(rows)
			},
			input:         "test@test.com",
			expectedError: nil,
		},
		{
			name: "GetWebUserByEmail failed with some sql error",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM web_user WHERE email = $1;").
					WithArgs("test@test.com").WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         "test@test.com",
			expectedError: fmt.Errorf("get web user by email: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetWebUserByEmail(tt.input)
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
