package pg

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type ReplyRepo struct {
	db *DB
}

func NewReplyRepo(db *DB) *ReplyRepo {
	return &ReplyRepo{db: db}
}

func (repo *ReplyRepo) CreateReply(reply *model.DBReply) error {
	_, err := repo.db.Exec(`
		INSERT INTO reply(user_id, message_id, title, image_url) VALUES ($1, $2, $3, $4);`,
		reply.UserID, reply.MessageID, reply.Title, reply.ImageURL,
	)
	if err != nil {
		return fmt.Errorf("create reply: %w", err)
	}

	return nil
}

func (repo *ReplyRepo) GetFullRepliesByMessageID(messageID int) ([]model.FullReply, error) {
	var replies []model.FullReply

	err := repo.db.Select(
		&replies,
		`SELECT r.id, r.title, r.image_url, 
		 u.id as user_id, u.fullname, u.image_url as user_image_url
		 FROM reply r 
		 LEFT JOIN tg_user u ON r.user_id = u.id 
		 WHERE r.message_id = $1 
		 ORDER BY r.id DESC NULLS LAST;`,
		messageID,
	)
	if err != nil {
		return nil, fmt.Errorf("get full replies by message id: %w", err)
	}

	if len(replies) == 0 {
		return nil, nil
	}

	return replies, nil
}
