package pg

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/VladPetriv/scanner_backend/internal/model"
	"github.com/VladPetriv/scanner_backend/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestMessagePg_GetMessage(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewMessageRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    *model.Message
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "channel_id", "title"}).
					AddRow(1, 1, 1, "test")

				mock.ExpectQuery("SELECT * FROM message WHERE id=$1;").
					WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  &model.Message{ID: 1, UserID: 1, ChannelID: 1, Title: "test"},
		},
		{
			name: "message not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "channel_id", "title"})

				mock.ExpectQuery("SELECT * FROM message WHERE id=$1;").
					WithArgs(404).WillReturnRows(rows)
			},
			input:   404,
			wantErr: true,
		},
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "channel_id", "title"})

				mock.ExpectQuery("SELECT * FROM message WHERE id=$1;").
					WithArgs().WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetMessage(tt.input)
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

func TestMessagePg_GetMessages(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewMessageRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		want    []model.Message
		wantErr bool
	}{
		{
			name: "ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "channel_id", "title"}).
					AddRow(1, 1, 1, "test1").
					AddRow(2, 2, 2, "test2")

				mock.ExpectQuery("SELECT * FROM message;").
					WillReturnRows(rows)
			},
			want: []model.Message{
				{ID: 1, UserID: 1, ChannelID: 1, Title: "test1"},
				{ID: 2, UserID: 2, ChannelID: 2, Title: "test2"},
			},
		},
		{
			name: "message not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "channel_id", "title"})

				mock.ExpectQuery("SELECT * FROM message;").
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetMessages()
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

func TestMessagePg_GetMessageByName(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewMessageRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    *model.Message
		wantErr bool
	}{
		{
			name: "ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "channel_id", "title"}).
					AddRow(1, 1, 1, "test1")

				mock.ExpectQuery("SELECT * FROM message WHERE title=$1;").
					WithArgs("test1").WillReturnRows(rows)
			},
			input: "test1",
			want:  &model.Message{ID: 1, UserID: 1, ChannelID: 1, Title: "test1"},
		},
		{
			name: "message not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "channel_id", "title"})

				mock.ExpectQuery("SELECT * FROM message WHERE title=$1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "channel_id", "title"})

				mock.ExpectQuery("SELECT * FROM message WHERE title=$1;").
					WithArgs().WillReturnRows(rows)
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetMessageByName(tt.input)
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

	r := NewMessageRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    []model.FullMessage
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "id", "name", "channelPhotoUrl", "id", "fullname", "photourl", "count"}).
					AddRow(1, "test1", 1, "test1", "test1.jpg", 1, "test1", "test1.jpg", 1).
					AddRow(2, "test2", 2, "test2", "test2.jpg", 2, "test2", "test2.jpg", 3)

				mock.ExpectQuery(
					`SELECT m.id, m.Title, c.id, c.Name, c.Photourl as channelPhotoUrl, u.id, u.Fullname, u.Photourl, (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m LEFT JOIN channel c ON c.id = m.channel_id LEFT JOIN tg_user u ON u.id = m.user_id
					ORDER BY m.id DESC NULLS LAST LIMIT 10 OFFSET $1;`,
				).WillReturnRows(rows)
			},
			want: []model.FullMessage{
				{ID: 1, Title: "test1", ChannelID: 1, ChannelName: "test1", ChannelPhotoURL: "test1.jpg", UserID: 1, FullName: "test1", PhotoURL: "test1.jpg", ReplieCount: 1},
				{ID: 2, Title: "test2", ChannelID: 2, ChannelName: "test2", ChannelPhotoURL: "test2.jpg", UserID: 2, FullName: "test2", PhotoURL: "test2.jpg", ReplieCount: 3},
			},
			input: 10,
		},
		{
			name: "full messages not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "id", "name", "channelPhotoUrl", "id", "fullname", "photourl"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, c.id, c.Name, c.Photourl as channelPhotoUrl, u.id, u.Fullname, u.Photourl, (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m LEFT JOIN channel c ON c.id = m.channel_id LEFT JOIN tg_user u ON u.id = m.user_id
					ORDER BY m.id DESC NULLS LAST LIMIT 10 OFFSET $1;`,
				).WillReturnRows(rows)
			},
			input:   10,
			wantErr: true,
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "id", "name", "channelPhotoUrl", "id", "fullname", "photourl"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, c.id, c.Name, c.Photourl as channelPhotoUrl, u.id, u.Fullname, u.Photourl, (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m LEFT JOIN channel c ON c.id = m.channel_id LEFT JOIN tg_user u ON u.id = m.user_id
					ORDER BY m.id DESC NULLS LAST LIMIT 10 OFFSET $1;`,
				).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessages(tt.input)
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

