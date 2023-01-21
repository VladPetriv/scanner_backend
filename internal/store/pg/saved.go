package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/VladPetriv/scanner_backend/internal/model"
)

type SavedRepo struct {
	db *DB
}

func NewSavedRepo(db *DB) *SavedRepo {
	return &SavedRepo{db: db}
}

func (repo SavedRepo) CreateSavedMessage(saved *model.Saved) error {
	_, err := repo.db.Exec(
		"INSERT INTO saved(user_id, message_id) VALUES ($1, $2);",
		saved.WebUserID, saved.MessageID,
	)
	if err != nil {
		return fmt.Errorf("create saved message: %w", err)
	}

	return nil
}

func (repo SavedRepo) GetSavedMessages(userID int) ([]model.Saved, error) {
	var savedMessages []model.Saved

	err := repo.db.Select(&savedMessages, "SELECT * FROM saved WHERE user_id = $1;", userID)
	if err != nil {
		return nil, fmt.Errorf("get saved messages: %w", err)
	}

	if len(savedMessages) == 0 {
		return nil, nil
	}

	return savedMessages, nil
}

func (repo SavedRepo) GetSavedMessageByID(id int) (*model.Saved, error) {
	var savedMessage model.Saved

	err := repo.db.Get(&savedMessage, "SELECT * FROM saved WHERE message_id = $1;", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("get saved message by id: %w", err)
	}

	return &savedMessage, nil
}

func (repo SavedRepo) DeleteSavedMessage(id int) error {
	_, err := repo.db.Exec("DELETE FROM saved WHERE id = $1;", id)
	if err != nil {
		return fmt.Errorf("delete saved message: %w", err)
	}

	return nil
}
