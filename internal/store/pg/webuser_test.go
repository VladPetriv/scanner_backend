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

func TestWebUserPg_CreateUser(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewWebUserRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		input   *model.WebUser
		want    int
		wantErr bool
	}{
		{
			name: "Ok: [user created]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectQuery("INSERT INTO web_user(email, password) VALUES ($1, $2) RETURNING id;").
					WithArgs("test@test.com", "test").WillReturnRows(rows)

			},
			input: &model.WebUser{Email: "test@test.com", Password: "test"},
			want:  1,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("INSERT INTO web_user(email, password) VALUES ($1, $2) RETURNING id;").
					WithArgs("test@test.com", "test").WillReturnError(fmt.Errorf("some error"))
			},
			input:   &model.WebUser{Email: "test@test.com", Password: "test"},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateWebUser(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestWebUserPg_GetWebUserByID(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewWebUserRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    *model.WebUser
		wantErr bool
	}{
		{
			name: "Ok: [user found]",
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
			name: "Error: [user not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"})

				mock.ExpectQuery("SELECT * FROM web_user WHERE id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input:   1,
			wantErr: true,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM web_user WHERE id = $1;").
					WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			input:   1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetWebUserByID(tt.input)
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

func TestWebUserPg_GetUserByEmail(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewWebUserRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    *model.WebUser
		wantErr bool
	}{
		{
			name: "Ok: [user found]",
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
			name: "Error: [user not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"})

				mock.ExpectQuery("SELECT * FROM web_user WHERE email = $1;").
					WithArgs("test@test.com").WillReturnRows(rows)
			},
			input:   "test@test.com",
			wantErr: true,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM web_user WHERE email = $1;").
					WithArgs("test@test.com").WillReturnError(fmt.Errorf("some error"))
			},
			input:   "test@test.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetWebUserByEmail(tt.input)
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
