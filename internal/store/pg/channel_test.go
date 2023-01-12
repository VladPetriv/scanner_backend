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

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewChannelRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         *model.DBChannel
		expectedError error
	}{
		{
			name: "CreateChannel successful",
			mock: func() {
				mock.ExpectExec(`
					INSERT INTO channel(name, title, image_url) VALUES ($1, $2, $3);`,
				).WithArgs("test", "test T", "test.jpg").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: &model.DBChannel{Name: "test", Title: "test T", ImageURL: "test.jpg"},
		},
		{
			name: "CreateChannel failed with some sql error",
			mock: func() {
				mock.ExpectExec(`
						INSERT INTO channel(name, title, image_url) VALUES ($1, $2, $3);`,
				).WithArgs("test", "test T", "test.jpg").WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         &model.DBChannel{Name: "test", Title: "test T", ImageURL: "test.jpg"},
			expectedError: fmt.Errorf("create channel: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err = r.CreateChannel(tt.input)
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
		defer db.Close()
	})
}

func Test_GetChannels(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewChannelRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		want          []model.Channel
		expectedError error
	}{
		{
			name: "GetChannels successful",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "image_url"}).
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
			name: "GetChannels failed with not found channels",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "image_url"})

				mock.ExpectQuery("SELECT * FROM channel;").
					WillReturnRows(rows)
			},
			expectedError: nil,
		},
		{
			name: "GetChannels failed with some sql error",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM channel;").
					WillReturnError(fmt.Errorf("some sql error"))
			},
			expectedError: fmt.Errorf("get channels: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannels()
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
		defer db.Close()
	})
}

func Test_GetChannelsByPage(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewChannelRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         int
		want          []model.Channel
		expectedError error
	}{
		{
			name: "GetChannelsByPage successful",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "image_url"}).
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
			name: "GetChannelsByPage failed with not found channels",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "image_url"})

				mock.ExpectQuery("SELECT * FROM channel LIMIT 10 OFFSET $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input:         1,
			expectedError: nil,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM channel LIMIT 10 OFFSET $1;").
					WithArgs(1).WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         1,
			expectedError: fmt.Errorf("get channels by page: %w", fmt.Errorf("some sql error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannelsByPage(tt.input)
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

func Test_GetChannelByName(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewChannelRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         string
		want          *model.Channel
		expectedError error
	}{
		{
			name: "GetChannelByName successful",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "image_url"}).
					AddRow(1, "test", "test", "test.jpg")

				mock.ExpectQuery("SELECT * FROM channel WHERE name = $1;").
					WithArgs("test").WillReturnRows(rows)
			},
			input: "test",
			want:  &model.Channel{ID: 1, Name: "test", Title: "test", ImageURL: "test.jpg"},
		},
		{
			name: "GetChannelByName failed with not found channel",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "image_url"})

				mock.ExpectQuery("SELECT * FROM channel WHERE name = $1;").
					WithArgs("test").WillReturnRows(rows)
			},
			input:         "test",
			expectedError: nil,
		},
		{
			name: "GetChannelByName failed with some sql error",
			mock: func() {
				mock.ExpectQuery("SELECT * FROM channel WHERE name = $1;").
					WithArgs("test").WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         "test",
			expectedError: fmt.Errorf("get channel by name: %w", fmt.Errorf("some sql error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannelByName(tt.input)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
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

func Test_GetChannelStats(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewChannelRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         int
		want          *model.Stat
		expectedError error
	}{
		{
			name: "GetChannelStats successful",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "count"}).
					AddRow(1, 12).
					AddRow(2, 0)

				mock.ExpectQuery(
					`SELECT m.id, COUNT(r.id) 
					FROM channel c LEFT JOIN message m ON m.channel_id = c.id 
					LEFT JOIN reply r ON r.message_id = m.id 
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
			name: "GetChannelStats failed with not found statistics",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "count"})

				mock.ExpectQuery(`SELECT m.id, COUNT(r.id) 
					FROM channel c LEFT JOIN message m ON m.channel_id = c.id 
					LEFT JOIN reply r ON r.message_id = m.id 
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
			name: "GetChannelStats failed with some sql error",
			mock: func() {
				mock.ExpectQuery(`SELECT m.id, COUNT(r.id) 
					FROM channel c LEFT JOIN message m ON m.channel_id = c.id 
					LEFT JOIN reply r ON r.message_id = m.id 
					WHERE c.id = $1 GROUP BY m.id;`,
				).WithArgs(1).WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         1,
			expectedError: fmt.Errorf("get channel statistic: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetChannelStats(tt.input)
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
}
