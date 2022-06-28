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
		return nil, nil
	}

	return channels, nil
}

func (repo *ChannelPgRepo) GetChannelsByPage(page int) ([]model.Channel, error) {
	channels := make([]model.Channel, 0)

	rows, err := repo.db.Query("SELECT * FROM channel LIMIT 10 OFFSET $1;", page)
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
		return nil, nil
	}

	return channels, nil
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

func (repo *ChannelPgRepo) GetChannelStats(channelID int) (*model.Stat, error) {
	var sum int
	stat := &model.Stat{}
	messageCount := make([]int, 0)
	replieCount := make([]int, 0)

	rows, err := repo.db.Query(
		"SELECT m.id, COUNT(r.id) FROM channel c LEFT JOIN message m ON m.channel_id = c.id LEFT JOIN replie r ON r.message_id = m.id WHERE c.id = $1 GROUP BY m.id;",
		channelID,
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting channel stats: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var mC int
		var rC int

		err := rows.Scan(&mC, &rC)
		if err != nil {
			continue
		}

		messageCount = append(messageCount, mC)
		replieCount = append(replieCount, rC)
	}

	stat.MessagesCount = len(messageCount)

	for _, replie := range replieCount {
		sum += replie
	}

	stat.RepliesCount = sum

	return stat, nil
}
