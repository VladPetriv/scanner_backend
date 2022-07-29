package pg

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
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

func (repo *MessageRepo) GetMessagesLength() (int, error) {
	var count int

	row := repo.db.QueryRow("SELECT COUNT(*) FROM message;")
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("error while getting messages length: %w", err)
	}

	return count, nil
}

func (repo *MessageRepo) GetFullMessages(page int) ([]model.FullMessage, error) {
	messages := make([]model.FullMessage, 0)

	rows, err := repo.db.Query(
		`SELECT m.id, m.Title, m.message_url, m.imageurl, 
		c.id, c.Name, c.imageurl as channelImageUrl, 
		u.id, u.Fullname, u.imageurl as userImageUrl, 
		(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
	  FROM message m 
		LEFT JOIN channel c ON c.id = m.channel_id 
		LEFT JOIN tg_user u ON u.id = m.user_id
		ORDER BY m.id DESC NULLS LAST LIMIT 10 OFFSET $1;`, page,
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting full messages: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		message := model.FullMessage{}

		err := rows.Scan(
			&message.ID, &message.Title, &message.MessageURL, &message.ImageURL, &message.ChannelID, &message.ChannelName,
			&message.ChannelImageURL, &message.UserID, &message.FullName, &message.UserImageURL, &message.ReplieCount,
		)
		if err != nil {
			continue
		}

		messages = append(messages, message)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages, nil
}

func (repo *MessageRepo) GetFullMessagesByChannelID(ID, limit, page int) ([]model.FullMessage, error) {
	messages := make([]model.FullMessage, 0)

	rows, err := repo.db.Query(
		`SELECT m.id, m.Title, m.message_url, m.imageurl, 
		c.id, c.Name, c.imageurl as channelImageUrl, 
		u.id, u.Fullname, u.imageurl as userImageUrl, 
		(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
		FROM message m 
		LEFT JOIN channel c ON c.id = m.channel_id 
		LEFT JOIN tg_user u ON u.id = m.user_id
		WHERE m.channel_id = $1 
		ORDER BY count DESC NULLS LAST LIMIT $2 OFFSET $3;`, ID, limit, page,
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting full messages by channel ID: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		message := model.FullMessage{}

		err := rows.Scan(
			&message.ID, &message.Title, &message.MessageURL, &message.ImageURL, &message.ChannelID, &message.ChannelName,
			&message.ChannelImageURL, &message.UserID, &message.FullName, &message.UserImageURL, &message.ReplieCount,
		)
		if err != nil {
			continue
		}

		messages = append(messages, message)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages, nil
}

func (repo *MessageRepo) GetMessagesLengthByChannelID(ID int) (int, error) {
	messages := make([]model.FullMessage, 0)

	rows, err := repo.db.Query(
		`SELECT m.id FROM message m LEFT JOIN channel c ON c.id = m.channel_id WHERE m.channel_id = $1;`, ID,
	)
	if err != nil {
		return 0, fmt.Errorf("error while getting full messages by channel ID: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		message := model.FullMessage{}

		err := rows.Scan(&message.ID)
		if err != nil {
			continue
		}

		messages = append(messages, message)
	}

	if len(messages) == 0 {
		return 0, nil
	}

	return len(messages), nil
}

func (repo *MessageRepo) GetFullMessagesByUserID(ID int) ([]model.FullMessage, error) {
	messages := make([]model.FullMessage, 0)

	rows, err := repo.db.Query(
		`SELECT m.id, m.Title, m.message_url, m.imageurl, 
		c.id, c.Name, c.Title, c.imageurl as channelImageUrl, 
		(SELECT COUNT(id) FROM replie WHERE message_id = m.id)
		FROM message m 
		LEFT JOIN channel c ON c.id = m.channel_id 
		LEFT JOIN tg_user u ON u.id = m.user_id
		WHERE m.user_id= $1 
		ORDER BY count DESC NULLS LAST;`, ID,
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting full messages by user ID: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		message := model.FullMessage{}

		err := rows.Scan(
			&message.ID, &message.Title, &message.MessageURL, &message.ImageURL, &message.ChannelID, &message.ChannelName,
			&message.ChannelTitle, &message.ChannelImageURL, &message.ReplieCount,
		)
		if err != nil {
			continue
		}

		messages = append(messages, message)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages, nil
}

func (repo *MessageRepo) GetFullMessageByMessageID(ID int) (*model.FullMessage, error) {
	message := &model.FullMessage{}

	rows, err := repo.db.Query(
		`SELECT m.id, m.Title, m.message_url, m.imageurl, 
		 c.id, c.Name, c.Title, c.imageurl as channelImageUrl, 
		 u.id, u.Fullname, u.imageurl as userImageUrl, 
		 (SELECT COUNT(id) FROM replie WHERE message_id = m.id)
		 FROM message m 
		 LEFT JOIN channel c ON c.id = m.channel_id 
		 LEFT JOIN tg_user u ON u.id = m.user_id
		 WHERE m.id = $1;`, ID,
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting full messages by message ID: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(
			&message.ID, &message.Title, &message.MessageURL, &message.ImageURL,
			&message.ChannelID, &message.ChannelName, &message.ChannelTitle, &message.ChannelImageURL,
			&message.UserID, &message.FullName, &message.UserImageURL, &message.ReplieCount,
		)
		if err != nil {
			continue
		}
	}

	if message.Title == "" {
		return nil, nil
	}

	return message, nil
}
