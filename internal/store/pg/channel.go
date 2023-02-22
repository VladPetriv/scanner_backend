package pg

import (
	"database/sql"
	"errors"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type ChannelPgRepo struct {
	db *DB
}

func NewChannelRepo(db *DB) *ChannelPgRepo {
	return &ChannelPgRepo{db}
}

func (repo ChannelPgRepo) CreateChannel(channel *model.DBChannel) error {
	_, err := repo.db.Exec(`
		INSERT INTO channel(name, title, image_url) VALUES ($1, $2, $3);`,
		channel.Name, channel.Title, channel.ImageURL,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo ChannelPgRepo) GetChannels() ([]model.Channel, error) {
	var channels []model.Channel

	err := repo.db.Select(&channels, "SELECT * FROM channel;")
	if err != nil {
		return nil, err
	}

	if len(channels) == 0 {
		return nil, nil
	}

	return channels, nil
}

func (repo ChannelPgRepo) GetChannelsByPage(page int) ([]model.Channel, error) {
	var channels []model.Channel

	err := repo.db.Select(&channels, "SELECT * FROM channel LIMIT 10 OFFSET $1;", page)
	if err != nil {
		return nil, err
	}

	if len(channels) == 0 {
		return nil, nil
	}

	return channels, nil
}

func (repo ChannelPgRepo) GetChannelByName(name string) (*model.Channel, error) {
	var channel model.Channel

	err := repo.db.Get(&channel, "SELECT * FROM channel WHERE name = $1;", name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &channel, nil
}

func (repo ChannelPgRepo) GetChannelStats(channelID int) (*model.Stat, error) {
	var sum int

	stat := &model.Stat{}

	messageCount := make([]int, 0)
	replyCount := make([]int, 0)

	rows, err := repo.db.Query(
		`SELECT m.id, COUNT(r.id) 
		 FROM channel c LEFT JOIN message m ON m.channel_id = c.id 
		 LEFT JOIN reply r ON r.message_id = m.id 
		 WHERE c.id = $1 GROUP BY m.id;`,
		channelID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var mC int
		var rC int

		err = rows.Scan(&mC, &rC)
		if err != nil {
			continue
		}

		messageCount = append(messageCount, mC)
		replyCount = append(replyCount, rC)
	}

	stat.MessagesCount = len(messageCount)

	for _, reply := range replyCount {
		sum += reply
	}

	stat.RepliesCount = sum

	return stat, nil
}
