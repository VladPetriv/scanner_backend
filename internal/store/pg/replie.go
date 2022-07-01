package pg

import (
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type ReplieRepo struct {
	db *DB
}

func NewReplieRepo(db *DB) *ReplieRepo {
	return &ReplieRepo{db: db}
}

func (repo *ReplieRepo) GetFullRepliesByMessageID(ID int) ([]model.FullReplie, error) {
	replies := make([]model.FullReplie, 0)

	rows, err := repo.db.Query(
		"SELECT r.id, r.title, u.id, u.fullname, u.imageurl FROM replie r LEFT JOIN tg_user u ON r.user_id = u.id WHERE r.message_id = $1 ORDER BY r.id DESC NULLS LAST;", ID,
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting full replies by message ID: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		replie := model.FullReplie{}
		err := rows.Scan(&replie.ID, &replie.Title, &replie.UserID, &replie.FullName, &replie.UserImageURL)
		if err != nil {
			continue
		}

		replies = append(replies, replie)
	}

	if len(replies) == 0 {
		return nil, nil
	}

	return replies, nil
}
