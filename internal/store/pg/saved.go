package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

var (
	ErrSavedMessagesNotFound = errors.New("saved messages not found")
	ErrSavedMessageNotFound  = errors.New("saved message not found")
)

type SavedRepo struct {
	db *DB
}

func NewSavedRepo(db *DB) *SavedRepo {
	return &SavedRepo{db: db}
}

func (repo *SavedRepo) GetSavedMessages(userID int) ([]model.Saved, error) {
	savedMessages := make([]model.Saved, 0, 5)

	err := repo.db.Select(&savedMessages, "SELECT * FROM saved WHERE user_id = $1;", userID)
	if err != nil {
		return nil, fmt.Errorf("error while getting saved messages: %w", err)
	}

	if len(savedMessages) == 0 {
		return nil, ErrSavedMessagesNotFound
	}

	return savedMessages, nil
}

func (repo *SavedRepo) GetSavedMessageByMessageID(messageID int) (*model.Saved, error) {
	var savedMessage model.Saved

	err := repo.db.Get(&savedMessage, "SELECT * FROM saved WHERE message_id = $1;", messageID)
	if err == sql.ErrNoRows {
		return nil, ErrSavedMessageNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("error while getting saved message: %w", err)
	}

	return &savedMessage, nil
}

func (repo *SavedRepo) CreateSavedMessage(saved *model.Saved) (int, error) {
	var id int

	row := repo.db.QueryRow("INSERT INTO saved(user_id, message_id) VALUES ($1, $2) RETURNING id;", saved.WebUserID, saved.MessageID)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error while creating saved message: %w", err)
	}

	return id, nil
}

func (repo *SavedRepo) DeleteSavedMessage(ID int) (int, error) {
	var id int

	row := repo.db.QueryRow("DELETE FROM saved WHERE id=$1 RETURNING id;", ID)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error while deleting saved message: %w", err)
	}

	return id, nil
}
