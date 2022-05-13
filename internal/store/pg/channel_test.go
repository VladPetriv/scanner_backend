package pg_test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
	"github.com/VladPetriv/scanner_backend/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestChannelPg_GetChannel(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewChannelRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    *model.Channel
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "photourl"}).
					AddRow(1, "test", "test", "test.jpg")

				mock.ExpectQuery("SELECT * FROM channel WHERE id=$1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  &model.Channel{ID: 1, Name: "test", Title: "test", PhotoURL: "test.jpg"},
		},
		{
			name: "channel not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "photourl"})

				mock.ExpectQuery("SELECT * FROM channel WHERE id=$1;").
					WithArgs(404).WillReturnRows(rows)
			},
			input:   404,
			wantErr: true,
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "photourl"})

				mock.ExpectQuery("SELECT * FROM channel WHERE id=$1;").
					WithArgs().WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannel(tt.input)
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

func TestChannelPg_GetChannels(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewChannelRepo(&pg.DB{DB: db})

	tests := []struct {
		name string
		mock func()
		want []model.Channel
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "photourl"}).
					AddRow(1, "test1", "test1", "test1.jpg").
					AddRow(2, "test2", "test2", "test2.jpg")

				mock.ExpectQuery("SELECT * FROM channel;").
					WillReturnRows(rows)
			},
			want: []model.Channel{
				{ID: 1, Name: "test1", Title: "test1", PhotoURL: "test1.jpg"},
				{ID: 2, Name: "test2", Title: "test2", PhotoURL: "test2.jpg"},
			},
		},
		{
			name: "channels not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "photourl"})
				mock.ExpectQuery("SELECT * FROM channel;").
					WillReturnRows(rows)
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannels()

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestChannelPg_GetChannelsByPage(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewChannelRepo(&pg.DB{DB: db})

	tests := []struct {
		name  string
		mock  func()
		input int
		want  []model.Channel
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "photourl"}).
					AddRow(1, "test1", "test1", "test1.jpg").
					AddRow(2, "test2", "test2", "test2.jpg")

				mock.ExpectQuery("SELECT * FROM channel LIMIT 10 OFFSET $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			want: []model.Channel{
				{ID: 1, Name: "test1", Title: "test1", PhotoURL: "test1.jpg"},
				{ID: 2, Name: "test2", Title: "test2", PhotoURL: "test2.jpg"},
			},
			input: 1,
		},
		{
			name: "channels not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "photourl"})

				mock.ExpectQuery("SELECT * FROM channel LIMIT 10 OFFSET $1;").
					WithArgs(1).WillReturnRows(rows)

			},
			input: 1,
			want:  nil,
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "photourl"})

				mock.ExpectQuery("SELECT * FROM channel LIMIT 10 OFFSET $1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannelsByPage(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
func TestChannelPg_GetByName(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewChannelRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    *model.Channel
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "photourl"}).
					AddRow(1, "test", "test", "test.jpg")

				mock.ExpectQuery("SELECT * FROM channel WHERE name=$1;").
					WithArgs("test").WillReturnRows(rows)
			},
			input: "test",
			want:  &model.Channel{ID: 1, Name: "test", Title: "test", PhotoURL: "test.jpg"},
		},
		{
			name: "channel not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "photourl"})

				mock.ExpectQuery("SELECT * FROM channel WHERE name=$1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "tilte", "photourl"})

				mock.ExpectQuery("SELECT * FROM channel WHERE name=$1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannelByName(tt.input)
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
