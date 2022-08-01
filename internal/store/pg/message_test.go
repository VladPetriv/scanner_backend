package pg_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/internal/store/pg"
	"github.com/VladPetriv/scanner_backend/pkg/util"
)

func Test_CreateMessage(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

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

func Test_GetMessagesCount(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		want    int
		wantErr bool
	}{
		{
			name: "Ok: [messages count found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).
					AddRow(10)

				mock.ExpectQuery("SELECT COUNT(*) FROM message;").
					WillReturnRows(rows)
			},
			want: 10,
		},
		{
			name: "Error: [message not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"})

				mock.ExpectQuery("SELECT COUNT(*) FROM message;").
					WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery("SELECT COUNT(*) FROM message;").
					WillReturnError(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetMessagesCount()
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

func Test_GetMessagesCountByChannelID(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    int
		wantErr bool
	}{
		{
			name: "Ok: [messages count found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).
					AddRow(2)

				mock.ExpectQuery("SELECT COUNT(*) FROM message m LEFT JOIN channel c ON c.id = m.channel_id WHERE m.channel_id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  2,
		},
		{
			name: "Error: [message count not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})

				mock.ExpectQuery("SELECT COUNT(*) FROM message m LEFT JOIN channel c ON c.id = m.channel_id WHERE m.channel_id = $1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input:   1,
			wantErr: true,
		},
		{
			name: "Error: [pq error]",
			mock: func() {
				mock.ExpectQuery("SELECT COUNT(*) FROM message m LEFT JOIN channel c ON c.id = m.channel_id WHERE m.channel_id = $1;").
					WithArgs(1).WillReturnError(sqlmock.ErrCancelled)
			},
			input:   1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt.mock()

		got, err := r.GetMessagesCountByChannelID(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}

		assert.NoError(t, mock.ExpectationsWereMet())
	}
}

func Test_GetFullMessagesByPage(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    []model.FullMessage
		wantErr bool
	}{
		{
			name: "Ok: [full messages found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "channelid", "channelname", "channelimageurl", "userid", "fullname", "userimageurl", "count"}).
					AddRow(1, "test1", "test1.com", "test1.jpg", 1, "test1", "test1.jpg", 1, "test1", "test1.jpg", 1).
					AddRow(2, "test2", "test2.com", "test2.jpg", 2, "test2", "test2.jpg", 2, "test2", "test2.jpg", 3)

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.imageurl, 
					c.id AS channelid, c.name AS channelname, c.imageurl AS channelimageurl, 
					u.id AS userid, u.fullname, u.imageurl AS userimageurl, 
					(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					ORDER BY m.id DESC NULLS LAST LIMIT 10 OFFSET $1;`,
				).WithArgs(10).WillReturnRows(rows)
			},
			input: 10,
			want: []model.FullMessage{
				{ID: 1, Title: "test1", MessageURL: "test1.com", ImageURL: "test1.jpg", ChannelID: 1, ChannelName: "test1", ChannelImageURL: "test1.jpg", UserID: 1, FullName: "test1", UserImageURL: "test1.jpg", ReplieCount: 1},
				{ID: 2, Title: "test2", MessageURL: "test2.com", ImageURL: "test2.jpg", ChannelID: 2, ChannelName: "test2", ChannelImageURL: "test2.jpg", UserID: 2, FullName: "test2", UserImageURL: "test2.jpg", ReplieCount: 3},
			},
		},
		{
			name: "Error: [full messages not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "channelid", "channelname", "channelimageurl", "userid", "fullname", "userimageurl"})

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.imageurl, 
					c.id AS channelid, c.name AS channelname, c.imageurl AS channelimageurl, 
					u.id AS userid, u.fullname, u.imageurl AS userimageurl, 
					(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					ORDER BY m.id DESC NULLS LAST LIMIT 10 OFFSET $1;`,
				).WithArgs(10).WillReturnRows(rows)
			},
			input:   10,
			wantErr: true,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.imageurl, 
					c.id AS channelid, c.name AS channelname, c.imageurl AS channelimageurl, 
					u.id AS userid, u.fullname, u.imageurl AS userimageurl, 
					(SELECT COUNT(*) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					ORDER BY m.id DESC NULLS LAST LIMIT 10 OFFSET $1;`,
				).WithArgs(10).WillReturnError(fmt.Errorf("some error"))
			},
			input:   10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessagesByPage(tt.input)
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

func Test_GetFullMessagesByChannelIDAndPage(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name      string
		mock      func()
		channelID int
		page      int
		want      []model.FullMessage
		wantErr   bool
	}{
		{
			name: "Ok: [full messages found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "channelid", "channelname", "channelimageurl", "userid", "fullname", "userimageurl", "count"}).
					AddRow(1, "test1", "test1.com", "test1.jpg", 1, "test", "test1.jpg", 1, "test1", "test1.jpg", 1).
					AddRow(2, "test2", "test2.com", "test2.jpg", 1, "test", "test2.jpg", 2, "test2", "test2.jpg", 2)

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.imageurl, 
					c.id AS channelid, c.name AS channelname, c.imageurl AS channelimageurl, 
					u.id AS userid, u.fullname, u.imageurl AS userimageurl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
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
			name: "Error: [full messages not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "channelid", "channelname", "channelimageurl", "userid", "fullname", "userimageurl"})

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.imageurl, 
					c.id AS channelid, c.name AS channelname, c.imageurl AS channelimageurl, 
					u.id AS userid, u.fullname, u.imageurl AS userimageurl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.channel_id = $1 
					ORDER BY count DESC NULLS LAST LIMIT 10 OFFSET $2;`,
				).WithArgs(1, 10).WillReturnRows(rows)
			},
			channelID: 1,
			page:      10,
			wantErr:   true,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.imageurl, 
					c.id AS channelid, c.name AS channelname, c.imageurl AS channelimageurl, 
					u.id AS userid, u.fullname, u.imageurl AS userimageurl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.channel_id = $1 
					ORDER BY count DESC NULLS LAST LIMIT 10 OFFSET $2;`,
				).WithArgs(1, 10).WillReturnError(fmt.Errorf("some error"))
			},
			channelID: 1,
			page:      10,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessagesByChannelIDAndPage(tt.channelID, tt.page)
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

func Test_GetFullMessagesByUserID(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    []model.FullMessage
		wantErr bool
	}{
		{
			name: "Ok: [full messages found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "channelid", "channelname", "channeltitle", "channelimageurl", "count"}).
					AddRow(1, "test1", "test1.com", "test1.jpg", 1, "test1", "test1", "test1.jpg", 1).
					AddRow(2, "test2", "test2.com", "test2.jpg", 2, "test2", "test2", "test2.jpg", 2)

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.imageurl, 
					c.id AS channelid, c.name AS channelname, c.Title AS channeltitle, c.imageurl AS channelimageurl, 
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
			name: "Error: [full messages not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "channelid", "channelname", "channeltitle", "channelimageurl", "count"})

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.imageurl, 
					c.id AS channelid, c.name AS channelname, c.Title AS channeltitle, c.imageurl AS channelimageurl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.user_id= $1 
					ORDER BY count DESC NULLS LAST;`,
				).WithArgs(1).WillReturnRows(rows)

			},
			input:   1,
			wantErr: true,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.imageurl, 
					c.id AS channelid, c.name AS channelname, c.Title AS channeltitle, c.imageurl AS channelimageurl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.user_id= $1 
					ORDER BY count DESC NULLS LAST;`,
				).WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			input:   1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessagesByUserID(tt.input)
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

func Test_GetFullMessageByMessageID(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewMessageRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    *model.FullMessage
		wantErr bool
	}{
		{
			name: "Ok: [message found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "channelid", "channelname", "channeltitle", "channelimageurl", "userid", "fullname", "userimageurl", "count"}).
					AddRow(1, "test1", "test.com", "test.jpg", 1, "test", "test1", "test1.jpg", 1, "test1 test", "test1.jpg", 2)

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.imageurl, 
					c.id AS channelid, c.name AS channelname, c.title as channeltitle, c.imageurl as channelimageurl, 
					u.id as userid, u.fullname, u.imageurl as userimageurl, 
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
			name: "Error: [message not found]",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "message_url", "imageurl", "channelid", "channelname", "channeltitle", "channelimageurl", "userid", "fullname", "userimageurl", "count"})

				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.imageurl, 
					c.id AS channelid, c.name AS channelname, c.title as channeltitle, c.imageurl as channelimageurl, 
					u.id as userid, u.fullname, u.imageurl as userimageurl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.id = $1;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input:   1,
			wantErr: true,
		},
		{
			name: "Error: [some sql error]",
			mock: func() {
				mock.ExpectQuery(
					`SELECT m.id, m.title, m.message_url, m.imageurl, 
					c.id AS channelid, c.name AS channelname, c.title as channeltitle, c.imageurl as channelimageurl, 
					u.id as userid, u.fullname, u.imageurl as userimageurl, 
					(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m 
					LEFT JOIN channel c ON c.id = m.channel_id 
					LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.id = $1;`,
				).WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			input:   1,
			wantErr: true,
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
