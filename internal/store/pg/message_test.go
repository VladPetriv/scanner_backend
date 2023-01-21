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

func Test_CreateMessage(t *testing.T) {
	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         *model.DBMessage
		want          int
		wantErr       bool
		expectedError error
	}{
		{
			name: "CreateMessage successful",
			mock: func() {
				row := sqlmock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectQuery(
					"INSERT INTO message(channel_id, user_id, title, message_url, image_url) VALUES ($1, $2, $3, $4, $5) RETURNING id;", //nolint:lll
				).WithArgs(1, 1, "test", "test.url", "test.jpg").WillReturnRows(row)
			},
			input: &model.DBMessage{ChannelID: 1, UserID: 1, Title: "test", MessageURL: "test.url", ImageURL: "test.jpg"},
			want:  1,
		},
		{
			name: "CreateMessage failed with some sql error",
			mock: func() {
				mock.ExpectQuery(
					"INSERT INTO message(channel_id, user_id, title, message_url, image_url) VALUES ($1, $2, $3, $4, $5) RETURNING id;", //nolint:lll
				).WithArgs(1, 1, "test", "test.url", "test.jpg").WillReturnError(fmt.Errorf("some sql error"))
			},
			input: &model.DBMessage{
				ChannelID:  1,
				UserID:     1,
				Title:      "test",
				MessageURL: "test.url",
				ImageURL:   "test.jpg",
			},
			wantErr:       true,
			expectedError: fmt.Errorf("create message: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateMessage(tt.input)
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

func Test_GetMessagesCount(t *testing.T) {
	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		want          int
		wantErr       bool
		expectedError error
	}{
		{
			name: "GetMessagesCount successful",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).
					AddRow(10)

				mock.ExpectQuery("SELECT COUNT(*) FROM message;").
					WillReturnRows(rows)
			},
			want: 10,
		},
		{
			name: "GetMessagesCount failed with not found messages",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"})

				mock.ExpectQuery("SELECT COUNT(*) FROM message;").
					WillReturnRows(rows)
			},
			expectedError: nil,
		},
		{
			name: "GetMessagesCount failed with some sql error",
			mock: func() {
				mock.ExpectQuery("SELECT COUNT(*) FROM message;").
					WillReturnError(fmt.Errorf("some sql error"))
			},
			expectedError: fmt.Errorf("get messages count: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetMessagesCount()
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

func Test_GetMessagesCountByChannelID(t *testing.T) {
	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         int
		want          int
		expectedError error
	}{
		{
			name: "GetMessagesCountByChannelID",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).
					AddRow(2)

				mock.ExpectQuery(
					"SELECT COUNT(*) FROM message m LEFT JOIN channel c ON c.id = m.channel_id WHERE m.channel_id = $1;",
				).WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  2,
		},
		{
			name: "GetMessagesCountByChannelID failed with not found messages count by channel id",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery(
					"SELECT COUNT(*) FROM message m LEFT JOIN channel c ON c.id = m.channel_id WHERE m.channel_id = $1;",
				).WithArgs(1).WillReturnRows(rows)
			},
			input:         1,
			expectedError: nil,
		},
		{
			name: "GetMessagesCountByChannelID failed with some sql error",
			mock: func() {
				mock.ExpectQuery(
					"SELECT COUNT(*) FROM message m LEFT JOIN channel c ON c.id = m.channel_id WHERE m.channel_id = $1;",
				).WithArgs(1).WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         1,
			expectedError: fmt.Errorf("get messages count by channel id: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		tt.mock()

		got, err := r.GetMessagesCountByChannelID(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.EqualValues(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		assert.NoError(t, mock.ExpectationsWereMet())
	}

	t.Cleanup(func() {
		db.Close()
	})
}
func Test_GetMessageByTitle(t *testing.T) {
	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         string
		want          *model.DBMessage
		expectedError error
	}{
		{
			name: "GetMessageByTitle successful",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "channel_id", "user_id", "title", "message_url", "image_url"}).
					AddRow(1, 1, 1, "test", "test", "test")

				mock.ExpectQuery(
					"SELECT * FROM message WHERE title = $1;",
				).WithArgs("test").WillReturnRows(rows)
			},
			input: "test",
			want:  &model.DBMessage{ID: 1, ChannelID: 1, UserID: 1, Title: "test", MessageURL: "test", ImageURL: "test"},
		},
		{
			name: "GetMessageByTitle failed with not found message",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "channel_id", "user_id", "title", "message_url", "image_url"})

				mock.ExpectQuery(
					"SELECT * FROM message WHERE title = $1;",
				).WithArgs("test").WillReturnRows(rows)
			},
			input:         "test",
			expectedError: nil,
		},
		{
			name: "GetMessageByTitle failed with some sql error",
			mock: func() {
				mock.ExpectQuery(
					"SELECT * FROM message WHERE title = $1;",
				).WithArgs("test").WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         "test",
			expectedError: fmt.Errorf("get message by title: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		tt.mock()

		got, err := r.GetMessageByTitle(tt.input)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.EqualValues(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		assert.NoError(t, mock.ExpectationsWereMet())
	}

	t.Cleanup(func() {
		db.Close()
	})
}

func Test_GetFullMessagesByPage(t *testing.T) {
	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         int
		want          []model.FullMessage
		expectedError error
	}{
		{
			name: "GetFullMessagesByPage successful",
			mock: func() {
				rows := sqlmock.NewRows(
					[]string{
						"id", "title", "message_url",
						"image_url", "channel_id", "channel_name",
						"channel_image_url", "user_id", "fullname",
						"user_image_url", "count",
					},
				).
					AddRow(1, "test1", "test1.com", "test1.jpg", 1, "test1", "test1.jpg", 1, "test1", "test1.jpg", 1).
					AddRow(2, "test2", "test2.com", "test2.jpg", 2, "test2", "test2.jpg", 2, "test2", "test2.jpg", 3)

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.image_url, 
					c.id AS channel_id, c.name AS channel_name, c.image_url AS channel_image_url, 
					u.id AS user_id, u.fullname, u.image_url AS user_image_url, 
					(SELECT COUNT(*) FROM reply WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					ORDER BY m.id DESC NULLS LAST LIMIT 10 OFFSET $1;`,
				).WithArgs(10).WillReturnRows(rows)
			},
			input: 10,
			want: []model.FullMessage{
				{
					ID: 1, Title: "test1", MessageURL: "test1.com", ImageURL: "test1.jpg",
					ChannelID: 1, ChannelName: "test1", ChannelImageURL: "test1.jpg",
					UserID: 1, FullName: "test1", UserImageURL: "test1.jpg",
					RepliesCount: 1,
				},
				{
					ID: 2, Title: "test2", MessageURL: "test2.com", ImageURL: "test2.jpg",
					ChannelID: 2, ChannelName: "test2", ChannelImageURL: "test2.jpg",
					UserID: 2, FullName: "test2", UserImageURL: "test2.jpg",
					RepliesCount: 3,
				},
			},
		},
		{
			name: "GetFullMessagesByPage failed with not found messages",
			mock: func() {
				rows := sqlmock.NewRows(
					[]string{
						"id", "title", "message_url",
						"image_url", "channel_id", "channel_name",
						"channel_image_url", "user_id", "fullname",
						"user_image_url", "count",
					},
				)

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.image_url, 
					c.id AS channel_id, c.name AS channel_name, c.image_url AS channel_image_url, 
					u.id AS user_id, u.fullname, u.image_url AS user_image_url, 
					(SELECT COUNT(*) FROM reply WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					ORDER BY m.id DESC NULLS LAST LIMIT 10 OFFSET $1;`,
				).WithArgs(10).WillReturnRows(rows)
			},
			input:         10,
			expectedError: nil,
		},
		{
			name: "GetFullMessagesByPage failed with some sql error",
			mock: func() {
				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.image_url, 
					c.id AS channel_id, c.name AS channel_name, c.image_url AS channel_image_url, 
					u.id AS user_id, u.fullname, u.image_url AS user_image_url, 
					(SELECT COUNT(*) FROM reply WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					ORDER BY m.id DESC NULLS LAST LIMIT 10 OFFSET $1;`,
				).WithArgs(10).WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         10,
			expectedError: fmt.Errorf("get full messages by page: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessagesByPage(tt.input)
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

func Test_GetFullMessagesByChannelIDAndPage(t *testing.T) {
	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		channelID     int
		page          int
		want          []model.FullMessage
		expectedError error
	}{
		{
			name: "GetFullMessagesByChannelIDAndPage successful",
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "title", "message_url", "image_url",
					"channel_id", "channel_name", "channel_image_url",
					"user_id", "fullname", "user_image_url",
					"count",
				}).
					AddRow(1, "test1", "test1.com", "test1.jpg", 1, "test", "test1.jpg", 1, "test1", "test1.jpg", 1).
					AddRow(2, "test2", "test2.com", "test2.jpg", 1, "test", "test2.jpg", 2, "test2", "test2.jpg", 2)

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.image_url, 
		 			c.id AS channel_id, c.name AS channel_name, c.image_url AS channel_image_url, 
		 			u.id AS user_id, u.fullname, u.image_url AS user_image_url, 
		 			(SELECT COUNT(id) FROM reply WHERE message_id = m.id)
		 			FROM message m 
		 			LEFT JOIN channel c ON c.id = m.channel_id 
		 			LEFT JOIN tg_user u ON u.id = m.user_id
	 	 			WHERE m.channel_id = $1 
		 			ORDER BY count DESC NULLS LAST LIMIT 10 OFFSET $2;`,
				).WithArgs(1, 10).WillReturnRows(rows)
			},
			channelID: 1,
			page:      10,
			want: []model.FullMessage{
				{
					ID: 1, Title: "test1", MessageURL: "test1.com", ImageURL: "test1.jpg",
					ChannelID: 1, ChannelName: "test", ChannelImageURL: "test1.jpg",
					UserID: 1, FullName: "test1", UserImageURL: "test1.jpg", RepliesCount: 1,
				},
				{
					ID: 2, Title: "test2", MessageURL: "test2.com", ImageURL: "test2.jpg",
					ChannelID: 1, ChannelName: "test", ChannelImageURL: "test2.jpg",
					UserID: 2, FullName: "test2", UserImageURL: "test2.jpg", RepliesCount: 2,
				},
			},
		},
		{
			name: "GetFullMessagesByChannelIDAndPage failed with not found messages",
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "title", "message_url", "image_url",
					"channel_id", "channel_name", "channel_image_url",
					"user_id", "fullname", "user_image_url",
					"count",
				})

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.image_url, 
		 			c.id AS channel_id, c.name AS channel_name, c.image_url AS channel_image_url, 
		 			u.id AS user_id, u.fullname, u.image_url AS user_image_url, 
		 			(SELECT COUNT(id) FROM reply WHERE message_id = m.id)
		 			FROM message m 
		 			LEFT JOIN channel c ON c.id = m.channel_id 
		 			LEFT JOIN tg_user u ON u.id = m.user_id
	 	 			WHERE m.channel_id = $1 
		 			ORDER BY count DESC NULLS LAST LIMIT 10 OFFSET $2;`,
				).WithArgs(1, 10).WillReturnRows(rows)
			},
			channelID:     1,
			page:          10,
			expectedError: nil,
		},
		{
			name: "GetFullMessagesByChannelIDAndPage failed with some sql error",
			mock: func() {
				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.image_url, 
		 			c.id AS channel_id, c.name AS channel_name, c.image_url AS channel_image_url, 
		 			u.id AS user_id, u.fullname, u.image_url AS user_image_url, 
		 			(SELECT COUNT(id) FROM reply WHERE message_id = m.id)
		 			FROM message m 
		 			LEFT JOIN channel c ON c.id = m.channel_id 
		 			LEFT JOIN tg_user u ON u.id = m.user_id
	 	 			WHERE m.channel_id = $1 
		 			ORDER BY count DESC NULLS LAST LIMIT 10 OFFSET $2;`,
				).WithArgs(1, 10).WillReturnError(fmt.Errorf("some sql error"))
			},
			channelID:     1,
			page:          10,
			expectedError: fmt.Errorf("get full messages by channel id and page: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessagesByChannelIDAndPage(tt.channelID, tt.page)
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

func Test_GetFullMessagesByUserID(t *testing.T) {
	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         int
		want          []model.FullMessage
		expectedError error
	}{
		{
			name: "GetFullMessagesByUserID successful",
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "title", "message_url", "image_url",
					"channel_id", "channel_name", "channel_title", "channel_image_url",
					"count",
				}).
					AddRow(1, "test1", "test1.com", "test1.jpg", 1, "test1", "test1", "test1.jpg", 1).
					AddRow(2, "test2", "test2.com", "test2.jpg", 2, "test2", "test2", "test2.jpg", 2)

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.image_url, 
		 			c.id AS channel_id, c.name AS channel_name, c.Title AS channel_title, c.image_url AS channel_image_url, 
		 			(SELECT COUNT(id) FROM reply WHERE message_id = m.id)
		 			FROM message m 
		 			LEFT JOIN channel c ON c.id = m.channel_id 
		 			LEFT JOIN tg_user u ON u.id = m.user_id
		 			WHERE m.user_id= $1 
					ORDER BY count DESC NULLS LAST;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want: []model.FullMessage{
				{
					ID: 1, Title: "test1", MessageURL: "test1.com", ImageURL: "test1.jpg",
					ChannelID: 1, ChannelName: "test1", ChannelTitle: "test1", ChannelImageURL: "test1.jpg",
					RepliesCount: 1,
				},
				{
					ID: 2, Title: "test2", MessageURL: "test2.com", ImageURL: "test2.jpg",
					ChannelID: 2, ChannelName: "test2", ChannelTitle: "test2", ChannelImageURL: "test2.jpg",
					RepliesCount: 2,
				},
			},
		},
		{
			name: "GetFullMessagesByUserID failed with not found messages",
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "title", "message_url", "image_url",
					"channel_id", "channel_name", "channel_title", "channel_image_url",
					"count",
				})

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.image_url, 
		 			c.id AS channel_id, c.name AS channel_name, c.Title AS channel_title, c.image_url AS channel_image_url, 
		 			(SELECT COUNT(id) FROM reply WHERE message_id = m.id)
		 			FROM message m 
		 			LEFT JOIN channel c ON c.id = m.channel_id 
		 			LEFT JOIN tg_user u ON u.id = m.user_id
		 			WHERE m.user_id= $1 
					ORDER BY count DESC NULLS LAST;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input:         1,
			expectedError: nil,
		},
		{
			name: "GetFullMessagesByUserID failed with some sql error",
			mock: func() {
				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.image_url, 
		 			c.id AS channel_id, c.name AS channel_name, c.Title AS channel_title, c.image_url AS channel_image_url, 
		 			(SELECT COUNT(id) FROM reply WHERE message_id = m.id)
		 			FROM message m 
		 			LEFT JOIN channel c ON c.id = m.channel_id 
		 			LEFT JOIN tg_user u ON u.id = m.user_id
		 			WHERE m.user_id= $1 
					ORDER BY count DESC NULLS LAST;`,
				).WithArgs(1).WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         1,
			expectedError: fmt.Errorf("get full messages by user id: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessagesByUserID(tt.input)
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

func Test_GetFullMessageByID(t *testing.T) {
	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         int
		want          *model.FullMessage
		expectedError error
	}{
		{
			name: "GetFullMessageByID successful",
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "title", "message_url", "image_url",
					"channel_id", "channel_name", "channel_title", "channel_image_url",
					"user_id", "fullname", "user_image_url",
					"count",
				}).
					AddRow(1, "test1", "test.com", "test.jpg", 1, "test", "test1", "test1.jpg", 1, "test1 test", "test1.jpg", 2)

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.image_url, 
		 			c.id AS channel_id, c.name AS channel_name, c.title as channel_title, c.image_url as channel_image_url, 
		 			u.id as user_id, u.fullname, u.image_url as user_image_url, 
		 			(SELECT COUNT(id) FROM reply WHERE message_id = m.id)
		 			FROM message m 
		 			LEFT JOIN channel c ON c.id = m.channel_id 
		 			LEFT JOIN tg_user u ON u.id = m.user_id
		 			WHERE m.id = $1;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want: &model.FullMessage{
				ID: 1, Title: "test1", MessageURL: "test.com", ImageURL: "test.jpg",
				ChannelID: 1, ChannelName: "test", ChannelTitle: "test1", ChannelImageURL: "test1.jpg",
				UserID: 1, FullName: "test1 test", UserImageURL: "test1.jpg",
				RepliesCount: 2,
			},
		},
		{
			name: "GetFullMessageByID failed with not found message",
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "title", "message_url", "image_url",
					"channel_id", "channel_name", "channel_title", "channel_image_url",
					"user_id", "fullname", "user_image_url",
					"count",
				})

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.image_url, 
		 			c.id AS channel_id, c.name AS channel_name, c.title as channel_title, c.image_url as channel_image_url, 
		 			u.id as user_id, u.fullname, u.image_url as user_image_url, 
		 			(SELECT COUNT(id) FROM reply WHERE message_id = m.id)
		 			FROM message m 
		 			LEFT JOIN channel c ON c.id = m.channel_id 
		 			LEFT JOIN tg_user u ON u.id = m.user_id
		 			WHERE m.id = $1;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input:         1,
			expectedError: nil,
		},
		{
			name: "GetFullMessageByID failed with some sql error",
			mock: func() {
				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.image_url, 
		 			c.id AS channel_id, c.name AS channel_name, c.title as channel_title, c.image_url as channel_image_url, 
		 			u.id as user_id, u.fullname, u.image_url as user_image_url, 
		 			(SELECT COUNT(id) FROM reply WHERE message_id = m.id)
		 			FROM message m 
		 			LEFT JOIN channel c ON c.id = m.channel_id 
		 			LEFT JOIN tg_user u ON u.id = m.user_id
		 			WHERE m.id = $1;`,
				).WithArgs(1).WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         1,
			expectedError: fmt.Errorf("get full message by id: %w", fmt.Errorf("some sql error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessageByID(tt.input)
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
