package pg_test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
	"github.com/VladPetriv/scanner_backend/pkg/util"
	"github.com/stretchr/testify/assert"
)

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
				rows := sqlmock.NewRows([]string{"id", "name", "title", "imageurl"}).
					AddRow(1, "test1", "test1", "test1.jpg").
					AddRow(2, "test2", "test2", "test2.jpg")

				mock.ExpectQuery("SELECT * FROM channel;").
					WillReturnRows(rows)
			},
			want: []model.Channel{
				{ID: 1, Name: "test1", Title: "test1", ImageURL: "test1.jpg"},
				{ID: 2, Name: "test2", Title: "test2", ImageURL: "test2.jpg"},
			},
		},
		{
			name: "channels not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "imageurl"})
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
				rows := sqlmock.NewRows([]string{"id", "name", "title", "imageurl"}).
					AddRow(1, "test1", "test1", "test1.jpg").
					AddRow(2, "test2", "test2", "test2.jpg")

				mock.ExpectQuery("SELECT * FROM channel LIMIT 10 OFFSET $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			want: []model.Channel{
				{ID: 1, Name: "test1", Title: "test1", ImageURL: "test1.jpg"},
				{ID: 2, Name: "test2", Title: "test2", ImageURL: "test2.jpg"},
			},
			input: 1,
		},
		{
			name: "channels not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "imageurl"})

				mock.ExpectQuery("SELECT * FROM channel LIMIT 10 OFFSET $1;").
					WithArgs(1).WillReturnRows(rows)

			},
			input: 1,
			want:  nil,
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "imageurl"})

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
				rows := sqlmock.NewRows([]string{"id", "name", "title", "imageurl"}).
					AddRow(1, "test", "test", "test.jpg")

				mock.ExpectQuery("SELECT * FROM channel WHERE name=$1;").
					WithArgs("test").WillReturnRows(rows)
			},
			input: "test",
			want:  &model.Channel{ID: 1, Name: "test", Title: "test", ImageURL: "test.jpg"},
		},
		{
			name: "channel not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "imageurl"})

				mock.ExpectQuery("SELECT * FROM channel WHERE name=$1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "tilte", "imageurl"})

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

func TestChannelPg_GetChannelStats(t *testing.T) {
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
		want  *model.Stat
	}{
		{
			name: "Ok: [Stat found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "count"}).
					AddRow(1, 12).
					AddRow(2, 0)

				mock.ExpectQuery("SELECT m.id, COUNT(r.id) FROM channel c LEFT JOIN message m ON m.channel_id = c.id LEFT JOIN replie r ON r.message_id = m.id WHERE c.id = $1 GROUP BY m.id;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want: &model.Stat{
				MessagesCount: 2,
				RepliesCount:  12,
			},
		},
		{
			name: "Error: [Empty field]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "count"})

				mock.ExpectQuery("SELECT m.id, COUNT(r.id) FROM channel c LEFT JOIN message m ON m.channel_id = c.id LEFT JOIN replie r ON r.message_id = m.id WHERE c.id = $1 GROUP BY m.id;").
					WithArgs().WillReturnRows(rows)
			},
			input: 1,
			want: &model.Stat{
				MessagesCount: 0,
				RepliesCount:  0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannelStats(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
