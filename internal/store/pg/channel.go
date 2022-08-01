package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

var (
	ErrChannelNotFound  = errors.New("channel not found")
	ErrChannelsNotFound = errors.New("channels not found")
)

type ChannelPgRepo struct {
	db *DB
}

func NewChannelRepo(db *DB) *ChannelPgRepo {
	return &ChannelPgRepo{db}
}

func (repo *ChannelPgRepo) CreateChannel(channel *model.DBChannel) error {
	var id int

	row := repo.db.QueryRow(`
		INSERT INTO channel(name, title, imageurl) 
		VALUES ($1, $2, $3) RETURNING id;`,
		channel.Name, channel.Title, channel.ImageURL,
	)
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	return nil
}

func (repo *ChannelPgRepo) GetChannels() ([]model.Channel, error) {
	channels := make([]model.Channel, 0, 10)

	err := repo.db.Select(&channels, "SELECT * FROM channel;")
	if err != nil {
		return nil, fmt.Errorf("error while getting channels: %w", err)
	}

	if len(channels) == 0 {
		return nil, ErrChannelsNotFound
	}

	return channels, nil
}

func (repo *ChannelPgRepo) GetChannelsByPage(page int) ([]model.Channel, error) {
	channels := make([]model.Channel, 0, 10)

	err := repo.db.Select(&channels, "SELECT * FROM channel LIMIT 10 OFFSET $1;", page)
	if err != nil {
		return nil, fmt.Errorf("error while getting channels by page: %w", err)
	}

	if len(channels) == 0 {
		return nil, ErrChannelsNotFound
	}

	return channels, nil
}

func (repo *ChannelPgRepo) GetChannelByName(name string) (*model.Channel, error) {
	var channel model.Channel

	err := repo.db.Get(&channel, "SELECT * FROM channel WHERE name=$1;", name)
	if err == sql.ErrNoRows {
		return nil, ErrChannelNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("error while getting channel by name: %w", err)
	}

	return &channel, nil
}

func (repo *ChannelPgRepo) GetChannelStats(channelID int) (*model.Stat, error) {
	var sum int

	stat := &model.Stat{}

	messageCount := make([]int, 0)
	replieCount := make([]int, 0)

	rows, err := repo.db.Query(
		`SELECT m.id, COUNT(r.id) 
		 FROM channel c LEFT JOIN message m ON m.channel_id = c.id 
		 LEFT JOIN replie r ON r.message_id = m.id 
		 WHERE c.id = $1 GROUP BY m.id;`,
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
