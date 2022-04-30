package pg

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type ChannelPgRepo struct {
	db *DB
}

func NewChannelRepo(db *DB) *ChannelPgRepo {
	return &ChannelPgRepo{db}
}

func (repo *ChannelPgRepo) GetChannels() ([]model.Channel, error) {
	channels := make([]model.Channel, 0)
	rows, err := repo.db.Query("SELECT * FROM channel;")
	if err != nil {
		return nil, fmt.Errorf("error while getting channels: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		channel := model.Channel{}
		err := rows.Scan(&channel.ID, &channel.Name, &channel.Title, &channel.PhotoURL)
		if err != nil {
			continue
		}

		channels = append(channels, channel)
	}
	if len(channels) == 0 {
		return nil, fmt.Errorf("channels not found")
	}

	return channels, nil
}

func (repo *ChannelPgRepo) GetChannel(channelID int) (*model.Channel, error) {
	channel := &model.Channel{}

	rows, err := repo.db.Query("SELECT * FROM channel WHERE id=$1;", channelID)
	if err != nil {
		return nil, fmt.Errorf("error while getting channel: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&channel.ID, &channel.Name, &channel.Title, &channel.PhotoURL)
		if err != nil {
			continue
		}
	}

	if channel.Name == "" {
		return nil, fmt.Errorf("channel not found")
	}

	return channel, nil
}
func (repo *ChannelPgRepo) GetChannelByName(name string) (*model.Channel, error) {
	channel := &model.Channel{}

	rows, err := repo.db.Query("SELECT * FROM channel WHERE name=$1;", name)
	if err != nil {
		return nil, fmt.Errorf("error while getting channel: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&channel.ID, &channel.Name, &channel.Title, &channel.PhotoURL)
		if err != nil {
			continue
		}
	}

	if channel.Name == "" {
		return nil, nil
	}

	return channel, nil
}
