package pg_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
	"github.com/VladPetriv/scanner_backend/pkg/util"
)

func TestMesasgePg_CreateMessage(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewMessageRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   *model.DBMessage
		want    int
		wantErr bool
	}{
		{
			name: "Ok: [message created]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1)

				mock.ExpectQuery(`
					INSERT INTO message(channel_id, user_id, title, message_url, imageurl) 
					VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
				).WithArgs(1, 1, "test", "test.url", "test.jpg").WillReturnRows(rows)
			},
			input: &model.DBMessage{ChannelID: 1, UserID: 1, Title: "test", MessageURL: "test.url", ImageURL: "test.jpg"},
			want:  1,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery(`
					INSERT INTO message(channel_id, user_id, title, message_url, imageurl) 
					VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
				).WithArgs(1, 1, "test", "test.url", "test.jpg").WillReturnError(fmt.Errorf("some sql error"))
			},
			input:   &model.DBMessage{ChannelID: 1, UserID: 1, Title: "test", MessageURL: "test.url", ImageURL: "test.jpg"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateMessage(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

}

func TestMessagePg_GetMessagesLength(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewMessageRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		want    int
		wantErr bool
	}{
		{
			name: "ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).
					AddRow(10)

				mock.ExpectQuery("SELECT COUNT(*) FROM message;").
					WillReturnRows(rows)
			},
			want: 10,
		},
		{
			name: "message not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"})

				mock.ExpectQuery("SELECT COUNT(*) FROM message;").
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetMessagesLength()
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

func TestMessagePg_GetFullMessages(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewMessageRepo(&pg.DB{DB: db})

	tests := []struct {
		name  string
		mock  func()
		input int
		want  []model.FullMessage
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "id", "name", "channelImageUrl", "id", "fullname", "userImageUrl", "count"}).
					AddRow(1, "test1", "test1.com", "test1.jpg", 1, "test1", "test1.jpg", 1, "test1", "test1.jpg", 1).
					AddRow(2, "test2", "test2.com", "test2.jpg", 2, "test2", "test2.jpg", 2, "test2", "test2.jpg", 3)

				mock.ExpectQuery(
					`SELECT m.id, m.Title, m.message_url, m.imageurl, 
					c.id, c.Name, c.imageurl as channelImageUrl, 
					u.id, u.Fullname, u.imageurl as userImageUrl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					ORDER BY m.id DESC NULLS LAST LIMIT 10 OFFSET $1;`,
				).WillReturnRows(rows)
			},
			want: []model.FullMessage{
				{ID: 1, Title: "test1", MessageURL: "test1.com", ImageURL: "test1.jpg", ChannelID: 1, ChannelName: "test1", ChannelImageURL: "test1.jpg", UserID: 1, FullName: "test1", UserImageURL: "test1.jpg", ReplieCount: 1},
				{ID: 2, Title: "test2", MessageURL: "test2.com", ImageURL: "test2.jpg", ChannelID: 2, ChannelName: "test2", ChannelImageURL: "test2.jpg", UserID: 2, FullName: "test2", UserImageURL: "test2.jpg", ReplieCount: 3},
			},
			input: 10,
		},
		{
			name: "full messages not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "id", "name", "channelImageUrl", "id", "fullname", "userImageUrl"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, m.message_url, m.imageurl, 
					c.id, c.Name, c.imageurl as channelImageUrl, 
					u.id, u.Fullname, u.imageurl as userImageUrl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					ORDER BY m.id DESC NULLS LAST LIMIT 10 OFFSET $1;`,
				).WillReturnRows(rows)
			},
			input: 10,
			want:  nil,
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "id", "name", "channelImageUrl", "id", "fullname", "userImageurl"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, m.message_url, m.imageurl, 
					c.id, c.Name, c.imageurl as channelImageUrl, 
					u.id, u.Fullname, u.imageurl as userImageUrl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					ORDER BY m.id DESC NULLS LAST LIMIT 10 OFFSET $1;`,
				).WillReturnRows(rows)
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessages(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMessagePg_GetFullMessagesByChannelID(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewMessageRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		ID      int
		page    int
		limit   int
		want    []model.FullMessage
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "id", "name", "channelImageUrl", "id", "fullname", "userImageUrl", "count"}).
					AddRow(1, "test1", "test1.com", "test1.jpg", 1, "test", "test1.jpg", 1, "test1", "test1.jpg", 1).
					AddRow(2, "test2", "test2.com", "test2.jpg", 1, "test", "test2.jpg", 2, "test2", "test2.jpg", 2)

				mock.ExpectQuery(
					`SELECT m.id, m.Title, m.message_url, m.imageurl, 
					c.id, c.Name, c.imageurl as channelImageUrl, 
					u.id, u.Fullname, u.imageurl as userImageUrl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.channel_id = $1 
					ORDER BY count DESC NULLS LAST LIMIT $2 OFFSET $3;`,
				).WithArgs(1, 10, 100).WillReturnRows(rows)

			},
			ID:    1,
			page:  10,
			limit: 100,
			want: []model.FullMessage{
				{
					ID: 1, Title: "test1", MessageURL: "test1.com", ImageURL: "test1.jpg",
					ChannelID: 1, ChannelName: "test", ChannelImageURL: "test1.jpg",
					UserID: 1, FullName: "test1", UserImageURL: "test1.jpg", ReplieCount: 1,
				},
				{
					ID: 2, Title: "test2", MessageURL: "test2.com", ImageURL: "test2.jpg",
					ChannelID: 1, ChannelName: "test", ChannelImageURL: "test2.jpg",
					UserID: 2, FullName: "test2", UserImageURL: "test2.jpg", ReplieCount: 2,
				},
			},
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "id", "name", "channelImageUrl", "id", "fullname", "userImageUrl"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, m.message_url, m.imageurl, 
					c.id, c.Name, c.imageurl as channelImageUrl, 
					u.id, u.Fullname, u.imageurl as userImageUrl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.channel_id = $1 
					ORDER BY count DESC NULLS LAST LIMIT $2 OFFSET $3;`,
				).WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "messages not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "id", "name", "channelImageUrl", "id", "fullname", "userImageUrl"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, m.message_url, m.imageurl, 
					c.id, c.Name, c.imageurl as channelImageUrl, 
					u.id, u.Fullname, u.imageurl as userImageUrl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.channel_id = $1 
					ORDER BY count DESC NULLS LAST LIMIT $2 OFFSET $3;`,
				).WithArgs(404, 10, 100).WillReturnRows(rows)
			},
			ID:    404,
			limit: 100,
			page:  10,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessagesByChannelID(tt.ID, tt.page, tt.limit)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMessagePg_GetFullMessagesByUserID(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewMessageRepo(&pg.DB{DB: db})

	tests := []struct {
		name  string
		mock  func()
		input int
		want  []model.FullMessage
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "id", "name", "title", "channelImageUrl", "count"}).
					AddRow(1, "test1", "test1.com", "test1.jpg", 1, "test1", "test1", "test1.jpg", 1).
					AddRow(2, "test2", "test2.com", "test2.jpg", 2, "test2", "test2", "test2.jpg", 2)

				mock.ExpectQuery(
					`SELECT m.id, m.Title, m.message_url, m.imageurl, 
					c.id, c.Name, c.Title, c.imageurl as channelImageUrl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.user_id= $1 
					ORDER BY count DESC NULLS LAST;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want: []model.FullMessage{
				{ID: 1, Title: "test1", MessageURL: "test1.com", ImageURL: "test1.jpg", ChannelID: 1, ChannelName: "test1", ChannelTitle: "test1", ChannelImageURL: "test1.jpg", ReplieCount: 1},
				{ID: 2, Title: "test2", MessageURL: "test2.com", ImageURL: "test2.jpg", ChannelID: 2, ChannelName: "test2", ChannelTitle: "test2", ChannelImageURL: "test2.jpg", ReplieCount: 2},
			},
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "id", "name", "title", "channelImageUrl", "count"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, m.message_url, m.imageurl, 
					c.id, c.Name, c.Title, c.imageurl as channelImageUrl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.user_id= $1 
					ORDER BY count DESC NULLS LAST;`,
				).WithArgs().WillReturnRows(rows)

			},
			want: nil,
		},
		{
			name: "messages not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "id", "name", "title", "channelImageUrl", "count"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, m.message_url, m.imageurl, 
					c.id, c.Name, c.Title, c.imageurl as channelImageUrl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.user_id= $1 
					ORDER BY count DESC NULLS LAST;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessagesByUserID(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMessagePg_GetFullMessageByMessageID(t *testing.T) {
	db, mock, err := util.CreateMock()

	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewMessageRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    *model.FullMessage
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "id", "name", "title", "channelImageUrl", "id", "fullname", "userImageUrl", "count"}).
					AddRow(1, "test1", "test.com", "test.jpg", 1, "test", "test1", "test1.jpg", 1, "test1 test", "test1.jpg", 2)

				mock.ExpectQuery(
					`SELECT m.id, m.Title, m.message_url, m.imageurl, 
					c.id, c.Name, c.Title, c.imageurl as channelImageUrl, 
					u.id, u.Fullname, u.imageurl as userImageUrl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.id = $1;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  &model.FullMessage{ID: 1, Title: "test1", MessageURL: "test.com", ImageURL: "test.jpg", ChannelID: 1, ChannelName: "test", ChannelTitle: "test1", ChannelImageURL: "test1.jpg", UserID: 1, FullName: "test1 test", UserImageURL: "test1.jpg", ReplieCount: 2},
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "id", "name", "title", "channelImageUrl", "id", "fullname", "userImageUrl", "count"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, m.message_url, m.imageurl, 
					c.id, c.Name, c.Title, c.imageurl as channelImageUrl, 
					u.id, u.Fullname, u.imageurl as userImageUrl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.id = $1;`,
				).WithArgs().WillReturnRows(rows)

			},
			want: nil,
		},
		{
			name: "message not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "id", "name", "title", "channelImageUrl", "id", "fullname", "userImageUrl", "count"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, m.message_url, m.imageurl, 
					c.id, c.Name, c.Title, c.imageurl as channelImageUrl, 
					u.id, u.Fullname, u.imageurl as userImageUrl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.id = $1;`,
				).WithArgs(404).WillReturnRows(rows)
			},
			input: 404,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessageByMessageID(tt.input)
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

func TestMessagePg_GetMessagesLengthByChannelID(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := pg.NewMessageRepo(&pg.DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    int
		wantErr bool
	}{
		{
			name: "Ok: [Messages found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1).
					AddRow(2)

				mock.ExpectQuery("SELECT m.id FROM message m LEFT JOIN channel c ON c.id = m.channel_id WHERE m.channel_id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  2,
		},
		{
			name: "Error: [Message not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("SELECT m.id FROM message m LEFT JOIN channel c ON c.id = m.channel_id WHERE m.channel_id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  0,
		},
		{
			name: "Error: [pq error]",
			mock: func() {
				mock.ExpectQuery("SELECT m.id FROM message m LEFT JOIN channel c ON c.id = m.channel_id WHERE m.channel_id = $1;").
					WithArgs(1).WillReturnError(sqlmock.ErrCancelled)
			},
			input:   1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt.mock()

		got, err := r.GetMessagesLengthByChannelID(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		assert.NoError(t, mock.ExpectationsWereMet())
	}
}
