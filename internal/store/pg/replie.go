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

func (repo *ReplieRepo) GetReplies() ([]model.Replie, error) {
	replies := make([]model.Replie, 0)
	rows, err := repo.db.Query("SELECT * FROM replie;")
	if err != nil {
		return nil, fmt.Errorf("error while getting replies: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		replie := model.Replie{}
		err := rows.Scan(&replie.ID, &replie.UserID, &replie.MessageID, &replie.Title)
		if err != nil {
			continue
		}

		replies = append(replies, replie)
	}

	if len(replies) == 0 {
		return nil, fmt.Errorf("replies not found")
	}

	return replies, nil
}

func (repo *ReplieRepo) GetReplie(replieId int) (*model.Replie, error) {
	replie := &model.Replie{}

	rows, err := repo.db.Query("SELECT * FROM replie WHERE id=$1;", replieId)
	if err != nil {
		return nil, fmt.Errorf("error while getting replie: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&replie.ID, &replie.UserID, &replie.MessageID, &replie.Title)
		if err != nil {
			continue
		}
	}

	if replie.Title == "" {
		return nil, fmt.Errorf("replie not found")
	}

	return replie, nil
}

func (repo *ReplieRepo) GetReplieByName(name string) (*model.Replie, error) {
	replie := &model.Replie{}

	rows, err := repo.db.Query("SELECT * FROM replie WHERE title=$1;", name)
	if err != nil {
		return nil, fmt.Errorf("error while getting replie: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&replie.ID, &replie.UserID, &replie.MessageID, &replie.Title)
		if err != nil {
			continue
		}
	}

	if replie.Title == "" {
		return nil, nil
	}

	return replie, nil
}

func (repo *ReplieRepo) GetFullRepliesByMessageID(ID int) ([]model.FullReplie, error) {
	replies := make([]model.FullReplie, 0)

	rows, err := repo.db.Query(
		"SELECT r.id, r.title, u.id, u.fullname, u.photourl FROM replie r LEFT JOIN tg_user u ON r.user_id = u.id WHERE r.message_id = $1 ORDER BY r.id DESC NULLS LAST;", ID,
	)
	if err != nil {
		return nil, fmt.Errorf("error while getting full replies by message ID: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		replie := model.FullReplie{}
		err := rows.Scan(&replie.ID, &replie.Title, &replie.UserID, &replie.FullName, &replie.PhotoURL)
		if err != nil {
			continue
		}

		replies = append(replies, replie)
	}

	if len(replies) == 0 {
		return nil, fmt.Errorf("replies not found")
	}

	return replies, nil
}
