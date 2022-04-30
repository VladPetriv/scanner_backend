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

func (repo *MessageRepo) GetMessages() ([]model.Message, error) {
	messages := make([]model.Message, 0)

	rows, err := repo.db.Query("SELECT * FROM message;")
	if err != nil {
		return nil, fmt.Errorf("error while getting messages: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		message := model.Message{}
		err := rows.Scan(&message.ID, &message.UserID, &message.ChannelID, &message.Title)
		if err != nil {
			continue
		}

		messages = append(messages, message)
	}

	if len(messages) == 0 {
		return nil, fmt.Errorf("messages not found")
	}

	return messages, nil
}

func (repo *MessageRepo) GetMessage(messageID int) (*model.Message, error) {
	message := &model.Message{}

	rows, err := repo.db.Query("SELECT * FROM message WHERE id=$1;", messageID)
	if err != nil {
		return nil, fmt.Errorf("error while getting message: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&message.ID, &message.UserID, &message.ChannelID, &message.Title)
		if err != nil {
			continue
		}
	}

	if message.Title == "" {
		return nil, fmt.Errorf("message not found")
	}

	return message, nil
}

func (repo *MessageRepo) GetMessageByName(name string) (*model.Message, error) {
	message := &model.Message{}

	rows, err := repo.db.Query("SELECT * FROM message WHERE title=$1;", name)
	if err != nil {
		return nil, fmt.Errorf("error while getting message: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&message.ID, &message.UserID, &message.ChannelID, &message.Title)
		if err != nil {
			continue
		}
	}

	if message.Title == "" {
		return nil, nil
	}

	return message, nil
}

func (repo *MessageRepo) GetFullMessages() ([]model.FullMessage, error) {
	messages := make([]model.FullMessage, 0)

	rows, err := repo.db.Query("SELECT m.id, m.Title, c.id, c.Name, c.Photourl as channelPhotoUrl, u.Fullname, u.Photourl, (SELECT COUNT(id) FROM replie WHERE message_id = m.id)  FROM message m LEFT JOIN channel c ON c.id = m.channel_id LEFT JOIN tg_user u ON u.id = m.user_id;")
	if err != nil {
		return nil, fmt.Errorf("error while getting full messages: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		message := model.FullMessage{}

		err := rows.Scan(&message.ID, &message.Title, &message.ChannelID, &message.ChannelName, &message.ChannelPhotoURL, &message.FullName, &message.PhotoURL, &message.ReplieCount)
		if err != nil {
			continue
		}

		messages = append(messages, message)
	}

	if len(messages) == 0 {
		return nil, fmt.Errorf("messages not found")
	}

	return messages, nil
}
