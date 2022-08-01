package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

var (
	ErrMessagesCountNotFound = errors.New("messages count not found")
	ErrMessagesNotFound      = errors.New("messages not found")
	ErrMessageNotFound       = errors.New("message not found")
)

type MessageRepo struct {
	db *DB
}

func NewMessageRepo(db *DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (repo *MessageRepo) CreateMessage(message *model.DBMessage) (int, error) {
	var id int

	row := repo.db.QueryRow(`
		INSERT INTO message(channel_id, user_id, title, message_url, imageurl) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		message.ChannelID, message.UserID, message.Title,
		message.MessageURL, message.ImageURL,
	)
	if err := row.Scan(&id); err != nil {
		return id, fmt.Errorf("failed to create message: %w", err)
	}

	return id, nil
}

func (repo *MessageRepo) GetMessageByTitle(title string) (*model.DBMessage, error) {
	var message model.DBMessage

	err := repo.db.Get(&message, "SELECT * FROM message WHERE title = $1;", title)
	if err == sql.ErrNoRows {
		return nil, ErrMessageNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("error while getting message by title: %w", err)
	}

	return &message, nil
}

func (repo *MessageRepo) GetMessagesCount() (int, error) {
	var count int

	err := repo.db.Get(&count, "SELECT COUNT(*) FROM message;")
	if err == sql.ErrNoRows {
		return 0, ErrMessagesCountNotFound
	}

	if err != nil {
		return 0, fmt.Errorf("error while getting messages count: %w", err)
	}

	return count, nil
}

func (repo *MessageRepo) GetMessagesCountByChannelID(channelID int) (int, error) {
	var count int
	err := repo.db.Get(
		&count,
		`SELECT COUNT(*) FROM message m LEFT JOIN channel c ON c.id = m.channel_id WHERE m.channel_id = $1;`,
		channelID,
	)
	if err == sql.ErrNoRows {
		return 0, ErrMessagesCountNotFound
	}

	if err != nil {
		return 0, fmt.Errorf("error while getting messages count by channel ID: %w", err)
	}

	return count, nil
}

func (repo *MessageRepo) GetFullMessagesByPage(page int) ([]model.FullMessage, error) {
	messages := make([]model.FullMessage, 0, 10)

	err := repo.db.Select(
		&messages,
		`SELECT m.id, m.title, m.message_url, m.imageurl, 
		 c.id AS channelid, c.name AS channelname, c.imageurl AS channelimageurl, 
		 u.id AS userid, u.fullname, u.imageurl AS userimageurl, 
		 (SELECT COUNT(*) FROM replie WHERE message_id = m.id)
		 FROM message m 
		 LEFT JOIN channel c ON c.id = m.channel_id 
		 LEFT JOIN tg_user u ON u.id = m.user_id
		 ORDER BY m.id DESC NULLS LAST LIMIT 10 OFFSET $1;`,
		page,
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting full messages by page: %w", err)
	}

	if len(messages) == 0 {
		return nil, ErrMessagesNotFound
	}

	return messages, nil
}

func (repo *MessageRepo) GetFullMessagesByChannelIDAndPage(channelID, page int) ([]model.FullMessage, error) {
	messages := make([]model.FullMessage, 0, 10)

	err := repo.db.Select(
		&messages,
		`SELECT m.id, m.title, m.message_url, m.imageurl, 
		 c.id AS channelid, c.name AS channelname, c.imageurl AS channelimageurl, 
		 u.id AS userid, u.fullname, u.imageurl AS userimageurl, 
		 (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
		 FROM message m 
		 LEFT JOIN channel c ON c.id = m.channel_id 
		 LEFT JOIN tg_user u ON u.id = m.user_id
	 	 WHERE m.channel_id = $1 
		 ORDER BY count DESC NULLS LAST LIMIT 10 OFFSET $2;`,
		channelID, page,
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting full messages by channel ID: %w", err)
	}

	if len(messages) == 0 {
		return nil, ErrMessagesNotFound
	}

	return messages, nil
}

func (repo *MessageRepo) GetFullMessagesByUserID(ID int) ([]model.FullMessage, error) {
	messages := make([]model.FullMessage, 0, 5)

	err := repo.db.Select(
		&messages,
		`SELECT m.id, m.title, m.message_url, m.imageurl, 
		 c.id AS channelid, c.name AS channelname, c.Title AS channeltitle, c.imageurl AS channelimageurl, 
		 (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
		 FROM message m 
		 LEFT JOIN channel c ON c.id = m.channel_id 
		 LEFT JOIN tg_user u ON u.id = m.user_id
		 WHERE m.user_id= $1 
		 ORDER BY count DESC NULLS LAST;`,
		ID,
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting full messages by user ID: %w", err)
	}

	if len(messages) == 0 {
		return nil, ErrMessagesNotFound
	}

	return messages, nil
}

func (repo *MessageRepo) GetFullMessageByMessageID(messageID int) (*model.FullMessage, error) {
	var message model.FullMessage

	err := repo.db.Get(
		&message,
		`SELECT m.id, m.title, m.message_url, m.imageurl, 
		 c.id AS channelid, c.name AS channelname, c.title as channeltitle, c.imageurl as channelimageurl, 
		 u.id as userid, u.fullname, u.imageurl as userimageurl, 
		 (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
		 FROM message m 
		 LEFT JOIN channel c ON c.id = m.channel_id 
		 LEFT JOIN tg_user u ON u.id = m.user_id
		 WHERE m.id = $1;`,
		messageID,
	)
	if err == sql.ErrNoRows {
		return nil, ErrMessageNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("error while getting full message by message ID: %w", err)
	}

	return &message, nil
}
