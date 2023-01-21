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

func Test_CreateReply(t *testing.T) {
	t.Parallel()

	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewReplyRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         *model.DBReply
		expecterError error
	}{
		{
			name: "CreateReply successful",
			mock: func() {
				mock.ExpectExec(`
					INSERT INTO reply(user_id, message_id, title, image_url) VALUES ($1, $2, $3, $4);`,
				).WithArgs(1, 1, "test", "test.jpg").WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: &model.DBReply{UserID: 1, MessageID: 1, Title: "test", ImageURL: "test.jpg"},
		},
		{
			name: "CreateReply failed with some sql error",
			mock: func() {
				mock.ExpectExec(`
					INSERT INTO reply(user_id, message_id, title, image_url) VALUES ($1, $2, $3, $4);`,
				).WithArgs(1, 1, "test", "test.jpg").WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         &model.DBReply{UserID: 1, MessageID: 1, Title: "test", ImageURL: "test.jpg"},
			expecterError: fmt.Errorf("create reply: %w", fmt.Errorf("some sql error")),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.mock()

			err = r.CreateReply(tt.input)
			if tt.expecterError != nil {
				assert.Error(t, err)
				assert.EqualValues(t, tt.expecterError, err)
			} else {
				assert.NoError(t, err)
				assert.Nil(t, err)
			}
		})
	}

	t.Cleanup(func() {
		db.Close()
	})
}

func Test_GetFullRepliesByMessageID(t *testing.T) {
	t.Parallel()

	db, mock, err := mocks.CreateMock()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "postgres")

	r := pg.NewReplyRepo(&pg.DB{DB: sqlxDB})

	tests := []struct {
		name          string
		mock          func()
		input         int
		want          []model.FullReply
		expectedError error
	}{
		{
			name: "GetFullRepliesByMessageID successful",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "image_url", "user_id", "fullname", "user_image_url"}).
					AddRow(1, "test1", "test_reply1.jpg", 1, "test1 test", "test1.jpg").
					AddRow(2, "test2", "test_reply2.jpg", 2, "test2 test", "test2.jpg")

				mock.ExpectQuery(
					`SELECT r.id, r.title, r.image_url, 
					u.id as user_id, u.fullname, u.image_url as user_image_url
					FROM reply r 
					LEFT JOIN tg_user u ON r.user_id = u.id 
					WHERE r.message_id = $1 
					ORDER BY r.id DESC NULLS LAST;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want: []model.FullReply{
				{ID: 1, Title: "test1", ImageURL: "test_reply1.jpg", UserID: 1, FullName: "test1 test", UserImageURL: "test1.jpg"},
				{ID: 2, Title: "test2", ImageURL: "test_reply2.jpg", UserID: 2, FullName: "test2 test", UserImageURL: "test2.jpg"},
			},
		},
		{
			name: "GetFullRepliesByMessageID failed with not found replies",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "image_url", "user_id", "fullname", "user_image_url"})

				mock.ExpectQuery(
					`SELECT r.id, r.title, r.image_url, 
		 			u.id as user_id, u.fullname, u.image_url as user_image_url
		 			FROM reply r 
		 			LEFT JOIN tg_user u ON r.user_id = u.id 
		 			WHERE r.message_id = $1 
		 			ORDER BY r.id DESC NULLS LAST;`,
				).WithArgs(1).WillReturnRows(rows)
			},
			input: 1,
			want:  nil,
		},
		{
			name: "GetFullRepliesByMessageID failed with some sql error",
			mock: func() {
				mock.ExpectQuery(
					`SELECT r.id, r.title, r.image_url, 
					u.id as user_id, u.fullname, u.image_url as user_image_url
					FROM reply r 
					LEFT JOIN tg_user u ON r.user_id = u.id 
					WHERE r.message_id = $1 
					ORDER BY r.id DESC NULLS LAST;`,
				).WithArgs(1).WillReturnError(fmt.Errorf("some sql error"))
			},
			input:         1,
			expectedError: fmt.Errorf("get full replies by message id: %w", fmt.Errorf("some sql error")),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.mock()

			got, err := r.GetFullRepliesByMessageID(tt.input)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}

	t.Cleanup(func() {
		db.Close()
	})
}
