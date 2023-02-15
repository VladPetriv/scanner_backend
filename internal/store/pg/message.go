package pg

import (
	"database/sql"
	"errors"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type MessageRepo struct {
	db *DB
}

func NewMessageRepo(db *DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (repo MessageRepo) CreateMessage(message *model.DBMessage) (int, error) {
	var id int

	row := repo.db.QueryRow(`
		INSERT INTO message(channel_id, user_id, title, message_url, image_url) VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		message.ChannelID, message.UserID, message.Title,
		message.MessageURL, message.ImageURL,
	)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo MessageRepo) GetMessagesCount() (int, error) {
	var count int

	err := repo.db.Get(&count, "SELECT COUNT(*) FROM message;")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return count, nil
}

func (repo MessageRepo) GetMessagesCountByChannelID(channelID int) (int, error) {
	var count int

	err := repo.db.Get(
		&count,
		`SELECT COUNT(*) FROM message m LEFT JOIN channel c ON c.id = m.channel_id WHERE m.channel_id = $1;`,
		channelID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return count, nil
}

func (repo MessageRepo) GetMessageByTitle(title string) (*model.DBMessage, error) {
	var message model.DBMessage

	err := repo.db.Get(&message, "SELECT * FROM message WHERE title = $1;", title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &message, nil
}

func (repo MessageRepo) GetFullMessagesByPage(page int) ([]model.FullMessage, error) {
	var messages []model.FullMessage

	err := repo.db.Select(
		&messages,
		`SELECT m.id, m.title, m.message_url, m.image_url, 
		 c.id AS channel_id, c.name AS channel_name, c.image_url AS channel_image_url, 
		 u.id AS user_id, u.fullname, u.image_url AS user_image_url, 
		 (SELECT COUNT(*) FROM reply WHERE message_id = m.id)
		 FROM message m 
		 LEFT JOIN channel c ON c.id = m.channel_id 
		 LEFT JOIN tg_user u ON u.id = m.user_id
		 ORDER BY m.id DESC NULLS LAST LIMIT 10 OFFSET $1;`,
		page,
	)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages, nil
}

func (repo MessageRepo) GetFullMessagesByChannelIDAndPage(channelID, page int) ([]model.FullMessage, error) {
	var messages []model.FullMessage

	err := repo.db.Select(
		&messages,
		`SELECT m.id, m.title, m.message_url, m.image_url, 
		 c.id AS channel_id, c.name AS channel_name, c.image_url AS channel_image_url, 
		 u.id AS user_id, u.fullname, u.image_url AS user_image_url, 
		 (SELECT COUNT(id) FROM reply WHERE message_id = m.id)
		 FROM message m 
		 LEFT JOIN channel c ON c.id = m.channel_id 
		 LEFT JOIN tg_user u ON u.id = m.user_id
	 	 WHERE m.channel_id = $1 
		 ORDER BY count DESC NULLS LAST LIMIT 10 OFFSET $2;`,
		channelID, page,
	)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages, nil
}

func (repo MessageRepo) GetFullMessagesByUserID(id int) ([]model.FullMessage, error) {
	var messages []model.FullMessage

	err := repo.db.Select(
		&messages,
		`SELECT m.id, m.title, m.message_url, m.image_url, 
		 c.id AS channel_id, c.name AS channel_name, c.Title AS channel_title, c.image_url AS channel_image_url, 
		 (SELECT COUNT(id) FROM reply WHERE message_id = m.id)
		 FROM message m 
		 LEFT JOIN channel c ON c.id = m.channel_id 
		 LEFT JOIN tg_user u ON u.id = m.user_id
		 WHERE m.user_id= $1 
		 ORDER BY count DESC NULLS LAST;`,
		id,
	)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages, nil
}

func (repo MessageRepo) GetFullMessageByID(messageID int) (*model.FullMessage, error) {
	var message model.FullMessage

	err := repo.db.Get(
		&message,
		`SELECT m.id, m.title, m.message_url, m.image_url, 
		 c.id AS channel_id, c.name AS channel_name, c.title as channel_title, c.image_url as channel_image_url, 
		 u.id as user_id, u.fullname, u.image_url as user_image_url, 
		 (SELECT COUNT(id) FROM reply WHERE message_id = m.id)
		 FROM message m 
		 LEFT JOIN channel c ON c.id = m.channel_id 
		 LEFT JOIN tg_user u ON u.id = m.user_id
		 WHERE m.id = $1;`,
		messageID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &message, nil
}
