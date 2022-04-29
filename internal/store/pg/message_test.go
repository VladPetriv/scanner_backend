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
		want    []model.FullMessage
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "name", "fullname", "photourl", "count"}).
					AddRow(1, "test1", "test1", "test1", "test1.jpg", 1).
					AddRow(2, "test2", "test2", "test2", "test2.jpg", 3)

				mock.ExpectQuery("SELECT m.id, m.Title, c.Name, u.Fullname, u.Photourl, (SELECT COUNT(id) FROM replie WHERE message_id = m.id)  FROM message m LEFT JOIN channel c ON c.id = m.channel_id LEFT JOIN tg_user u ON u.id = m.user_id ORDER BY m.id;").
					WillReturnRows(rows)
			},
			want: []model.FullMessage{
				{ID: 1, Title: "test1", ChannelName: "test1", FullName: "test1", PhotoURL: "test1.jpg", ReplieCount: 1},
				{ID: 2, Title: "test2", ChannelName: "test2", FullName: "test2", PhotoURL: "test2.jpg", ReplieCount: 3},
			},
		},
		{
			name: "full messages not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "name", "fullname", "photourl"})

				mock.ExpectQuery("SELECT m.id, m.Title, c.Name, u.Fullname, u.Photourl, (SELECT COUNT(id) FROM replie WHERE message_id = m.id)  FROM message m LEFT JOIN channel c ON c.id = m.channel_id LEFT JOIN tg_user u ON u.id = m.user_id ORDER BY m.id;").
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetFullMessages()
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
