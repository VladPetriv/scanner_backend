package pg_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
	"github.com/VladPetriv/scanner_backend/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestWebUserPg_CreateUser(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewWebUserRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   *model.WebUser
		want    int
		wantErr bool
	}{
		{
			name: "Ok: [User created]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectQuery("INSERT INTO web_user(email, password) VALUES ($1, $2) RETURNING id;").
					WithArgs("test@test.com", "test").WillReturnRows(rows)

			},
			input: &model.WebUser{Email: "test@test.com", Password: "test"},
			want:  1,
		},
		{
			name: "Error: [Empty fields]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("INSERT INTO web_user(email, password) VALUES ($1, $2) RETURNING id;").
					WithArgs("", "").WillReturnRows(rows)
			},
			input:   &model.WebUser{},
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

func TestWebUserPg_GetUser(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewWebUserRepo(&pg.DB{DB: db})

	tests := []struct {
		name  string
		mock  func()
		input int
		want  *model.WebUser
	}{
		{
			name: "Ok: [User found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"}).
					AddRow(1, "test@test.com", "test")

				mock.ExpectQuery("SELECT * FROM web_user WHERE id=$1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  &model.WebUser{ID: 1, Email: "test@test.com", Password: "test"},
		},
		{
			name: "Error: [User not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"})

				mock.ExpectQuery("SELECT * FROM web_user WHERE id=$1;").
					WithArgs(404).WillReturnRows(rows)
			},
			input: 404,
			want:  nil,
		},
		{
			name: "Error: [Empty field]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"})

				mock.ExpectQuery("SELECT * FROM web_user WHERE id=$1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetWebUser(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)

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

	r := pg.NewWebUserRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    *model.WebUser
		wantErr bool
	}{
		{
			name: "Ok: [User found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"}).
					AddRow(1, "test@test.com", "test")

				mock.ExpectQuery("SELECT * FROM web_user WHERE email=$1;").
					WithArgs("test@test.com").WillReturnRows(rows)
			},
			input: "test@test.com",
			want:  &model.WebUser{ID: 1, Email: "test@test.com", Password: "test"},
		},
		{
			name: "Error: [User not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"})

				mock.ExpectQuery("SELECT * FROM web_user WHERE email=$1;").
					WithArgs("test@test.com").WillReturnRows(rows)
			},
			input: "test@test.com",
			want:  nil,
		},
		{
			name: "Error: [Empty field]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password"})

				mock.ExpectQuery("SELECT * FROM web_user WHERE email=$1;").
					WithArgs("").WillReturnRows(rows)
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetWebUserByEmail(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
