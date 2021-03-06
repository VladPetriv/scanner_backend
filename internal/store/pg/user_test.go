package pg_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
	"github.com/VladPetriv/scanner_backend/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestUserPg_GetUsers(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewUserRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		want    []model.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "imageurl"}).
					AddRow(1, "test1", "test test1", "test1.jpg").
					AddRow(2, "test2", "test test2", "test2.jpg")

				mock.ExpectQuery("SELECT * FROM tg_user;").
					WillReturnRows(rows)
			},
			want: []model.User{
				{ID: 1, Username: "test1", FullName: "test test1", ImageURL: "test1.jpg"},
				{ID: 2, Username: "test2", FullName: "test test2", ImageURL: "test2.jpg"},
			},
		},
		{
			name: "users not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "imageurl"})

				mock.ExpectQuery("SELECT * FROM tg_user;").
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUsers()

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

func TestUserPg_GetUserByUsername(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewUserRepo(&pg.DB{DB: db})

	tests := []struct {
		name  string
		mock  func()
		input string
		want  *model.User
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "imageurl"}).
					AddRow(1, "test", "test test", "test.jpg")

				mock.ExpectQuery("SELECT * FROM tg_user WHERE username=$1;").
					WithArgs("test").WillReturnRows(rows)
			},
			input: "test",
			want: &model.User{
				ID:       1,
				Username: "test",
				FullName: "test test",
				ImageURL: "test.jpg",
			},
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "imageurl"})

				mock.ExpectQuery("SELECT * FROM tg_user WHERE username=$1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "user not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "imageurl"})

				mock.ExpectQuery("SELECT * FROM tg_user WHERE username=$1;").
					WithArgs("lost").WillReturnRows(rows)
			},
			input: "lost",
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUserByUsername(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserPg_GetUserByID(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewUserRepo(&pg.DB{DB: db})

	tests := []struct {
		name  string
		mock  func()
		input int
		want  *model.User
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "imageurl"}).
					AddRow(1, "test", "test test", "test.jpg")
				mock.ExpectQuery("SELECT * FROM tg_user WHERE id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  &model.User{ID: 1, Username: "test", FullName: "test test", ImageURL: "test.jpg"},
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "imageurl"})

				mock.ExpectQuery("SELECT * FROM tg_user WHERE id = $1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "user not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "fullname", "imageurl"})

				mock.ExpectQuery("SELECT * FROM tg_user WHERE id = $1;").
					WithArgs(404).WillReturnRows(rows)
			},
			input: 404,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetUserByID(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