func TestMessagePg_GetFullMessagesByChannelID(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewMessageRepo(&DB{DB: db})

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
				rows := sqlmock.NewRows([]string{"id", "title", "id", "name", "channelPhotoUrl", "id", "fullname", "photourl", "count"}).
					AddRow(1, "test1", 1, "test", "test1.jpg", 1, "test1", "test1.jpg", 1).
					AddRow(2, "test2", 1, "test", "test2.jpg", 2, "test2", "test2.jpg", 2)

				mock.ExpectQuery(
					`SELECT m.id, m.Title, c.id, c.Name, c.Photourl as channelPhotoUrl, u.id, u.Fullname, u.Photourl, (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m LEFT JOIN channel c ON c.id = m.channel_id LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.channel_id = $1 ORDER BY count DESC NULLS LAST LIMIT $2 OFFSET $3;`,
				).WithArgs(1, 10, 100).WillReturnRows(rows)

			},
			ID:    1,
			page:  10,
			limit: 100,
			want: []model.FullMessage{
				{ID: 1, Title: "test1", ChannelID: 1, ChannelName: "test", ChannelPhotoURL: "test1.jpg", UserID: 1, FullName: "test1", PhotoURL: "test1.jpg", ReplieCount: 1},
				{ID: 2, Title: "test2", ChannelID: 1, ChannelName: "test", ChannelPhotoURL: "test2.jpg", UserID: 2, FullName: "test2", PhotoURL: "test2.jpg", ReplieCount: 2},
			},
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "id", "name", "channelPhotoUrl", "id", "fullname", "photourl"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, c.id, c.Name, c.Photourl as channelPhotoUrl, u.id, u.Fullname, u.Photourl, (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m LEFT JOIN channel c ON c.id = m.channel_id LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.channel_id = $1 ORDER BY count DESC NULLS LAST LIMIT $2 OFFSET $3;`,
				).WithArgs().WillReturnRows(rows)

			},
			wantErr: true,
		},
		{
			name: "messages not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "id", "name", "channelPhotoUrl", "id", "fullname", "photourl"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, c.id, c.Name, c.Photourl as channelPhotoUrl, u.id, u.Fullname, u.Photourl, (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m LEFT JOIN channel c ON c.id = m.channel_id LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.channel_id = $1 ORDER BY count DESC NULLS LAST LIMIT $2 OFFSET $3;`,
				).WithArgs(404, 10, 100).WillReturnRows(rows)
			},
			ID:      404,
			limit:   100,
			page:    10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessagesByChannelID(tt.ID, tt.page, tt.limit)
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

func TestMessagePg_GetFullMessagesByUserID(t *testing.T) {
	db, mock, err := util.CreateMock()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewMessageRepo(&DB{DB: db})

	tests := []struct {
		name    string
		mock    func()
		input   int
		want    []model.FullMessage
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "id", "name", "title", "channelPhotoUrl", "count"}).
					AddRow(1, "test1", 1, "test1", "test1", "test1.jpg", 1).
					AddRow(2, "test2", 2, "test2", "test2", "test2.jpg", 2)

				mock.ExpectQuery(
					`SELECT m.id, m.Title, c.id, c.Name, c.Title, c.Photourl as channelPhotoUrl, (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m LEFT JOIN channel c ON c.id = m.channel_id LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.user_id= $1 ORDER BY count DESC NULLS LAST;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want: []model.FullMessage{
				{ID: 1, Title: "test1", ChannelID: 1, ChannelName: "test1", ChannelTitle: "test1", ChannelPhotoURL: "test1.jpg", ReplieCount: 1},
				{ID: 2, Title: "test2", ChannelID: 2, ChannelName: "test2", ChannelTitle: "test2", ChannelPhotoURL: "test2.jpg", ReplieCount: 2},
			},
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "id", "name", "title", "channelPhotoUrl", "count"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, c.id, c.Name, c.Title, c.Photourl as channelPhotoUrl, (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m LEFT JOIN channel c ON c.id = m.channel_id LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.user_id= $1 ORDER BY count DESC NULLS LAST;`,
				).WithArgs().WillReturnRows(rows)

			},
			wantErr: true,
		},
		{
			name: "messages not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "id", "name", "title", "channelPhotoUrl", "count"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, c.id, c.Name, c.Title, c.Photourl as channelPhotoUrl, (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m LEFT JOIN channel c ON c.id = m.channel_id LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.user_id= $1 ORDER BY count DESC NULLS LAST;`,
				).WithArgs(1).WillReturnRows(rows)
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

func TestMessagePg_GetFullMessageByMessageID(t *testing.T) {
	db, mock, err := util.CreateMock()

	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewMessageRepo(&DB{DB: db})

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
				rows := sqlmock.NewRows([]string{"id", "title", "id", "name", "title", "channelPhotoUrl", "id", "fullname", "photourl", "count"}).
					AddRow(1, "test1", 1, "test", "test1", "test1.jpg", 1, "test1 test", "test1.jpg", 2)

				mock.ExpectQuery(
					`SELECT m.id, m.Title, c.id, c.Name, c.Title, c.Photourl as channelPhotoUrl, u.id, u.Fullname, u.Photourl, (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m LEFT JOIN channel c ON c.id = m.channel_id LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.id = $1;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  &model.FullMessage{ID: 1, Title: "test1", ChannelID: 1, ChannelName: "test", ChannelTitle: "test1", ChannelPhotoURL: "test1.jpg", UserID: 1, FullName: "test1 test", PhotoURL: "test1.jpg", ReplieCount: 2},
		},
		{
			name: "empty field",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "id", "name", "title", "channelPhotoUrl", "id", "fullname", "photourl", "count"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, c.id, c.Name, c.Title, c.Photourl as channelPhotoUrl, u.id, u.Fullname, u.Photourl, (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m LEFT JOIN channel c ON c.id = m.channel_id LEFT JOIN tg_user u ON u.id = m.user_id
					WHERE m.id = $1;`,
				).WithArgs().WillReturnRows(rows)

			},
			want: nil,
		},
		{
			name: "message not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "id", "name", "title", "channelPhotoUrl", "id", "fullname", "photourl", "count"})

				mock.ExpectQuery(
					`SELECT m.id, m.Title, c.id, c.Name, c.Title, c.Photourl as channelPhotoUrl, u.id, u.Fullname, u.Photourl, (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
					FROM message m LEFT JOIN channel c ON c.id = m.channel_id LEFT JOIN tg_user u ON u.id = m.user_id
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
