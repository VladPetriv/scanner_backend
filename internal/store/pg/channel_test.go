package pg_test

import (
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
	"github.com/VladPetriv/scanner_backend/pkg/util"
)

func Test_CreateChannel(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewChannelRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		input   *model.DBChannel
		wantErr bool
	}{
		{
			name: "Ok: [channel created]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)

				mock.ExpectQuery(`
					INSERT INTO channel(name, title, imageurl) 
					VALUES ($1, $2, $3) RETURNING id;`,
				).WithArgs("test", "test T", "test.jpg").WillReturnRows(rows)
			},
			input: &model.DBChannel{Name: "test", Title: "test T", ImageURL: "test.jpg"},
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery(`
					INSERT INTO channel(name, title, imageurl) 
					VALUES ($1, $2, $3) RETURNING id;`,
				).WithArgs("test", "test T", "test.jpg").WillReturnError(fmt.Errorf("some sql error"))
			},
			input:   &model.DBChannel{Name: "test", Title: "test T", ImageURL: "test.jpg"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.CreateChannel(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_GetChannels(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewChannelRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		want    []model.Channel
		wantErr bool
	}{
		{
			name: "Ok: [channels found]",
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
			name: "Error: channels not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "imageurl"})

				mock.ExpectQuery("SELECT * FROM channel;").
					WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM channel;").
					WillReturnError(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannels()

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

func Test_GetChannelsByPage(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewChannelRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    []model.Channel
		wantErr bool
	}{
		{
			name: "Ok: [channels found]",
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
			name: "Error: [channels not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "imageurl"})

				mock.ExpectQuery("SELECT * FROM channel LIMIT 10 OFFSET $1;").
					WithArgs(1).WillReturnRows(rows)

			},
			input:   1,
			wantErr: true,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM channel LIMIT 10 OFFSET $1;").
					WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			input:   1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannelsByPage(tt.input)
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

func Test_GetChannelByName(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewChannelRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    *model.Channel
		wantErr bool
	}{
		{
			name: "Ok: [channel found]",
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
			name: "Error: [channel not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "imageurl"})

				mock.ExpectQuery("SELECT * FROM channel WHERE name=$1;").
					WithArgs("test").WillReturnRows(rows)
			},
			input:   "test",
			wantErr: true,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM channel WHERE name=$1;").
					WithArgs("test").WillReturnError(fmt.Errorf("some error"))
			},
			input:   "test",
			wantErr: true,
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

func Test_GetChannelStats(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewChannelRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    *model.Stat
		wantErr bool
	}{
		{
			name: "Ok: [channel stat found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "count"}).
					AddRow(1, 12).
					AddRow(2, 0)

				mock.ExpectQuery(`SELECT m.id, COUNT(r.id) 
					FROM channel c LEFT JOIN message m ON m.channel_id = c.id 
					LEFT JOIN replie r ON r.message_id = m.id 
					WHERE c.id = $1 GROUP BY m.id;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want: &model.Stat{
				MessagesCount: 2,
				RepliesCount:  12,
			},
		},
		{
			name: "Error: [channel stat not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "count"})

				mock.ExpectQuery(`SELECT m.id, COUNT(r.id) 
					FROM channel c LEFT JOIN message m ON m.channel_id = c.id 
					LEFT JOIN replie r ON r.message_id = m.id 
					WHERE c.id = $1 GROUP BY m.id;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want: &model.Stat{
				MessagesCount: 0,
				RepliesCount:  0,
			},
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery(`SELECT m.id, COUNT(r.id) 
					FROM channel c LEFT JOIN message m ON m.channel_id = c.id 
					LEFT JOIN replie r ON r.message_id = m.id 
					WHERE c.id = $1 GROUP BY m.id;`,
				).WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			input:   1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannelStats(tt.input)
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
